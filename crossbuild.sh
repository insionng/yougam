echo "::Building Yougam::"
echo "::Disable Cgo::"
find ./ -type f -path "*.go"|xargs sed -i 's:_ "yougam/libraries/mattn/go-sqlite3" // need cgo!://_ "yougam/libraries/mattn/go-sqlite3" // need cgo!:g'

echo "::Building Yougam for windows::"
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ./bin/yougam-win-32bit.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/yougam-win-64bit.exe

echo "::Building Yougam for linux::"
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./bin/yougam-linux-32bit.bin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/yougam-linux-64bit.bin

echo "::Building Yougam for darwin::"
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -ldflags "-s -w" -o ./bin/yougam-darwin-32bit.bin
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/yougam-darwin-64bit.bin

echo "::Building Yougam for freebsd::"
CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-s -w" -o ./bin/yougam-freebsd-32bit.bin
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/yougam-freebsd-64bit.bin

echo "::Building Yougam for linux ARM 5/6/7::"
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w" -o ./bin/yougam-linux-arm7.bin
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "-s -w" -o ./bin/yougam-linux-arm6.bin
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-s -w" -o ./bin/yougam-linux-arm5.bin

echo "::Enable Cgo::"
find ./ -type f -path "*.go"|xargs sed -i 's://_ "yougam/libraries/mattn/go-sqlite3" // need cgo!:_ "yougam/libraries/mattn/go-sqlite3" // need cgo!:g'
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/yougam-cgo-linux-64bit.bin

echo "Okay!"

echo "Moving the output binary file into root directory:./yougam/"
echo "And then setting the database configure in file:./yougam/conf/config.conf"
echo "Enjoy it!"