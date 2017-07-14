package tests

import (
	"testing"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

type IDGonicMapper struct {
	ID int64
}

func testID(engine *xorm.Engine, t *testing.T) {
	testGonicMapperID(engine, t)
	testSameMapperID(engine, t)
}

func testGonicMapperID(engine *xorm.Engine, t *testing.T) {
	oldMapper := engine.ColumnMapper
	engine.SetMapper(core.LintGonicMapper)
	defer engine.SetMapper(oldMapper)

	err := engine.CreateTables(new(IDGonicMapper))
	if err != nil {
		t.Fatal(err)
	}

	tables, err := engine.DBMetas()
	if err != nil {
		t.Fatal(err)
	}

	for _, tb := range tables {
		if tb.Name == "id_gonic_mapper" {
			if len(tb.PKColumns()) != 1 && !tb.PKColumns()[0].IsPrimaryKey && !tb.PKColumns()[0].IsPrimaryKey {
				t.Fatal(tb)
			}
			return
		}
	}

	t.Fatal("not table id_gonic_mapper")
}

type IDSameMapper struct {
	ID int64
}

func testSameMapperID(engine *xorm.Engine, t *testing.T) {
	oldMapper := engine.ColumnMapper
	engine.SetMapper(core.SameMapper{})
	defer engine.SetMapper(oldMapper)

	err := engine.CreateTables(new(IDSameMapper))
	if err != nil {
		t.Fatal(err)
	}

	tables, err := engine.DBMetas()
	if err != nil {
		t.Fatal(err)
	}

	for _, tb := range tables {
		if tb.Name == "IDSameMapper" {
			if len(tb.PKColumns()) != 1 && !tb.PKColumns()[0].IsPrimaryKey && !tb.PKColumns()[0].IsPrimaryKey {
				t.Fatal(tb)
			}
			return
		}
	}
	t.Fatal("not table IDSameMapper")
}
