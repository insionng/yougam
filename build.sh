echo "::Building Yougam::"

find ./ -type f -path "*.go"|xargs sed -i 's://_ "yougam/libraries/mattn/go-sqlite3" // need cgo!:_ "yougam/libraries/mattn/go-sqlite3" // need cgo!:g'
#重复替换多次避免遗留
find ./ -type f -path "*.go"|xargs sed -i 's://_ "yougam/libraries/mattn/go-sqlite3" // need cgo!:_ "yougam/libraries/mattn/go-sqlite3" // need cgo!:g'
echo "::Disable Cgo::"
find ./ -type f -path "*.go"|xargs sed -i 's:_ "yougam/libraries/mattn/go-sqlite3" // need cgo!://_ "yougam/libraries/mattn/go-sqlite3" // need cgo!:g'

CGO_ENABLED="0" go build ./
echo "Okay!"