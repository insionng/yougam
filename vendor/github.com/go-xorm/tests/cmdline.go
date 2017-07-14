package tests

import (
	"flag"
)

var ConnectionPort string

func init() {
	flag.StringVar(&ConnectionPort, "port", "", "")
	flag.Parse()
}
