package xorm

import (
    "reflect"

    "github.com/insionng/yougam/libraries/go-xorm/core"
)

var (
	ptrPkType = reflect.TypeOf(&core.PK{})
	pkType    = reflect.TypeOf(core.PK{})
)
