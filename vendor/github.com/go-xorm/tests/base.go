package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

const (
	CreateTableMySql = "CREATE TABLE IF NOT EXISTS `big_struct` (`id` BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL, `name` VARCHAR(255) NULL, `title` VARCHAR(255) NULL, `age` VARCHAR(255) NULL, `alias` VARCHAR(255) NULL, `nick_name` VARCHAR(255) NULL);"
	DropTableMySql   = "DROP TABLE IF EXISTS `big_struct`;"
)

var ShowTestSql bool = true

/*
CREATE TABLE `userinfo` (
    `id` INT(10) NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NULL,
    `departname` VARCHAR(64) NULL,
    `created` DATE NULL,
    PRIMARY KEY (`uid`)
);
CREATE TABLE `userdeatail` (
    `id` INT(10) NULL,
    `intro` TEXT NULL,
    `profile` TEXT NULL,
    PRIMARY KEY (`uid`)
);
*/

type Userinfo struct {
	Uid        int64  `xorm:"id pk not null autoincr"`
	Username   string `xorm:"unique"`
	Departname string
	Alias      string `xorm:"-"`
	Created    time.Time
	Detail     Userdetail `xorm:"detail_id int(11)"`
	Height     float64
	Avatar     []byte
	IsMan      bool
}

type Userdetail struct {
	Id      int64
	Intro   string `xorm:"text"`
	Profile string `xorm:"varchar(2000)"`
}

type Picture struct {
	Id          int64
	Url         string `xorm:"unique"` //image's url
	Title       string
	Description string
	Created     time.Time `xorm:"created"`
	ILike       int
	PageView    int
	From_url    string
	Pre_url     string `xorm:"unique"` //pre view image's url
	Uid         int64
}

type Numeric struct {
	Numeric float64 `xorm:"numeric(26,2)"`
}

func NewCacher() core.Cacher {
	return xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 10000)
}

type Article struct {
	Id      int32  `xorm:"pk INT autoincr"`
	Name    string `xorm:"VARCHAR(45)"`
	Img     string `xorm:"VARCHAR(100)"`
	Aside   string `xorm:"VARCHAR(200)"`
	Desc    string `xorm:"VARCHAR(200)"`
	Content string `xorm:"TEXT"`
	Status  int8   `xorm:"TINYINT(4)"`
}

type Limit struct {
	Id      int64
	Name    string
	Updated time.Time `xorm:"updated"`
}

