#!/bin/bash
if [ "$#" == "0" ];then
    tests=('sqlite3' 'mysql' 'mymysql' 'postgres')
else
    tests=($@)
fi

mysqlcmd='for db in xorm_test xorm_test1 xorm_test2 xorm_test3
do
    mysql -h"$MYSQL_PORT_3306_TCP_ADDR" \
          -P"$MYSQL_PORT_3306_TCP_PORT" \
          -uroot --password="" \
          --execute="CREATE DATABASE IF NOT EXISTS $db DEFAULT CHARACTER SET utf8;"
done'

postgrescmd='for db in xorm_test
do
    psql -h "$POSTGRES_PORT_5432_TCP_ADDR" \
         -p "$POSTGRES_PORT_5432_TCP_PORT" \
         -U postgres \
         -c "DROP DATABASE IF EXISTS $db;"
    psql -h "$POSTGRES_PORT_5432_TCP_ADDR" \
         -p "$POSTGRES_PORT_5432_TCP_PORT" \
         -U postgres \
         -c "CREATE DATABASE $db;"
done'

parsename() {
    echo "$1" | sed -re 's/:.*$//'
}

parsever() {
    ver="$(echo "$1" | sed -re 's/^[^:]+(:([0-9]+[.][0-9]+))?$/\2/')"
    echo "${ver:=$2}"
}

getport() {
    num="$(echo "$2" | sed -re 's:[^0-9]::g')"
    expr "10000" "+" "$1" "+" "${num:=0}"
}

runtest() {
    t="$(parsename $1)"
    p="$2"
    ( cd $t && ./${t}.sh -timeout 20m -args -port $p )
}

for tst in "${tests[@]}"
do
    echo "* executing tests for $tst"
    case $tst in
        sqlite3)
            echo "* getting dependencies"
            go get -v github.com/mattn/go-sqlite3 || exit 1

            echo "* running tests"
            ./sqlite3.sh || exit 1
            ;;

        mysql*|mymysql*)
            version=$(parsever $tst "5.5")
            image="mysql:${version}"
            name="xorm_mysql_${version}"
            port="$(getport 1000 ${version})"

            if ! docker ps -a | grep -q "$name"
            then
                echo "* starting container $name"
                docker run --name $name \
                       -e MYSQL_ROOT_PASSWORD="" \
                       -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
                       -d \
                       -p $port:3306 \
                       $image \
                    || exit 1
            fi

            echo "* polling server"
            while :
            do
                echo -n "."
                docker run --link $name:mysql --rm mysql \
                       bash -c "$mysqlcmd" &>/dev/null \
                    && break
                sleep 1
            done
            echo
            echo "* created tables"

            echo "* getting dependencies"
            case $tst in
                mysql*)
                    go get -v github.com/go-sql-driver/mysql || exit 1
                    ;;
                mymysql*)
                    go get -v github.com/ziutek/mymysql/godrv || exit 1
                    ;;
            esac

            echo "* running tests"
            runtest $tst $port || exit 1
            ;;

        postgres)
            version=$(parsever $tst "latest")
            image="postgres:${version}"
            name="xorm_postgres_${version}"
            port="$(getport 2000 ${version})"

            if ! docker ps -a | grep -q "$name"
            then
                echo "* starting container $name"
                docker run --name $name \
                       -e POSTGRES_PASSWORD="" \
                       -d \
                       -p $port:5432 \
                       $image \
                    || exit 1
            fi

            echo "* polling server"
            while :
            do
                echo -n "."
                docker run --link $name:postgres --rm postgres \
                       bash -c "$postgrescmd" &>/dev/null \
                    && break
                sleep 1
            done
            echo
            echo "* created tables"

            echo "* getting dependencies"
            go get -v github.com/lib/pq || exit 1

            echo "* running tests"
            runtest $tst $port || exit 1
            ;;

        *)
            echo "unknown test type $t" 1>&2
            exit 1
    esac
done
