package tests

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-xorm/xorm"
)

func table(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables("user_user")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Table("user_user").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func createMultiTables(engine *xorm.Engine, t *testing.T) {
	session := engine.NewSession()
	defer session.Close()

	user := &Userinfo{}
	err := session.Begin()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for i := 0; i < 10; i++ {
		tableName := fmt.Sprintf("user_%v", i)

		err = session.DropTable(tableName)
		if err != nil {
			session.Rollback()
			t.Error(err)
			panic(err)
		}

		err = session.Table(tableName).CreateTable(user)
		if err != nil {
			session.Rollback()
			t.Error(err)
			panic(err)
		}
	}
	err = session.Commit()
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func directCreateTable(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&Userinfo{}, &Userdetail{}, &Numeric{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Sync(&Userinfo{}, &Userdetail{}, new(Picture), new(Numeric))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	isEmpty, err := engine.IsTableEmpty(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !isEmpty {
		err = errors.New("userinfo should empty")
		t.Error(err)
		panic(err)
	}

	tbName := engine.TableMapper.Obj2Table("Userinfo")
	isEmpty, err = engine.IsTableEmpty(tbName)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !isEmpty {
		err = errors.New("userinfo should empty")
		t.Error(err)
		panic(err)
	}

	err = engine.DropTables(&Userinfo{}, &Userdetail{}, new(Numeric))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Userinfo{}, &Userdetail{}, new(Numeric))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateIndexes(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateIndexes(&Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateUniques(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateUniques(&Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

type CustomTableName struct {
	Id   int64
	Name string
}

func (c *CustomTableName) TableName() string {
	return "customtablename"
}

func testCustomTableName(engine *xorm.Engine, t *testing.T) {
	c := new(CustomTableName)
	err := engine.DropTables(c)
	if err != nil {
		t.Error(err)
	}

	err = engine.CreateTables(c)
	if err != nil {
		t.Error(err)
	}
}

func testDump(engine *xorm.Engine, t *testing.T) {
	fp := engine.Dialect().URI().DbName + ".sql"
	os.Remove(fp)
	err := engine.DumpAllToFile(fp)
	if err != nil {
		t.Error(err)
		fmt.Println(err)
	}
}

type IndexOrUnique struct {
	Id        int64
	Index     int `xorm:"index"`
	Unique    int `xorm:"unique"`
	Group1    int `xorm:"index(ttt)"`
	Group2    int `xorm:"index(ttt)"`
	UniGroup1 int `xorm:"unique(lll)"`
	UniGroup2 int `xorm:"unique(lll)"`
}

func testIndexAndUnique(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.DropTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.CreateTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.CreateIndexes(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateUniques(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}
}

func testMetaInfo(engine *xorm.Engine, t *testing.T) {
	tables, err := engine.DBMetas()
	if err != nil {
		t.Error(err)
		panic(err)
	}

	for _, table := range tables {
		fmt.Println(table.Name)
		//for _, col := range table.Columns() {
		//TODO: engine.dialect show exported
		//fmt.Println(col.String(engine.dia))
		//}

		for _, index := range table.Indexes {
			fmt.Println(index.Name, index.Type, strings.Join(index.Cols, ","))
		}
	}
}