func limit(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Limit(2, 1).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	err = engine.DropTables(new(Limit))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Sync2(new(Limit))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(&Limit{Name: "1"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	// TODO: support limit for update
	/*_, err = engine.Limit(1).Update(&Limit{Name: "2"})
	if err != nil {
		t.Error(err)
		panic(err)
	}*/
}

func tableOp(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{Username: "tablexiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
	tableName := fmt.Sprintf("user_%v", len(user.Username))
	cnt, err := engine.Table(tableName).Insert(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	has, err := engine.Table(tableName).Get(&Userinfo{Username: "tablexiao"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("Get has return false")
		t.Error(err)
		panic(err)
		return
	}

	users := make([]Userinfo, 0)
	err = engine.Table(tableName).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	id := user.Uid
	cnt, err = engine.Table(tableName).Id(id).Update(&Userinfo{Username: "tableda"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Table(tableName).Id(id).Delete(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testCharst(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables("user_charset")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Charset("utf8").Table("user_charset").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testStoreEngine(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables("user_store_engine")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.StoreEngine("InnoDB").Table("user_store_engine").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

type tempUser struct {
	Id       int64
	Username string
}

func testCols(engine *xorm.Engine, t *testing.T) {
	users := []Userinfo{}
	err := engine.Cols("id, username").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(users)

	tmpUsers := []tempUser{}
	err = engine.NoCache().Table("userinfo").Cols("id, username").Find(&tmpUsers)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(tmpUsers)

	user := &Userinfo{Uid: 1, Alias: "", Height: 0}
	affected, err := engine.Cols("departname, height").Id(1).Update(user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println("===================", user, affected)
}

func testColsSameMapper(engine *xorm.Engine, t *testing.T) {
	users := []Userinfo{}
	err := engine.Cols("id, Username").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(users)

	tmpUsers := []tempUser{}
	// TODO: should use cache
	err = engine.NoCache().Table("Userinfo").Cols("id, Username").Find(&tmpUsers)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(tmpUsers)

	user := &Userinfo{Uid: 1, Alias: "", Height: 0}
	affected, err := engine.Cols("Departname, Height").Update(user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println("===================", user, affected)
}

type UserCU struct {
	Id      int64
	Name    string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func testCreatedAndUpdated(engine *xorm.Engine, t *testing.T) {
	u := new(UserCU)
	err := engine.DropTables(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	u.Name = "sss"
	cnt, err := engine.Insert(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	u.Name = "xxx"
	cnt, err = engine.Id(u.Id).Update(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	u.Id = 0
	u.Created = time.Now().Add(-time.Hour * 24 * 365)
	u.Updated = u.Created
	fmt.Println(u)
	cnt, err = engine.NoAutoTime().Insert(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
}

type StrangeName struct {
	Id_t int64 `xorm:"pk autoincr"`
	Name string
}

func testStrangeName(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(new(StrangeName))
	if err != nil {
		t.Error(err)
	}

	err = engine.CreateTables(new(StrangeName))
	if err != nil {
		t.Error(err)
	}

	_, err = engine.Insert(&StrangeName{Name: "sfsfdsfds"})
	if err != nil {
		t.Error(err)
	}

	beans := make([]StrangeName, 0)
	err = engine.Find(&beans)
	if err != nil {
		t.Error(err)
	}
}

func testPrefixTableName(engine *xorm.Engine, t *testing.T) {
	/*tempEngine, err := xorm.NewEngine(engine.DriverName(), engine.DataSourceName())
	//tempEngine, err := engine.Clone()
	if err != nil {
		t.Error(err)
		panic(err)
		return
	}
	defer tempEngine.Close()

	tempEngine.ShowSQL = true
	mapper := core.NewPrefixMapper(core.SnakeMapper{}, "xlw_")
	tempEngine.SetTableMapper(mapper)
	exist, err := tempEngine.IsTableExist(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
		return
	}
	if exist {
		err = tempEngine.DropTables(&Userinfo{})
		if err != nil {
			t.Error(err)
			panic(err)
			return
		}
	}

	err = tempEngine.CreateTables(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}*/
}

type CreatedUpdated struct {
	Id       int64
	Name     string
	Value    float64   `xorm:"numeric"`
	Created  time.Time `xorm:"created"`
	Created2 time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

func testCreatedUpdated(engine *xorm.Engine, t *testing.T) {
	err := engine.Sync(&CreatedUpdated{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	c := &CreatedUpdated{Name: "test"}
	_, err = engine.Insert(c)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	c2 := new(CreatedUpdated)
	has, err := engine.Id(c.Id).Get(c2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		panic(errors.New("no id"))
	}

	c2.Value -= 1
	_, err = engine.Id(c2.Id).Update(c2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

type Lowercase struct {
	Id    int64
	Name  string
	ended int64 `xorm:"-"`
}

func testLowerCase(engine *xorm.Engine, t *testing.T) {
	err := engine.Sync(&Lowercase{})
	_, err = engine.Where("(id) > 0").Delete(&Lowercase{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	_, err = engine.Insert(&Lowercase{ended: 1})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ls := make([]Lowercase, 0)
	err = engine.Find(&ls)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(ls) != 1 {
		err = errors.New("should be 1")
		t.Error(err)
		panic(err)
	}
}

func BaseTestAll(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- directCreateTable --------------")
	directCreateTable(engine, t)
	fmt.Println("-------------- insert --------------")
	insert(engine, t)
	fmt.Println("-------------- testInsertDefault --------------")
	testInsertDefault(engine, t)
	fmt.Println("-------------- insertAutoIncr --------------")
	insertAutoIncr(engine, t)
	fmt.Println("-------------- insertMulti --------------")
	insertMulti(engine, t)
	fmt.Println("-------------- insertTwoTable --------------")
	insertTwoTable(engine, t)
	fmt.Println("-------------- testDelete --------------")
	testDelete(engine, t)
	fmt.Println("-------------- get --------------")
	get(engine, t)
	fmt.Println("-------------- testCascade --------------")
	testCascade(engine, t)
	fmt.Println("-------------- testFind --------------")
	testFind(engine, t)
	fmt.Println("-------------- count --------------")
	count(engine, t)
	fmt.Println("-------------- where --------------")
	where(engine, t)
	fmt.Println("-------------- in --------------")
	in(engine, t)
	fmt.Println("-------------- limit --------------")
	limit(engine, t)
	fmt.Println("-------------- testCustomTableName --------------")
	testCustomTableName(engine, t)
	fmt.Println("-------------- testDump --------------")
	testDump(engine, t)
	fmt.Println("-------------- testConversion --------------")
	testConversion(engine, t)
	fmt.Println("-------------- testJsonField --------------")
	testJsonField(engine, t)
	fmt.Println("-------------- testSum --------------")
	testSum(engine, t)
}

func BaseTestAll2(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- table --------------")
	table(engine, t)
	fmt.Println("-------------- createMultiTables --------------")
	createMultiTables(engine, t)
	fmt.Println("-------------- tableOp --------------")
	tableOp(engine, t)
	fmt.Println("-------------- testCharst --------------")
	testCharst(engine, t)
	fmt.Println("-------------- testStoreEngine --------------")
	testStoreEngine(engine, t)
	fmt.Println("-------------- testExtends --------------")
	testExtends(engine, t)
	fmt.Println("-------------- testColTypes --------------")
	testColTypes(engine, t)
	fmt.Println("-------------- testCustomType --------------")
	testCustomType(engine, t)
	fmt.Println("-------------- testCreatedAndUpdated --------------")
	testCreatedAndUpdated(engine, t)
	fmt.Println("-------------- testIndexAndUnique --------------")
	testIndexAndUnique(engine, t)
	fmt.Println("-------------- testIntId --------------")
	testIntId(engine, t)
	fmt.Println("-------------- testInt16Id --------------")
	testInt16Id(engine, t)
	fmt.Println("-------------- testInt32Id --------------")
	testInt32Id(engine, t)
	fmt.Println("-------------- testUintId --------------")
	testUintId(engine, t)
	fmt.Println("-------------- testUint16Id --------------")
	testUint16Id(engine, t)
	fmt.Println("-------------- testUint32Id --------------")
	testUint32Id(engine, t)
	fmt.Println("-------------- testUint64Id --------------")
	testUint64Id(engine, t)
	fmt.Println("-------------- testMyIntId --------------")
	testMyIntId(engine, t)
	fmt.Println("-------------- testMyStringId --------------")
	testMyStringId(engine, t)
	fmt.Println("-------------- testMetaInfo --------------")
	testMetaInfo(engine, t)
	fmt.Println("-------------- testIterate --------------")
	testIterate(engine, t)
	fmt.Println("-------------- testRows --------------")
	testRows(engine, t)
	fmt.Println("-------------- testStrangeName --------------")
	testStrangeName(engine, t)
	fmt.Println("-------------- testVersion --------------")
	testVersion(engine, t)
	fmt.Println("-------------- testDistinct --------------")
	testDistinct(engine, t)
	fmt.Println("-------------- testUseBool --------------")
	testUseBool(engine, t)
	fmt.Println("-------------- testBool --------------")
	testBool(engine, t)
	fmt.Println("-------------- testTime --------------")
	testTime(engine, t)
	fmt.Println("-------------- testPrefixTableName --------------")
	testPrefixTableName(engine, t)
	fmt.Println("-------------- testCreatedUpdated --------------")
	testCreatedUpdated(engine, t)
	fmt.Println("-------------- testLowercase ---------------")
	testLowerCase(engine, t)
	fmt.Println("-------------- processors --------------")
	testProcessors(engine, t)
	fmt.Println("-------------- transaction --------------")
	transaction(engine, t)
	fmt.Println("-------------- testCacheDomain --------------")
	testCacheDomain(engine, t)
	fmt.Println("-------------- testDeleted --------------")
	testDeleted(engine, t)
	fmt.Println("-------------- testCompositeKey --------------")
	testCompositeKey(engine, t)
	fmt.Println("-------------- testCompositeKey2 --------------")
	testCompositeKey2(engine, t)
	fmt.Println("-------------- testCompositeKey3 --------------")
	testCompositeKey3(engine, t)
	fmt.Println("-------------- testStringPK --------------")
	testStringPK(engine, t)
	fmt.Println("-------------- testForUpdate --------------")
	testForUpdate(engine, t)
	fmt.Println("-------------- testID --------------")
	testID(engine, t)
}

// !nash! the 3rd set of the test is intended for non-cache enabled engine
func BaseTestAll3(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- processors TX --------------")
	testProcessorsTx(engine, t)
	fmt.Println("-------------- insert pointer data --------------")
	testPointerData(engine, t)
	fmt.Println("-------------- insert pointer to aliased types --------------")
	testPointersToAliases(engine, t)
	fmt.Println("-------------- insert null data --------------")
	testNullValue(engine, t)
	fmt.Println("-------------- testNoCacheDomain --------------")
	testNoCacheDomain(engine, t)
	fmt.Println("-------------- testCache2 --------------")
	testCache2(engine, t)
}

func BaseTestAllSnakeMapper(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- query --------------")
	testQuery(engine, t)
	fmt.Println("-------------- find3 --------------")
	find3(engine, t)
	fmt.Println("-------------- exec --------------")
	exec(engine, t)
	fmt.Println("-------------- update --------------")
	update(engine, t)
	fmt.Println("-------------- order --------------")
	order(engine, t)
	fmt.Println("-------------- join --------------")
	join(engine, t)
	fmt.Println("-------------- having --------------")
	having(engine, t)
	fmt.Println("-------------- combineTransaction --------------")
	combineTransaction(engine, t)
	fmt.Println("-------------- testCols --------------")
	testCols(engine, t)
	fmt.Println("-------------- testNullStruct --------------")
	TestNullStruct(engine, t)
	fmt.Println("-------------- testBuilder --------------")
	testBuilder(engine, t)
}

func BaseTestAllSameMapper(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- query --------------")
	testQuerySameMapper(engine, t)
	fmt.Println("-------------- exec --------------")
	execSameMapper(engine, t)
	fmt.Println("-------------- update --------------")
	updateSameMapper(engine, t)
	fmt.Println("-------------- order --------------")
	orderSameMapper(engine, t)
	fmt.Println("-------------- join --------------")
	joinSameMapper(engine, t)
	fmt.Println("-------------- having --------------")
	havingSameMapper(engine, t)
	fmt.Println("-------------- combineTransaction --------------")
	combineTransactionSameMapper(engine, t)
	fmt.Println("-------------- testCols --------------")
	testColsSameMapper(engine, t)
}
