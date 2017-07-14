package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/insionng/yougam/libraries/cli"
	"github.com/insionng/yougam/modules/app"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/macaron.v1"
)

func Test_Yougam(t *testing.T) {
	Convey("App service", t, func() {

		cliApp := cli.NewApp()
		cliApp.Name = "yougam"
		cliApp.Usage = "yougam 0.0.0.0 8000 true"
		cliApp.Action = func(self *cli.Context) {

			fmt.Println("------------------App-Test------------------")

			if e := app.App(self); e != nil {
				log.Fatalln(e)
			}
		}

		go cliApp.Run(os.Args)

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
		So(err, ShouldBeNil)

		m := macaron.New()
		m.ServeHTTP(resp, req)

	})

}
