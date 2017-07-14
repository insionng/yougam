package tests

import (
	"errors"
	"testing"

	. "github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

func testBuilder(engine *xorm.Engine, t *testing.T) {
	testBuilder1(engine, t)
}

const (
	OpEqual int = iota
	OpGreatThan
	OpLessThan
)

type Condition struct {
	Id        int64
	TableName string
	ColName   string
	Op        int
	Value     string
}

func testBuilder1(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&Condition{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(&Condition{TableName: "table1", ColName: "col1", Op: OpEqual, Value: "1"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var cond Condition
	has, err := engine.Where(Eq{"col_name": "col1"}).Get(&cond)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		err = errors.New("records should exist")
		t.Error(err)
		panic(err)
	}

	has, err = engine.Where(Eq{"col_name": "col1"}.And(Eq{"op": OpEqual})).Get(&cond)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		err = errors.New("records should exist")
		t.Error(err)
		panic(err)
	}

	has, err = engine.Where(Eq{"col_name": "col1", "op": OpEqual, "value": "1"}).Get(&cond)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		err = errors.New("records should exist")
		t.Error(err)
		panic(err)
	}

	has, err = engine.Where(Eq{"col_name": "col1"}.And(Neq{"op": OpEqual})).Get(&cond)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if has {
		err = errors.New("records should not exist")
		t.Error(err)
		panic(err)
	}

	var conds []Condition
	err = engine.Where(Eq{"col_name": "col1"}.And(Eq{"op": OpEqual})).Find(&conds)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(conds) != 1 {
		err = errors.New("records should exist")
		t.Error(err)
		panic(err)
	}

	conds = make([]Condition, 0)
	err = engine.Where(Like{"col_name", "col"}).Find(&conds)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(conds) != 1 {
		err = errors.New("records should exist")
		t.Error(err)
		panic(err)
	}

	conds = make([]Condition, 0)
	err = engine.Where(Expr("col_name = ?", "col1")).Find(&conds)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(conds) != 1 {
		err = errors.New("records should exist")
		t.Error(err)
		panic(err)
	}

	conds = make([]Condition, 0)
	err = engine.Where(In("col_name", "col1", "col2")).Find(&conds)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(conds) != 1 {
		err = errors.New("records should exist")
		t.Error(err)
		panic(err)
	}

	// complex condtions
	var where = NewCond()
	if true {
		where = where.And(Eq{"col_name": "col1"})
		where = where.Or(And(In("col_name", "col1", "col2"), Expr("col_name = ?", "col1")))
	}

	conds = make([]Condition, 0)
	err = engine.Where(where).Find(&conds)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(conds) != 1 {
		err = errors.New("records should exist")
		t.Error(err)
		panic(err)
	}
}
