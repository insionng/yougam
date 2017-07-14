package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
)

func testDelete(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{Uid: 1}
	cnt, err := engine.Delete(&user)
	if err != nil {
		t.Fatal("delete failed:", err)
	}
	if cnt != 1 {
		t.Fatal("delete failed: deleted 0 rows")
	}

	user.Uid = 0
	user.IsMan = true
	has, err := engine.Id(3).Get(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if has {
		//var tt time.Time
		//user.Created = tt
		cnt, err := engine.Id(3).Delete(new(Userinfo))
		if err != nil {
			t.Fatal("delete failed:", err)
		}
		if cnt != 1 {
			t.Fatal("delete failed: deleted 0 rows")
		}
	}
}

type Deleted struct {
	Id        int64 `xorm:"pk"`
	Name      string
	DeletedAt time.Time `xorm:"deleted"`
}

func testDeleted(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&Deleted{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Deleted{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.InsertOne(&Deleted{Id: 1, Name: "11111"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.InsertOne(&Deleted{Id: 2, Name: "22222"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.InsertOne(&Deleted{Id: 3, Name: "33333"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	// Test normal Find()
	var records1 []Deleted
	err = engine.Where("`"+engine.ColumnMapper.Obj2Table("Id")+"` > 0").Find(&records1, &Deleted{})
	if len(records1) != 3 {
		t.Fatalf("Find failed: expected=%d, actual=%d, err=%v", 3, len(records1), err)
	}

	// Test normal Get()
	record1 := &Deleted{}
	has, err := engine.Id(1).Get(record1)
	if !has {
		t.Fatalf("Get failed: expected=%v, actual=%v, err=%v", true, has, err)
	}
	//fmt.Println("----- get:", record1)

	// Test Delete() with deleted
	affected, err := engine.Id(1).Delete(&Deleted{})
	if affected != 1 {
		t.Fatalf("Delete failed: expected=%v, actual=%v, err=%v", 1, affected, err)
	}
	has, err = engine.Id(1).Get(&Deleted{})
	if has {
		t.Fatalf("Delete failed. Must not get any records.")
	}
	var records2 []Deleted
	err = engine.Where("`" + engine.ColumnMapper.Obj2Table("Id") + "` > 0").Find(&records2)
	if len(records2) != 2 {
		t.Fatalf("Find() failed.")
	}

	// Test no rows affected after Delete() again.
	affected, err = engine.Id(1).Delete(&Deleted{})
	if affected != 0 {
		t.Fatalf("Delete failed. No rows must be affected: expected=%v, actual=%v, err=%v", 0, affected, err)
	}

	// Deleted.DeletedAt must not be updated.
	affected, err = engine.Id(2).Update(&Deleted{Name: "2", DeletedAt: time.Now()})
	if affected != 1 {
		t.Fatalf("Update failed: expected=%v, actual=%v, err=%v", 1, affected, err)
	}
	record2 := &Deleted{}
	has, err = engine.Id(2).Get(record2)
	if !record2.DeletedAt.IsZero() {
		t.Fatalf("Update failed: DeletedAt must be zero value. actual=%v", record2.DeletedAt)
	}

	// Test find all records whatever `deleted`.
	var unscopedRecords1 []Deleted
	err = engine.Unscoped().Where("`"+engine.ColumnMapper.Obj2Table("Id")+"` > 0").Find(&unscopedRecords1, &Deleted{})
	if len(unscopedRecords1) != 3 {
		fmt.Printf("unscopedRecords1 = %v\n", unscopedRecords1)
		t.Fatalf("Find failed: all records must be selected when engine.Unscoped()")
	}

	// Delete() must really delete a record with Unscoped()
	affected, err = engine.Unscoped().Id(1).Delete(&Deleted{})
	if affected != 1 {
		t.Fatalf("Delete failed")
	}
	var unscopedRecords2 []Deleted
	err = engine.Unscoped().Where("`"+engine.ColumnMapper.Obj2Table("Id")+"` > 0").Find(&unscopedRecords2, &Deleted{})
	if len(unscopedRecords2) != 2 {
		t.Fatalf("Find failed: Only 2 records must be selected when engine.Unscoped()")
	}

	var records3 []Deleted
	err = engine.Where("`"+engine.ColumnMapper.Obj2Table("Id")+"` > 0").And("`"+engine.ColumnMapper.Obj2Table("Id")+"`> 1").
		Or("`"+engine.ColumnMapper.Obj2Table("Id")+"` = ?", 3).Find(&records3)
	if len(records3) != 2 {
		t.Fatalf("Find failed: expected=%d, actual=%d, err=%v", 2, len(records3), err)
	}
}
