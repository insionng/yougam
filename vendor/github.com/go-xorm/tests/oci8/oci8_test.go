package tests

import (
	"database/sql"
	"testing"

	. ".."
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-oci8"
)

var connStr string = "anonymous/123456@192.168.59.103:49161/xe"

func newOci8Engine() (*xorm.Engine, error) {
	orm, err := xorm.NewEngine("oci8", connStr)
	if err != nil {
		return nil, err
	}
	orm.ShowSQL(ShowTestSql)

	tables, err := orm.DBMetas()
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		_, err = orm.Exec("drop table \"" + table.Name + "\"")
		if err != nil {
			return nil, err
		}
	}

	return orm, err
}

func newOci8DriverDB() (*sql.DB, error) {
	return sql.Open("oci8", connStr)
}

func TestOci8(t *testing.T) {
	engine, err := newOci8Engine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	BaseTestAll(engine, t)
	UserTest1(engine, t)
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestOci8WithCache(t *testing.T) {
	engine, err := newOci8Engine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	BaseTestAll(engine, t)
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
}

func TestOci8SameMapper(t *testing.T) {
	engine, err := newOci8Engine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetMapper(core.SameMapper{})

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestOci8WithCacheSameMapper(t *testing.T) {
	engine, err := newOci8Engine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())
	engine.SetMapper(core.SameMapper{})

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
}

const (
	createTableOci8 = `CREATE TABLE IF NOT EXISTS "big_struct" ("id" SERIAL PRIMARY KEY  NOT NULL, "name" VARCHAR(255) NULL, "title" VARCHAR(255) NULL, "age" VARCHAR(255) NULL, "alias" VARCHAR(255) NULL, "nick_name" VARCHAR(255) NULL);`
	dropTableOci8   = `DROP TABLE IF EXISTS "big_struct";`
)

func BenchmarkOci8DriverInsert(t *testing.B) {
	DoBenchDriver(newOci8DriverDB, createTableOci8, dropTableOci8,
		DoBenchDriverInsert, t)
}

func BenchmarkOci8DriverFind(t *testing.B) {
	DoBenchDriver(newOci8DriverDB, createTableOci8, dropTableOci8,
		DoBenchDriverFind, t)
}

func BenchmarkOci8NoCacheInsert(t *testing.B) {
	engine, err := newOci8Engine()

	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchInsert(engine, t)
}

func BenchmarkOci8NoCacheFind(t *testing.B) {
	engine, err := newOci8Engine()

	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchFind(engine, t)
}

func BenchmarkOci8NoCacheFindPtr(t *testing.B) {
	engine, err := newOci8Engine()

	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchFindPtr(engine, t)
}

func BenchmarkOci8CacheInsert(t *testing.B) {
	engine, err := newOci8Engine()

	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchInsert(engine, t)
}

func BenchmarkOci8CacheFind(t *testing.B) {
	engine, err := newOci8Engine()

	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchFind(engine, t)
}

func BenchmarkOci8CacheFindPtr(t *testing.B) {
	engine, err := newOci8Engine()

	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchFindPtr(engine, t)
}
