package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
)

func insert(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{0, "xiaolunwen", "dev", "lunny", time.Now(),
		Userdetail{Id: 1}, 1.78, []byte{1, 2, 3}, true}
	cnt, err := engine.Insert(&user)
	fmt.Println(user.Uid)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
	if user.Uid <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	user.Uid = 0
	cnt, err = engine.Insert(&user)
	if err == nil {
		err = errors.New("insert failed but no return error")
		t.Error(err)
		panic(err)
	}
	if cnt != 0 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	testInsertCreated(engine, t)
	testCreatedJsonTime(engine, t)
}

func insertAutoIncr(engine *xorm.Engine, t *testing.T) {
	// auto increment insert
	user := Userinfo{Username: "xiaolunwen2", Departname: "dev", Alias: "lunny", Created: time.Now(),
		Detail: Userdetail{Id: 1}, Height: 1.78, Avatar: []byte{1, 2, 3}, IsMan: true}
	cnt, err := engine.Insert(&user)
	fmt.Println(user.Uid)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
	if user.Uid <= 0 {
		t.Error(errors.New("not return id error"))
	}
}

type DefaultInsert struct {
	Id      int64
	Status  int `xorm:"default -1"`
	Name    string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func testInsertDefault(engine *xorm.Engine, t *testing.T) {
	di := new(DefaultInsert)
	err := engine.Sync2(di)
	if err != nil {
		t.Error(err)
	}

	var di2 = DefaultInsert{Name: "test"}
	_, err = engine.Omit(engine.ColumnMapper.Obj2Table("Status")).Insert(&di2)
	if err != nil {
		t.Error(err)
	}

	has, err := engine.Desc("(id)").Get(di)
	if err != nil {
		t.Error(err)
	}
	if !has {
		err = errors.New("error with no data")
		t.Error(err)
		panic(err)
	}
	if di.Status != -1 {
		err = errors.New("inserted error data")
		t.Error(err)
		panic(err)
	}
	if di2.Updated.Unix() != di.Updated.Unix() {
		err = errors.New("updated should equal")
		t.Error(err, di.Updated, di2.Updated)
		panic(err)
	}
	if di2.Created.Unix() != di.Created.Unix() {
		err = errors.New("created should equal")
		t.Error(err, di.Created, di2.Created)
		panic(err)
	}

	testInsertDefault2(engine, t)
}

type DefaultInsert2 struct {
	Id        int64
	Name      string
	Url       string    `xorm:"text"`
	CheckTime time.Time `xorm:"not null default '2000-01-01 00:00:00' TIMESTAMP"`
}

func testInsertDefault2(engine *xorm.Engine, t *testing.T) {
	di := new(DefaultInsert2)
	err := engine.Sync2(di)
	if err != nil {
		t.Error(err)
	}

	var di2 = DefaultInsert2{Name: "test"}
	_, err = engine.Omit(engine.ColumnMapper.Obj2Table("CheckTime")).Insert(&di2)
	if err != nil {
		t.Error(err)
	}

	has, err := engine.Desc("(id)").Get(di)
	if err != nil {
		t.Error(err)
	}
	if !has {
		err = errors.New("error with no data")
		t.Error(err)
		panic(err)
	}

	has, err = engine.NoAutoCondition().Desc("(id)").Get(&di2)
	if err != nil {
		t.Error(err)
	}

	if !has {
		err = errors.New("error with no data")
		t.Error(err)
		panic(err)
	}

	if *di != di2 {
		err = fmt.Errorf("%v is not equal to %v", di, di2)
		t.Error(err)
		panic(err)
	}

	/*if di2.Updated.Unix() != di.Updated.Unix() {
		err = errors.New("updated should equal")
		t.Error(err, di.Updated, di2.Updated)
		panic(err)
	}
	if di2.Created.Unix() != di.Created.Unix() {
		err = errors.New("created should equal")
		t.Error(err, di.Created, di2.Created)
		panic(err)
	}*/
}

type CreatedInsert struct {
	Id      int64
	Created time.Time `xorm:"created"`
}

type CreatedInsert2 struct {
	Id      int64
	Created int64 `xorm:"created"`
}

type CreatedInsert3 struct {
	Id      int64
	Created int `xorm:"created bigint"`
}

type CreatedInsert4 struct {
	Id      int64
	Created int `xorm:"created"`
}

type CreatedInsert5 struct {
	Id      int64
	Created time.Time `xorm:"created bigint"`
}

type CreatedInsert6 struct {
	Id      int64
	Created time.Time `xorm:"created bigint"`
}

func testInsertCreated(engine *xorm.Engine, t *testing.T) {
	di := new(CreatedInsert)
	err := engine.Sync2(di)
	if err != nil {
		t.Fatal(err)
	}
	ci := &CreatedInsert{}
	_, err = engine.Insert(ci)
	if err != nil {
		t.Fatal(err)
	}

	has, err := engine.Desc("(id)").Get(di)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci.Created.Unix() != di.Created.Unix() {
		t.Fatal("should equal:", ci, di)
	}
	fmt.Println("ci:", ci, "di:", di)

	di2 := new(CreatedInsert2)
	err = engine.Sync2(di2)
	if err != nil {
		t.Fatal(err)
	}
	ci2 := &CreatedInsert2{}
	_, err = engine.Insert(ci2)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Desc("(id)").Get(di2)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci2.Created != di2.Created {
		t.Fatal("should equal:", ci2, di2)
	}
	fmt.Println("ci2:", ci2, "di2:", di2)

	di3 := new(CreatedInsert3)
	err = engine.Sync2(di3)
	if err != nil {
		t.Fatal(err)
	}
	ci3 := &CreatedInsert3{}
	_, err = engine.Insert(ci3)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Desc("(id)").Get(di3)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci3.Created != di3.Created {
		t.Fatal("should equal:", ci3, di3)
	}
	fmt.Println("ci3:", ci3, "di3:", di3)

	di4 := new(CreatedInsert4)
	err = engine.Sync2(di4)
	if err != nil {
		t.Fatal(err)
	}
	ci4 := &CreatedInsert4{}
	_, err = engine.Insert(ci4)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Desc("(id)").Get(di4)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci4.Created != di4.Created {
		t.Fatal("should equal:", ci4, di4)
	}
	fmt.Println("ci4:", ci4, "di4:", di4)

	di5 := new(CreatedInsert5)
	err = engine.Sync2(di5)
	if err != nil {
		t.Fatal(err)
	}
	ci5 := &CreatedInsert5{}
	_, err = engine.Insert(ci5)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Desc("(id)").Get(di5)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci5.Created.Unix() != di5.Created.Unix() {
		t.Fatal("should equal:", ci5, di5)
	}
	fmt.Println("ci5:", ci5, "di5:", di5)

	di6 := new(CreatedInsert6)
	err = engine.Sync2(di6)
	if err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().Add(-time.Hour)
	ci6 := &CreatedInsert6{Created: oldTime}
	_, err = engine.Insert(ci6)
	if err != nil {
		t.Fatal(err)
	}

	has, err = engine.Desc("(id)").Get(di6)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci6.Created.Unix() != di6.Created.Unix() {
		t.Fatal("should equal:", ci6, di6)
	}
	fmt.Println("ci6:", ci6, "di6:", di6)
}

type JsonTime time.Time

func (j JsonTime) format() string {
	t := time.Time(j)
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02")
}

func (j JsonTime) MarshalText() ([]byte, error) {
	return []byte(j.format()), nil
}

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

type MyJsonTime struct {
	Id      int64    `json:"id"`
	Created JsonTime `xorm:"created" json:"created_at"`
}

func testCreatedJsonTime(engine *xorm.Engine, t *testing.T) {
	di5 := new(MyJsonTime)
	err := engine.Sync2(di5)
	if err != nil {
		t.Fatal(err)
	}
	ci5 := &MyJsonTime{}
	_, err = engine.Insert(ci5)
	if err != nil {
		t.Fatal(err)
	}
	has, err := engine.Desc("(id)").Get(di5)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if time.Time(ci5.Created).Unix() != time.Time(di5.Created).Unix() {
		t.Fatal("should equal:", time.Time(ci5.Created).Unix(), time.Time(di5.Created).Unix())
	}
	fmt.Println("ci5:", ci5, "di5:", di5)

	var dis = make([]MyJsonTime, 0)
	err = engine.Find(&dis)
	if err != nil {
		t.Fatal(err)
	}
}

func insertMulti(engine *xorm.Engine, t *testing.T) {
	users := []Userinfo{
		{Username: "xlw", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		{Username: "xlw2", Departname: "dev", Alias: "lunny3", Created: time.Now()},
		{Username: "xlw11", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		{Username: "xlw22", Departname: "dev", Alias: "lunny3", Created: time.Now()},
	}
	cnt, err := engine.Insert(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != int64(len(users)) {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	users2 := []*Userinfo{
		&Userinfo{Username: "1xlw", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		&Userinfo{Username: "1xlw2", Departname: "dev", Alias: "lunny3", Created: time.Now()},
		&Userinfo{Username: "1xlw11", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		&Userinfo{Username: "1xlw22", Departname: "dev", Alias: "lunny3", Created: time.Now()},
	}

	cnt, err = engine.Insert(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != int64(len(users2)) {
		err = errors.New(fmt.Sprintf("insert not returned %v", len(users2)))
		t.Error(err)
		panic(err)
		return
	}
}

func insertTwoTable(engine *xorm.Engine, t *testing.T) {
	userdetail := Userdetail{ /*Id: 1, */ Intro: "I'm a very beautiful women.", Profile: "sfsaf"}
	userinfo := Userinfo{Username: "xlw3", Departname: "dev", Alias: "lunny4", Created: time.Now(), Detail: userdetail}

	cnt, err := engine.Insert(&userinfo, &userdetail)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if userinfo.Uid <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	if userdetail.Id <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	if cnt != 2 {
		err = errors.New("insert not returned 2")
		t.Error(err)
		panic(err)
		return
	}
}
