package tests

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	. "github.com/go-xorm/tests"
	"github.com/go-xorm/xorm"
)

func connStr(db string) string {
	conn := "root:@"
	if ConnectionPort != "" {
		conn += "tcp(127.0.0.1:" + ConnectionPort + ")"
	}
	return conn + "/" + db + "?charset=utf8"
}

/*
CREATE DATABASE IF NOT EXISTS xorm_test CHARACTER SET
utf8 COLLATE utf8_general_ci;
*/
func TestMysql(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", connStr("xorm_test"))
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.ShowSQL(ShowTestSql)

	BaseTestAll(engine, t)
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestMysqlSameMapper(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", connStr("xorm_test1"))
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.ShowSQL(ShowTestSql)
	engine.SetMapper(core.SameMapper{})

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestMysqlGonicMapper(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", connStr("xorm_test1"))
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.ShowSQL(ShowTestSql)
	engine.SetMapper(core.GonicMapper{})

	BaseTestAll(engine, t)
	//BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestMysqlWithCache(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", connStr("xorm_test2"))
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL(ShowTestSql)

	BaseTestAll(engine, t)
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
}

func TestMysqlWithCacheSameMapper(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", connStr("xorm_test3"))
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetMapper(core.SameMapper{})
	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL(ShowTestSql)

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
}

func newMysqlEngine() (*xorm.Engine, error) {
	return xorm.NewEngine("mysql", connStr("xorm_test"))
}

func mysqlDdlImport() error {
	engine, err := xorm.NewEngine("mysql", connStr(""))
	if err != nil {
		return err
	}
	engine.ShowSQL(ShowTestSql)

	sqlResults, err := engine.ImportFile("../testdata/mysql_ddl.sql")
	if err != nil {
		return err
	}
	engine.Logger().Debug("sql results:", sqlResults)
	engine.Close()
	return nil
}

func newMysqlDriverDB() (*sql.DB, error) {
	return sql.Open("mysql", connStr("xorm_test"))
}

func BenchmarkMysqlDriverInsert(t *testing.B) {
	DoBenchDriver(newMysqlDriverDB, CreateTableMySql, DropTableMySql,
		DoBenchDriverInsert, t)
}

func BenchmarkMysqlDriverFind(t *testing.B) {
	DoBenchDriver(newMysqlDriverDB, CreateTableMySql, DropTableMySql,
		DoBenchDriverFind, t)
}

func BenchmarkMysqlNoCacheInsert(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.ShowSQL(true)
	DoBenchInsert(engine, t)
}

func BenchmarkMysqlNoCacheFind(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.ShowSQL(true)
	DoBenchFind(engine, t)
}

func BenchmarkMysqlNoCacheFindPtr(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.ShowSQL(true)
	DoBenchFindPtr(engine, t)
}

func BenchmarkMysqlCacheInsert(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL(true)

	DoBenchInsert(engine, t)
}

func BenchmarkMysqlCacheFind(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL(true)

	DoBenchFind(engine, t)
}

func BenchmarkMysqlCacheFindPtr(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL(true)

	DoBenchFindPtr(engine, t)
}
