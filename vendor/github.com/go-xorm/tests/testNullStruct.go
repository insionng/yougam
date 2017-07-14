package tests

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"strconv"
	"strings"
	"testing"
)

type NullType struct {
	Id           int `xorm:"pk autoincr"`
	Name         sql.NullString
	Age          sql.NullInt64
	Height       sql.NullFloat64
	IsMan        sql.NullBool `xorm:"null"`
	CustomStruct CustomStruct `xorm:"valchar(64) null"`
}

type CustomStruct struct {
	Year  int
	Month int
	Day   int
}

func (CustomStruct) String() string {
	return "CustomStruct"
}

func (m *CustomStruct) Scan(value interface{}) error {
	if value == nil {
		m.Year, m.Month, m.Day = 0, 0, 0
		return nil
	}

	if s, ok := value.([]byte); ok {
		seps := strings.Split(string(s), "/")
		m.Year, _ = strconv.Atoi(seps[0])
		m.Month, _ = strconv.Atoi(seps[1])
		m.Day, _ = strconv.Atoi(seps[2])
	}

	return errors.New("sacn dat not fit []byte")
}

func (m CustomStruct) Value() (driver.Value, error) {
	return fmt.Sprintf("%d/%d/%d", m.Year, m.Month, m.Day), nil
}

func TestNullStruct(engine *xorm.Engine, t *testing.T) {
	TestDropNullStructTable(engine, t)
	TestCreateNullStructTable(engine, t)

	TestNullStructInsert(engine, t)
	TestNullStructUpdate(engine, t)
	TestNullStructFind(engine, t)
	TestNullStructIterate(engine, t)
	TestNullStructCount(engine, t)
	TestNullStructRows(engine, t)
	TestNullStructDelete(engine, t)

	TestDropNullStructTable(engine, t)
}

func TestCreateNullStructTable(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(new(NullType))
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func TestDropNullStructTable(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(new(NullType))
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func TestNullStructInsert(engine *xorm.Engine, t *testing.T) {
	if true {
		item := new(NullType)
		_, err := engine.Insert(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		fmt.Println(item)
		if item.Id != 1 {
			err = errors.New("insert error")
			t.Error(err)
			panic(err)
		}
	}

	if true {
		item := NullType{
			Name:   sql.NullString{"haolei", true},
			Age:    sql.NullInt64{34, true},
			Height: sql.NullFloat64{1.72, true},
			IsMan:  sql.NullBool{true, true},
		}
		_, err := engine.Insert(&item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		fmt.Println(item)
		if item.Id != 2 {
			err = errors.New("insert error")
			t.Error(err)
			panic(err)
		}
	}

	if true {
		items := []NullType{}

		for i := 0; i < 5; i++ {
			item := NullType{
				Name:         sql.NullString{"haolei_" + fmt.Sprint(i+1), true},
				Age:          sql.NullInt64{30 + int64(i), true},
				Height:       sql.NullFloat64{1.5 + 1.1*float64(i), true},
				IsMan:        sql.NullBool{true, true},
				CustomStruct: CustomStruct{i, i + 1, i + 2},
			}

			items = append(items, item)
		}

		_, err := engine.Insert(&items)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		fmt.Println(items)
	}
}

func TestNullStructUpdate(engine *xorm.Engine, t *testing.T) {
	if true { // 测试可插入NULL
		item := new(NullType)
		item.Age = sql.NullInt64{23, true}
		item.Height = sql.NullFloat64{0, false} // update to NULL

		affected, err := engine.Id(2).Cols("height", "is_man").Update(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if affected != 1 {
			err := errors.New("update failed")
			t.Error(err)
			panic(err)
		}
	}

	if true { // 测试In update
		item := new(NullType)
		affected, err := engine.In("id", 3, 4).Cols("height", "is_man").Update(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if affected != 2 {
			err := errors.New("update failed")
			t.Error(err)
			panic(err)
		}
	}

	if true { // 测试where
		item := new(NullType)
		item.Name = sql.NullString{"nullname", true}
		item.IsMan = sql.NullBool{true, true}
		item.Age = sql.NullInt64{34, true}

		_, err := engine.Where("age > ?", 34).Update(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
	}

	if true { // 修改全部时，插入空值
		item := &NullType{
			Name:   sql.NullString{"winxxp", true},
			Age:    sql.NullInt64{30, true},
			Height: sql.NullFloat64{1.72, true},
			// IsMan:  sql.NullBool{true, true},
		}

		_, err := engine.AllCols().Id(6).Update(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		fmt.Println(item)
	}

}

func TestNullStructFind(engine *xorm.Engine, t *testing.T) {
	if true {
		item := new(NullType)
		has, err := engine.Id(1).Get(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if !has {
			t.Error(errors.New("no find id 1"))
			panic(err)
		}
		fmt.Println(item)
		if item.Id != 1 || item.Name.Valid || item.Age.Valid || item.Height.Valid ||
			item.IsMan.Valid {
			err = errors.New("insert error")
			t.Error(err)
			panic(err)
		}
	}

	if true {
		item := new(NullType)
		item.Id = 2

		has, err := engine.Get(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if !has {
			t.Error(errors.New("no find id 2"))
			panic(err)
		}
		fmt.Println(item)
	}

	if true {
		item := make([]NullType, 0)

		err := engine.Id(2).Find(&item)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		fmt.Println(item)
	}

	if true {
		item := make([]NullType, 0)

		err := engine.Asc("age").Find(&item)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		for k, v := range item {
			fmt.Println(k, v)
		}
	}
}

func TestNullStructIterate(engine *xorm.Engine, t *testing.T) {
	if true {
		err := engine.Where("age IS NOT NULL").OrderBy("age").Iterate(new(NullType),
			func(i int, bean interface{}) error {
				nultype := bean.(*NullType)
				fmt.Println(i, nultype)
				return nil
			})
		if err != nil {
			t.Error(err)
			panic(err)
		}
	}
}

func TestNullStructCount(engine *xorm.Engine, t *testing.T) {
	if true {
		item := new(NullType)
		total, err := engine.Where("age IS NOT NULL").Count(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		fmt.Println(total)
	}
}

func TestNullStructRows(engine *xorm.Engine, t *testing.T) {
	item := new(NullType)
	rows, err := engine.Where("id > ?", 1).Rows(item)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(item)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		fmt.Println(item)
	}
}

func TestNullStructDelete(engine *xorm.Engine, t *testing.T) {
	item := new(NullType)

	_, err := engine.Id(1).Delete(item)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Where("id > ?", 1).Delete(item)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}
