package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
)

type Condi map[string]interface{}

type UpdateAllCols struct {
	Id     int64
	Bool   bool
	String string
	Ptr    *string
}

type UpdateMustCols struct {
	Id     int64
	Bool   bool
	String string
}

type UpdateIncr struct {
	Id  int64
	Cnt int
}

func update(engine *xorm.Engine, t *testing.T) {
	testUpdate1(engine, t)
	testUpdateUpdated(engine, t)
	testUpdateMap(engine, t)
}

func testUpdateMap(engine *xorm.Engine, t *testing.T) {
	_, err := engine.Table("update_must_cols").Where("id =?", 1).Update(map[string]interface{}{
		"bool": true,
	})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testUpdate1(engine *xorm.Engine, t *testing.T) {
	var ori Userinfo
	has, err := engine.Get(&ori)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		t.Error(errors.New("not exist"))
		panic(errors.New("not exist"))
	}
	// update by id
	user := Userinfo{Username: "xxx", Height: 1.2}
	cnt, err := engine.Id(ori.Uid).Update(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	condi := Condi{"username": "zzz", "departname": ""}
	cnt, err = engine.Table(&user).Id(ori.Uid).Update(&condi)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	cnt, err = engine.Update(&Userinfo{Username: "yyy"}, &user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	total, err := engine.Count(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != total {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	// nullable update
	{
		user := &Userinfo{Username: "not null data", Height: 180.5}
		_, err := engine.Insert(user)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		userID := user.Uid

		has, err := engine.Id(userID).
			And("username = ?", user.Username).
			And("height = ?", user.Height).
			And("departname = ?", "").
			And("detail_id = ?", 0).
			And("is_man = ?", 0).
			And("created IS NOT NULL").
			Get(&Userinfo{})
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if !has {
			err = errors.New("cannot insert properly")
			t.Error(err)
			panic(err)
		}

		updatedUser := &Userinfo{Username: "null data"}
		cnt, err = engine.Id(userID).
			Nullable("height", "departname", "is_man", "created").
			Update(updatedUser)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			err = errors.New("update not returned 1")
			t.Error(err)
			panic(err)
		}

		has, err = engine.Id(userID).
			And("username = ?", updatedUser.Username).
			And("height IS NULL").
			And("departname IS NULL").
			And("is_man IS NULL").
			And("created IS NULL").
			And("detail_id = ?", 0).
			Get(&Userinfo{})
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if !has {
			err = errors.New("cannot update with null properly")
			t.Error(err)
			panic(err)
		}

		cnt, err = engine.Id(userID).Delete(&Userinfo{})
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			err = errors.New("delete not returned 1")
			t.Error(err)
			panic(err)
		}
	}

	err = engine.StoreEngine("Innodb").Sync2(&Article{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	defer func() {
		err = engine.DropTables(&Article{})
		if err != nil {
			t.Error(err)
			panic(err)
		}
	}()

	a := &Article{0, "1", "2", "3", "4", "5", 2}
	cnt, err = engine.Insert(a)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New(fmt.Sprintf("insert not returned 1 but %d", cnt))
		t.Error(err)
		panic(err)
	}

	if a.Id == 0 {
		err = errors.New("insert returned id is 0")
		t.Error(err)
		panic(err)
	}

	cnt, err = engine.Id(a.Id).Update(&Article{Name: "6"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New(fmt.Sprintf("insert not returned 1 but %d", cnt))
		t.Error(err)
		panic(err)
		return
	}

	var s = "test"

	col1 := &UpdateAllCols{Ptr: &s}
	err = engine.Sync(col1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(col1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	col2 := &UpdateAllCols{col1.Id, true, "", nil}
	_, err = engine.Id(col2.Id).AllCols().Update(col2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	col3 := &UpdateAllCols{}
	has, err = engine.Id(col2.Id).Get(col3)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		err = errors.New(fmt.Sprintf("cannot get id %d", col2.Id))
		t.Error(err)
		panic(err)
		return
	}

	if *col2 != *col3 {
		err = errors.New(fmt.Sprintf("col2 should eq col3"))
		t.Error(err)
		panic(err)
		return
	}

	{

		col1 := &UpdateMustCols{}
		err = engine.Sync(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		_, err = engine.Insert(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		col2 := &UpdateMustCols{col1.Id, true, ""}
		boolStr := engine.ColumnMapper.Obj2Table("Bool")
		stringStr := engine.ColumnMapper.Obj2Table("String")
		_, err = engine.Id(col2.Id).MustCols(boolStr, stringStr).Update(col2)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		col3 := &UpdateMustCols{}
		has, err := engine.Id(col2.Id).Get(col3)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		if !has {
			err = errors.New(fmt.Sprintf("cannot get id %d", col2.Id))
			t.Error(err)
			panic(err)
			return
		}

		if *col2 != *col3 {
			err = errors.New(fmt.Sprintf("col2 should eq col3"))
			t.Error(err)
			panic(err)
			return
		}
	}

	{

		col1 := &UpdateIncr{}
		err = engine.Sync(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		_, err = engine.Insert(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		cnt, err := engine.Id(col1.Id).Incr("cnt").Update(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			err = errors.New("update incr failed")
			t.Error(err)
			panic(err)
		}

		newCol := new(UpdateIncr)
		has, err := engine.Id(col1.Id).Get(newCol)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if !has {
			err = errors.New("has incr failed")
			t.Error(err)
			panic(err)
		}
		if 1 != newCol.Cnt {
			err = fmt.Errorf("incr failed %v %v %v", newCol.Cnt, newCol, col1)
			t.Error(err)
			panic(err)
		}
	}
}

type UpdatedUpdate struct {
	Id      int64
	Updated time.Time `xorm:"updated"`
}

type UpdatedUpdate2 struct {
	Id      int64
	Updated int64 `xorm:"updated"`
}

type UpdatedUpdate3 struct {
	Id      int64
	Updated int `xorm:"updated bigint"`
}

type UpdatedUpdate4 struct {
	Id      int64
	Updated int `xorm:"updated"`
}

type UpdatedUpdate5 struct {
	Id      int64
	Updated time.Time `xorm:"updated bigint"`
}

func testUpdateUpdated(engine *xorm.Engine, t *testing.T) {
	di := new(UpdatedUpdate)
	err := engine.Sync2(di)
	if err != nil {
		t.Fatal(err)
	}

	_, err = engine.Insert(&UpdatedUpdate{})
	if err != nil {
		t.Fatal(err)
	}

	ci := &UpdatedUpdate{}
	_, err = engine.Id(1).Update(ci)
	if err != nil {
		t.Fatal(err)
	}

	has, err := engine.Id(1).Get(di)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci.Updated.Unix() != di.Updated.Unix() {
		t.Fatal("should equal:", ci, di)
	}
	fmt.Println("ci:", ci, "di:", di)

	di2 := new(UpdatedUpdate2)
	err = engine.Sync2(di2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = engine.Insert(&UpdatedUpdate2{})
	if err != nil {
		t.Fatal(err)
	}
	ci2 := &UpdatedUpdate2{}
	_, err = engine.Id(1).Update(ci2)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Id(1).Get(di2)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci2.Updated != di2.Updated {
		t.Fatal("should equal:", ci2, di2)
	}
	fmt.Println("ci2:", ci2, "di2:", di2)

	di3 := new(UpdatedUpdate3)
	err = engine.Sync2(di3)
	if err != nil {
		t.Fatal(err)
	}

	_, err = engine.Insert(&UpdatedUpdate3{})
	if err != nil {
		t.Fatal(err)
	}
	ci3 := &UpdatedUpdate3{}
	_, err = engine.Id(1).Update(ci3)
	if err != nil {
		t.Fatal(err)
	}

	has, err = engine.Id(1).Get(di3)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci3.Updated != di3.Updated {
		t.Fatal("should equal:", ci3, di3)
	}
	fmt.Println("ci3:", ci3, "di3:", di3)

	di4 := new(UpdatedUpdate4)
	err = engine.Sync2(di4)
	if err != nil {
		t.Fatal(err)
	}

	_, err = engine.Insert(&UpdatedUpdate4{})
	if err != nil {
		t.Fatal(err)
	}

	ci4 := &UpdatedUpdate4{}
	_, err = engine.Id(1).Update(ci4)
	if err != nil {
		t.Fatal(err)
	}

	has, err = engine.Id(1).Get(di4)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci4.Updated != di4.Updated {
		t.Fatal("should equal:", ci4, di4)
	}
	fmt.Println("ci4:", ci4, "di4:", di4)

	di5 := new(UpdatedUpdate5)
	err = engine.Sync2(di5)
	if err != nil {
		t.Fatal(err)
	}

	_, err = engine.Insert(&UpdatedUpdate5{})
	if err != nil {
		t.Fatal(err)
	}
	ci5 := &UpdatedUpdate5{}
	_, err = engine.Id(1).Update(ci5)
	if err != nil {
		t.Fatal(err)
	}

	has, err = engine.Id(1).Get(di5)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci5.Updated.Unix() != di5.Updated.Unix() {
		t.Fatal("should equal:", ci5, di5)
	}
	fmt.Println("ci5:", ci5, "di5:", di5)
}

func updateSameMapper(engine *xorm.Engine, t *testing.T) {
	var ori Userinfo
	has, err := engine.Get(&ori)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		t.Error(errors.New("not exist"))
		panic(errors.New("not exist"))
	}
	// update by id
	user := Userinfo{Username: "xxx", Height: 1.2}
	cnt, err := engine.Id(ori.Uid).Update(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	condi := Condi{"Username": "zzz", "Departname": ""}
	cnt, err = engine.Table(&user).Id(ori.Uid).Update(&condi)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	cnt, err = engine.Update(&Userinfo{Username: "yyy"}, &user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	total, err := engine.Count(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != total {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	err = engine.Sync(&Article{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	defer func() {
		err = engine.DropTables(&Article{})
		if err != nil {
			t.Error(err)
			panic(err)
		}
	}()

	a := &Article{0, "1", "2", "3", "4", "5", 2}
	cnt, err = engine.Insert(a)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New(fmt.Sprintf("insert not returned 1 but %d", cnt))
		t.Error(err)
		panic(err)
	}

	if a.Id == 0 {
		err = errors.New("insert returned id is 0")
		t.Error(err)
		panic(err)
	}

	cnt, err = engine.Id(a.Id).Update(&Article{Name: "6"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New(fmt.Sprintf("insert not returned 1 but %d", cnt))
		t.Error(err)
		panic(err)
		return
	}

	col1 := &UpdateAllCols{}
	err = engine.Sync(col1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(col1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	col2 := &UpdateAllCols{col1.Id, true, "", nil}
	_, err = engine.Id(col2.Id).AllCols().Update(col2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	col3 := &UpdateAllCols{}
	has, err = engine.Id(col2.Id).Get(col3)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		err = errors.New(fmt.Sprintf("cannot get id %d", col2.Id))
		t.Error(err)
		panic(err)
		return
	}

	if *col2 != *col3 {
		err = errors.New(fmt.Sprintf("col2 should eq col3"))
		t.Error(err)
		panic(err)
		return
	}

	{
		col1 := &UpdateMustCols{}
		err = engine.Sync(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		_, err = engine.Insert(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		col2 := &UpdateMustCols{col1.Id, true, ""}
		boolStr := engine.ColumnMapper.Obj2Table("Bool")
		stringStr := engine.ColumnMapper.Obj2Table("String")
		_, err = engine.Id(col2.Id).MustCols(boolStr, stringStr).Update(col2)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		col3 := &UpdateMustCols{}
		has, err := engine.Id(col2.Id).Get(col3)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		if !has {
			err = errors.New(fmt.Sprintf("cannot get id %d", col2.Id))
			t.Error(err)
			panic(err)
			return
		}

		if *col2 != *col3 {
			err = errors.New(fmt.Sprintf("col2 should eq col3"))
			t.Error(err)
			panic(err)
			return
		}
	}

	{

		col1 := &UpdateIncr{}
		err = engine.Sync(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		_, err = engine.Insert(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		cnt, err := engine.Id(col1.Id).Incr("`Cnt`").Update(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			err = errors.New("update incr failed")
			t.Error(err)
			panic(err)
		}

		newCol := new(UpdateIncr)
		has, err := engine.Id(col1.Id).Get(newCol)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if !has {
			err = errors.New("has incr failed")
			t.Error(err)
			panic(err)
		}
		if 1 != newCol.Cnt {
			err = errors.New("incr failed")
			t.Error(err)
			panic(err)
		}
	}
}

func testUseBool(engine *xorm.Engine, t *testing.T) {
	cnt1, err := engine.Count(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	users := make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	var fNumber int64
	for _, u := range users {
		if u.IsMan == false {
			fNumber += 1
		}
	}

	cnt2, err := engine.UseBool().Update(&Userinfo{IsMan: true})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if fNumber != cnt2 {
		fmt.Println("cnt1", cnt1, "fNumber", fNumber, "cnt2", cnt2)
		/*err = errors.New("Updated number is not corrected.")
		  t.Error(err)
		  panic(err)*/
	}

	_, err = engine.Update(&Userinfo{IsMan: true})
	if err == nil {
		err = errors.New("error condition")
		t.Error(err)
		panic(err)
	}
}

func testBool(engine *xorm.Engine, t *testing.T) {
	_, err := engine.UseBool().Update(&Userinfo{IsMan: true})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	users := make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		if !user.IsMan {
			err = errors.New("update bool or find bool error")
			t.Error(err)
			panic(err)
		}
	}

	_, err = engine.UseBool().Update(&Userinfo{IsMan: false})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	users = make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		if user.IsMan {
			err = errors.New("update bool or find bool error")
			t.Error(err)
			panic(err)
		}
	}
}
