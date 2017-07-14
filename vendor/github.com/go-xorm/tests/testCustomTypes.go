package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

func testCustomType(engine *xorm.Engine, t *testing.T) {
	testCustomType1(engine, t)
	testCustomType2(engine, t)
}

type ConvString string

func (s *ConvString) FromDB(data []byte) error {
	*s = ConvString("prefix---" + string(data))
	return nil
}

func (s *ConvString) ToDB() ([]byte, error) {
	return []byte(string(*s)), nil
}

type ConvConfig struct {
	Name string
	Id   int64
}

func (s *ConvConfig) FromDB(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *ConvConfig) ToDB() ([]byte, error) {
	return json.Marshal(s)
}

type SliceType []*ConvConfig

func (s *SliceType) FromDB(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *SliceType) ToDB() ([]byte, error) {
	return json.Marshal(s)
}

type ConvStruct struct {
	Conv  ConvString
	Conv2 *ConvString
	Cfg1  ConvConfig
	Cfg2  *ConvConfig     `xorm:"TEXT"`
	Cfg3  core.Conversion `xorm:"BLOB"`
	Slice SliceType
}

func (c *ConvStruct) BeforeSet(name string, cell xorm.Cell) {
	if name == "cfg3" || name == "Cfg3" {
		c.Cfg3 = new(ConvConfig)
	}
}

func testConversion(engine *xorm.Engine, t *testing.T) {
	c := new(ConvStruct)
	err := engine.DropTables(c)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	err = engine.Sync(c)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var s ConvString = "sssss"
	c.Conv = "tttt"
	c.Conv2 = &s
	c.Cfg1 = ConvConfig{"mm", 1}
	c.Cfg2 = &ConvConfig{"xx", 2}
	c.Cfg3 = &ConvConfig{"zz", 3}
	c.Slice = []*ConvConfig{{"yy", 4}, {"ff", 5}}

	_, err = engine.Insert(c)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	c1 := new(ConvStruct)
	_, err = engine.Get(c1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if string(c1.Conv) != "prefix---tttt" {
		err = fmt.Errorf("get conversion error prefix---tttt != %s", c1.Conv)
		t.Error(err)
		panic(err)
	}

	if c1.Conv2 == nil || *c1.Conv2 != "prefix---"+s {
		err = fmt.Errorf("get conversion error2, %v", *c1.Conv2)
		t.Error(err)
		panic(err)
	}

	if c1.Cfg1 != c.Cfg1 {
		err = fmt.Errorf("get conversion error3, %v", c1.Cfg1)
		t.Error(err)
		panic(err)
	}

	if c1.Cfg2 == nil || *c1.Cfg2 != *c.Cfg2 {
		err = fmt.Errorf("get conversion error4, %v", *c1.Cfg2)
		t.Error(err)
		panic(err)
	}

	if c1.Cfg3 == nil || *c1.Cfg3.(*ConvConfig) != *c.Cfg3.(*ConvConfig) {
		err = fmt.Errorf("get conversion error5, %v", *c1.Cfg3.(*ConvConfig))
		t.Error(err)
		panic(err)
	}

	if len(c1.Slice) != 2 {
		err = fmt.Errorf("get conversion error6, should be 2")
		t.Error(err)
		panic(err)
	}

	if *c1.Slice[0] != *c.Slice[0] ||
		*c1.Slice[1] != *c.Slice[1] {
		err = fmt.Errorf("get conversion error7, should be %v", c1.Slice)
		t.Error(err)
		panic(err)
	}
}

type MyInt int
type MyUInt uint
type MyFloat float64
type MyString string

type MyStruct struct {
	Type      MyInt
	U         MyUInt
	F         MyFloat
	S         MyString
	IA        []MyInt
	UA        []MyUInt
	FA        []MyFloat
	SA        []MyString
	NameArray []string
	Name      string
	UIA       []uint
	UIA8      []uint8
	UIA16     []uint16
	UIA32     []uint32
	UIA64     []uint64
	UI        uint
	//C64       complex64
	MSS map[string]string
}

func testCustomType1(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&MyStruct{})
	if err != nil {
		t.Error(err)
		panic(err)
		return
	}

	err = engine.CreateTables(&MyStruct{})
	i := MyStruct{Name: "Test", Type: MyInt(1)}
	i.U = 23
	i.F = 1.34
	i.S = "fafdsafdsaf"
	i.UI = 2
	i.IA = []MyInt{1, 3, 5}
	i.UIA = []uint{1, 3}
	i.UIA16 = []uint16{2}
	i.UIA32 = []uint32{4, 5}
	i.UIA64 = []uint64{6, 7, 9}
	i.UIA8 = []uint8{1, 2, 3, 4}
	i.NameArray = []string{"ssss", "fsdf", "lllll, ss"}
	i.MSS = map[string]string{"s": "sfds,ss", "x": "lfjljsl"}
	cnt, err := engine.Insert(&i)
	if err != nil {
		t.Error(err)
		panic(err)
		return
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	fmt.Println(i)
	i.NameArray = []string{}
	i.MSS = map[string]string{}
	i.F = 0
	has, err := engine.Get(&i)
	if err != nil {
		t.Error(err)
		panic(err)
	} else if !has {
		t.Error(errors.New("should get one record"))
		panic(err)
	}

	ss := []MyStruct{}
	err = engine.Find(&ss)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(ss)

	sss := MyStruct{}
	has, err = engine.Get(&sss)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(sss)

	if has {
		sss.NameArray = []string{}
		sss.MSS = map[string]string{}
		cnt, err := engine.Delete(&sss)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			t.Error(errors.New("delete error"))
			panic(err)
		}
	}
}

type Status struct {
	Name  string
	Color string
}

var (
	_        core.Conversion   = &Status{}
	Registed Status            = Status{"Registed", "white"}
	Approved Status            = Status{"Approved", "green"}
	Removed  Status            = Status{"Removed", "red"}
	Statuses map[string]Status = map[string]Status{
		Registed.Name: Registed,
		Approved.Name: Approved,
		Removed.Name:  Removed,
	}
)

func (s *Status) FromDB(bytes []byte) error {
	if r, ok := Statuses[string(bytes)]; ok {
		*s = r
		return nil
	} else {
		return errors.New("no this data")
	}
}

func (s *Status) ToDB() ([]byte, error) {
	return []byte(s.Name), nil
}

type UserCus struct {
	Id     int64
	Name   string
	Status Status `xorm:"varchar(40)"`
}

func testCustomType2(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&UserCus{})
	if err != nil {
		t.Fatal(err)
	}

	tableName := engine.TableMapper.Obj2Table("UserCus")
	_, err = engine.Exec("delete from " + engine.Quote(tableName))
	if err != nil {
		t.Fatal(err)
	}

	if engine.Dialect().DBType() == core.MSSQL {
		return
		_, err = engine.Exec("set IDENTITY_INSERT " + tableName + " on")
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = engine.Insert(&UserCus{1, "xlw", Registed})
	if err != nil {
		t.Fatal(err)
	}

	user := UserCus{}
	exist, err := engine.Id(1).Get(&user)
	if err != nil {
		t.Fatal(err)
	}

	if !exist {
		t.Fatal("user not exist")
	}

	fmt.Println(user)

	users := make([]UserCus, 0)
	err = engine.Where("`"+engine.ColumnMapper.Obj2Table("Status")+"` = ?", "Registed").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users) != 1 {
		t.Error("users should has 1 record.")
		panic("")
	}

	fmt.Println(users)
}
