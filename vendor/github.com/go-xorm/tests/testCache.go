package tests

import (
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

type CacheDomain struct {
	Id   int64 `xorm:"pk cache"`
	Name string
}

func testCacheDomain(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&CacheDomain{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	table := engine.TableInfo(&CacheDomain{})
	if table.Cacher == nil {
		err = fmt.Errorf("table cache is nil")
		t.Error(err)
		panic(err)
	}
}

type ActiveData struct {
	Area string
}

func testCache2(engine *xorm.Engine, t *testing.T) {
	/*err := engine.Table("active-2016-01-02").CreateTable(&ActiveData{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	l := make([]*ActiveData, 0)
	sess := engine.NoCache().Table("active-2016-01-02")
	sess.Where("area=", "1").Find(&l)
	sql, _ := sess.LastSQL()
	fmt.Println("----", sql)
	if !strings.Contains(sql, "active-2016-01-02") {
		err = errors.New("no record")
		t.Error(err)
		panic(err)
	}*/
}

type NoCacheDomain struct {
	Id   int64 `xorm:"pk nocache"`
	Name string
}

func testNoCacheDomain(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&NoCacheDomain{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	table := engine.TableInfo(&NoCacheDomain{})
	if table.Cacher != nil {
		err = fmt.Errorf("table cache exist")
		t.Error(err)
		panic(err)
	}
}
