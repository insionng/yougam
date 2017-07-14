package main

import (
	"github.com/insionng/yougam/libraries/Sirupsen/logrus"
	"github.com/insionng/yougam/libraries/mattn/go-colorable"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())

	logrus.Info("succeeded")
	logrus.Warn("not correct")
	logrus.Error("something error")
	logrus.Fatal("panic")
}
