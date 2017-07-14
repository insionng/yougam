package tests

import (
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

func testIterate(engine *xorm.Engine, t *testing.T) {
	err := engine.Omit("is_man").Iterate(new(Userinfo), func(idx int, bean interface{}) error {
		user := bean.(*Userinfo)
		fmt.Println(idx, "--", user)
		return nil
	})

	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testRows(engine *xorm.Engine, t *testing.T) {
	rows, err := engine.Omit("is_man").Rows(new(Userinfo))
	if err != nil {
		t.Error(err)
		panic(err)
	}
	defer rows.Close()

	idx := 0
	user := new(Userinfo)
	for rows.Next() {
		err = rows.Scan(user)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		fmt.Println(idx, "--", user)
		idx++
	}
}
