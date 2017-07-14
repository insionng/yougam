package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

func testFind(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- find --------------")
	find(engine, t)
	fmt.Println("-------------- find2 --------------")
	find2(engine, t)
	fmt.Println("-------------- findMap --------------")
	findMap(engine, t)
	fmt.Println("-------------- findMap2 --------------")
	findMap2(engine, t)
	fmt.Println("-------------- findInts --------------")
	testFindInts(engine, t)
}

func where(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Where("(id) > ?", 2).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	err = engine.Where("(id) > ?", 2).And("(id) < ?", 10).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func find(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}

	users2 := make([]Userinfo, 0)
	userinfo := engine.TableMapper.Obj2Table("Userinfo")
	err = engine.Sql("select * from " + engine.Quote(userinfo)).Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func find2(engine *xorm.Engine, t *testing.T) {
	users := make([]*Userinfo, 0)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

type Team struct {
	Id int64
}

type TeamUser struct {
	OrgId  int64
	Uid    int64
	TeamId int64
}

func find3(engine *xorm.Engine, t *testing.T) {
	err := engine.Sync2(new(Team), new(TeamUser))
	if err != nil {
		t.Error(err)
		panic(err.Error())
	}

	var teams []Team
	err = engine.Cols("`team`.id").
		Where("`team_user`.org_id=?", 1).
		And("`team_user`.uid=?", 2).
		Join("INNER", "`team_user`", "`team_user`.team_id=`team`.id").
		Find(&teams)
	if err != nil {
		t.Error(err)
		panic(err.Error())
	}
}

func findMap(engine *xorm.Engine, t *testing.T) {
	users := make(map[int64]Userinfo)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

func findMap2(engine *xorm.Engine, t *testing.T) {
	users := make(map[int64]*Userinfo)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for id, user := range users {
		fmt.Println(id, user)
	}
}

func testDistinct(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	departname := engine.TableMapper.Obj2Table("Departname")
	err := engine.Distinct(departname).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users) != 1 {
		t.Error(err)
		panic(errors.New("should be one record"))
	}

	fmt.Println(users)

	type Depart struct {
		Departname string
	}

	users2 := make([]Depart, 0)
	err = engine.Distinct(departname).Table(new(Userinfo)).Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users2) != 1 {
		t.Error(err)
		panic(errors.New("should be one record"))
	}
	fmt.Println(users2)
}

func order(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.OrderBy("id desc").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	users2 := make([]Userinfo, 0)
	err = engine.Asc("id", "username").Desc("height").Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users2)
}

func having(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.GroupBy("username").Having("username='xlw'").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	/*users = make([]Userinfo, 0)
	err = engine.Cols("id, username").GroupBy("username").Having("username='xlw'").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)*/
}

func orderSameMapper(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.OrderBy("(id) desc").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	users2 := make([]Userinfo, 0)
	err = engine.Asc("(id)", "Username").Desc("Height").Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users2)
}

func havingSameMapper(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.GroupBy("`Username`").Having("`Username`='xlw'").Find(&users)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(users)
}

func testFindInts(engine *xorm.Engine, t *testing.T) {
	userinfo := engine.TableMapper.Obj2Table("Userinfo")
	var idsInt64 []int64
	err := engine.Table(userinfo).Cols("id").Desc("id").Find(&idsInt64)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsInt64)

	var idsInt32 []int32
	err = engine.Table(userinfo).Cols("id").Desc("id").Find(&idsInt32)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsInt32)
}
