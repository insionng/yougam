package tests

import (
	"database/sql"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	. "github.com/go-xorm/tests"
	"github.com/go-xorm/xorm"
)

var mssqlConnStr = "server=localhost;user id=sa;password=;database=xorm_test"

func newMssqlEngine() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mssql", mssqlConnStr)
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(ShowTestSql)
	return engine, nil
}

func TestMssql(t *testing.T) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	BaseTestAll(engine, t)
	BaseTestAll2(engine, t)
}

func TestMssqlWithCache(t *testing.T) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	BaseTestAll(engine, t)
	BaseTestAll2(engine, t)
}

func newMssqlDriverDB() (*sql.DB, error) {
	return sql.Open("mssql", mssqlConnStr)
}

const (
	createTableMssql = `IF NOT EXISTS (SELECT [name] FROM sys.tables WHERE [name] = 'big_struct' ) CREATE TABLE
		"big_struct" ("id" BIGINT PRIMARY KEY IDENTITY NOT NULL, "name" VARCHAR(255) NULL, "title" VARCHAR(255) NULL,
		"age" VARCHAR(255) NULL, "alias" VARCHAR(255) NULL, "nick_name" VARCHAR(255) NULL);
		`

	dropTableMssql = "IF EXISTS (SELECT * FROM sysobjects WHERE id = object_id(N'big_struct') and OBJECTPROPERTY(id, N'IsUserTable') = 1) DROP TABLE IF EXISTS `big_struct`;"
)

func BenchmarkMssqlDriverInsert(t *testing.B) {
	DoBenchDriver(newMssqlDriverDB, createTableMssql, dropTableMssql,
		DoBenchDriverInsert, t)
}

func BenchmarkMssqlDriverFind(t *testing.B) {
	DoBenchDriver(newMssqlDriverDB, createTableMssql, dropTableMssql,
		DoBenchDriverFind, t)
}

func BenchmarkMssqlNoCacheInsert(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	DoBenchInsert(engine, t)
}

func BenchmarkMssqlNoCacheFind(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	DoBenchFind(engine, t)
}

func BenchmarkMssqlNoCacheFindPtr(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	DoBenchFindPtr(engine, t)
}

func BenchmarkMssqlCacheInsert(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	DoBenchInsert(engine, t)
}

func BenchmarkMssqlCacheFind(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	DoBenchFind(engine, t)
}

func BenchmarkMssqlCacheFindPtr(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	DoBenchFindPtr(engine, t)
}
