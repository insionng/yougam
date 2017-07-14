package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"github.com/insionng/yougam/libraries/cli"
	"github.com/insionng/yougam/modules/app"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

//main enter
func main() {

	cliApp := cli.NewApp()
	cliApp.Name = "yougam"
	cliApp.Usage = "yougam 0.0.0.0 8000 true"
	cliApp.Action = func(self *cli.Context) {

		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Usage: ./yougam.bin host port bool")
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("yougam.bin: yougam's Application.")
		fmt.Println("host: Listening on the host.")
		fmt.Println("port: Listening on the port.")
		fmt.Println("bool: The true/false for production/development version.")
		fmt.Println("-----------------------------------------------------------")

		if e := app.App(self); e != nil {
			log.Fatalln(e)
		}
	}
	cliApp.Run(os.Args)
}
