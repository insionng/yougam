package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

type ForUpdate struct {
	Id   int64 `xorm:"pk"`
	Name string
}

func setupForUpdate(engine *xorm.Engine) error {
	v := new(ForUpdate)
	err := engine.DropTables(v)
	if err != nil {
		return err
	}
	err = engine.CreateTables(v)
	if err != nil {
		return err
	}

	list := []ForUpdate{
		{1, "data1"},
		{2, "data2"},
		{3, "data3"},
	}

	for _, f := range list {
		_, err = engine.Insert(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestForUpdate(engine *xorm.Engine, t *testing.T) {
	testForUpdate(engine, t)
}

func testForUpdate(engine *xorm.Engine, t *testing.T) {
	if engine.DriverName() == "tidb" || engine.DriverName() == "sqlite3" || engine.Dialect().DBType() == core.MSSQL {
		return
	}
	err := setupForUpdate(engine)
	if err != nil {
		t.Error(err)
		return
	}

	session1 := engine.NewSession()
	session2 := engine.NewSession()
	session3 := engine.NewSession()
	defer session1.Close()
	defer session2.Close()
	defer session3.Close()

	// start transaction
	err = session1.Begin()
	if err != nil {
		t.Error(err)
		return
	}

	// use lock
	fList := make([]ForUpdate, 0)
	session1.ForUpdate()
	session1.Where("(id) = ?", 1)
	err = session1.Find(&fList)
	switch {
	case err != nil:
		t.Error(err)
		return
	case len(fList) != 1:
		t.Errorf("find not returned single row")
		return
	case fList[0].Name != "data1":
		t.Errorf("for_update.name must be `data1`")
		return
	}

	// wait for lock
	wg := &sync.WaitGroup{}

	// lock is used
	wg.Add(1)
	go func() {
		f2 := new(ForUpdate)
		session2.Where("(id) = ?", 1).ForUpdate()
		has, err := session2.Get(f2) // wait release lock
		switch {
		case err != nil:
			t.Error(err)
		case !has:
			t.Errorf("cannot find target row. for_update.id = 1")
		case f2.Name != "updated by session1":
			t.Errorf("read lock failed")
		}
		wg.Done()
	}()

	// lock is NOT used
	wg.Add(1)
	go func() {
		f3 := new(ForUpdate)
		session3.Where("(id) = ?", 1)
		has, err := session3.Get(f3) // wait release lock
		switch {
		case err != nil:
			t.Error(err)
		case !has:
			t.Errorf("cannot find target row. for_update.id = 1")
		case f3.Name != "data1":
			t.Errorf("read lock failed")
		}
		wg.Done()
	}()

	// wait for go rountines
	time.Sleep(50 * time.Millisecond)

	f := new(ForUpdate)
	f.Name = "updated by session1"
	session1.Where("(id) = ?", 1)
	session1.Update(f)

	// release lock
	err = session1.Commit()
	if err != nil {
		t.Error(err)
		return
	}

	wg.Wait()
}
