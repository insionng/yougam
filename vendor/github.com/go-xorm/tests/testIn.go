package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

func in(engine *xorm.Engine, t *testing.T) {
	var usrs []Userinfo
	err := engine.Limit(3).Find(&usrs)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(usrs) != 3 {
		err = errors.New("there are not 3 records")
		t.Error(err)
		panic(err)
	}

	var ids []int64
	var idsStr string
	for _, u := range usrs {
		ids = append(ids, u.Uid)
		idsStr = fmt.Sprintf("%d,", u.Uid)
	}
	idsStr = idsStr[:len(idsStr)-1]

	users := make([]Userinfo, 0)
	err = engine.In("(id)", ids[0], ids[1], ids[2]).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
	if len(users) != 3 {
		err = errors.New("in uses should be " + idsStr + " total 3")
		t.Error(err)
		panic(err)
	}

	users = make([]Userinfo, 0)
	err = engine.In("(id)", ids).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
	if len(users) != 3 {
		err = errors.New("in uses should be " + idsStr + " total 3")
		t.Error(err)
		panic(err)
	}

	for _, user := range users {
		if user.Uid != ids[0] && user.Uid != ids[1] && user.Uid != ids[2] {
			err = errors.New("in uses should be " + idsStr + " total 3")
			t.Error(err)
			panic(err)
		}
	}

	users = make([]Userinfo, 0)
	var idsInterface []interface{}
	for _, id := range ids {
		idsInterface = append(idsInterface, id)
	}

	department := "`" + engine.ColumnMapper.Obj2Table("Departname") + "`"
	err = engine.Where(department+" = ?", "dev").In("(id)", idsInterface...).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	if len(users) != 3 {
		err = errors.New("in uses should be " + idsStr + " total 3")
		t.Error(err)
		panic(err)
	}

	for _, user := range users {
		if user.Uid != ids[0] && user.Uid != ids[1] && user.Uid != ids[2] {
			err = errors.New("in uses should be " + idsStr + " total 3")
			t.Error(err)
			panic(err)
		}
	}

	dev := engine.ColumnMapper.Obj2Table("Dev")

	err = engine.In("(id)", 1).In("(id)", 2).In(department, dev).Find(&users)

	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	cnt, err := engine.In("(id)", ids[0]).Update(&Userinfo{Departname: "dev-"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update records not 1")
		t.Error(err)
		panic(err)
	}

	user := new(Userinfo)
	has, err := engine.Id(ids[0]).Get(user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("get record not 1")
		t.Error(err)
		panic(err)
	}
	if user.Departname != "dev-" {
		err = errors.New("update not success")
		t.Error(err)
		panic(err)
	}

	cnt, err = engine.In("(id)", ids[0]).Update(&Userinfo{Departname: "dev"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update records not 1")
		t.Error(err)
		panic(err)
	}

	cnt, err = engine.In("(id)", ids[1]).Delete(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("deleted records not 1")
		t.Error(err)
		panic(err)
	}
}
