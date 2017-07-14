# xorm tests

Xorm is a simple and powerful ORM for Go.

[![Build Status](https://drone.io/github.com/go-xorm/tests/status.png)](https://drone.io/github.com/go-xorm/tests/latest)  [![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/go-xorm/xorm) [![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/lunny/xorm/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

This is the tests project for xorm.

Please add your test codes here if you want to pull request to xorm.

### How to run the tests

* test sqlite

 ```
 ./sqlite3.sh
 ```

* test mysql or mymysql

 Create empty databases `xorm_test`, `xorm_test1`, `xorm_test2`, `xorm_test3` on your mysql server and make an account root, and let passwd empty on localhost, and then run:

 ```
 ./mysql.sh
 ./mymysql.sh
 ```

* test postgres

 Create empty database `xorm_test` on your postgres and and let passwd empty for default account on localhost, and then run:

 ```
 ./postgres.sh
 ```

### Running tests by name

You can also use `run_tests.sh` script:

```
./run_tests.sh <db>  # e.g. mysql
```

Run all tests:

```
./run_tests.sh
```

### Running tests using docker

You can also use `run_tests_docker.sh` script that will pull and run preconfigured images with database engines, and run tests on them. With this approach, you don't need to configure anything, just install docker and run the script.

Run specific test:

```
./run_tests_docker.sh <db>  # e.g. mysql
```

Run specific test with given database version:

```
./run_tests_docker.sh <db>:<version>  # e.g. mysql:5.5
```

Run all tests:

```
./run_tests_docker.sh
```
