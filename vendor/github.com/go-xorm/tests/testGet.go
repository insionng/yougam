package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

type NoIdUser struct {
	User   string `xorm:"unique"`
	Remain int64
	Total  int64
}

func get(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{Uid: 2}

	has, err := engine.Get(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if has {
		fmt.Println(user)
	} else {
		fmt.Println("no record id is 2")
	}

	err = engine.Sync(&NoIdUser{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	userCol := engine.ColumnMapper.Obj2Table("User")

	_, err = engine.Where("`"+userCol+"` = ?", "xlw").Delete(&NoIdUser{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	cnt, err := engine.Insert(&NoIdUser{"xlw", 20, 100})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
	}

	noIdUser := new(NoIdUser)
	has, err = engine.Where("`"+userCol+"` = ?", "xlw").Get(noIdUser)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		err = errors.New("get not returned 1")
		t.Error(err)
		panic(err)
	}
	fmt.Println(noIdUser)
}
