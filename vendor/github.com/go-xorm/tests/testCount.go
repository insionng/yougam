package tests

import (
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

func count(engine *xorm.Engine, t *testing.T) {
	colName := engine.ColumnMapper.Obj2Table("Departname")
	sess := engine.Where("`"+colName+"` = ?", "dev")
	total, err := sess.Clone().Count(new(Userinfo))
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Printf("Total %d records!!!\n", total)

	var users []Userinfo
	err = sess.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Printf("Total %d records!!!\n", total)
}
