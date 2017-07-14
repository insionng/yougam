//TODO 下一版本重构数据库表结构 采用空间换时间 通过数据的冗余，实现去关联查询

package models

import (
	_ "github.com/insionng/yougam/libraries/go-sql-driver/mysql"
	_ "github.com/insionng/yougam/libraries/lib/pq"
	//_ "github.com/insionng/yougam/libraries/mattn/go-sqlite3" // need cgo!

	//"github.com/go-xorm/core"
	//"github.com/go-xorm/xorm"
	_ "github.com/insionng/yougam/libraries/go-xorm/tidb"
	_ "github.com/insionng/yougam/libraries/pingcap/tidb"
	//_ "github.com/go-xorm/tidb"

	"github.com/insionng/yougam/libraries/go-xorm/core"
	"github.com/insionng/yougam/libraries/go-xorm/xorm"

	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"
	"github.com/insionng/yougam/helper"
)

var (
	Engine    *xorm.Engine
	HasEngine bool

	DataType  = helper.DataType
	DBConnect = helper.DBConnect
)

/*
	设计索引时，最好能够选择具有唯一性的字段或者重复性比较少的字段
	如此设置索引对于数据库性能来说才有比较大的价值

	通常情况下，只有当经常查询索引列中的数据时，才需要在表上创建索引。
	索引将占用磁盘空间，并且降低添加、删除和更新行的速度。
	不过在多数情况下，索引所带来的数据检索速度的优势大大超过它的不足之处。
	然而，如果应用程序非常频繁地更新数据，或磁盘空间有限，那么最好限制索引的数量。

	频繁变化的字段尽量不要设置索引
*/

//Usergroup,Pid:root
type Usergroup struct {
	Id             int64
	Pid            int64 `xorm:"index"`
	Uid            int64 `xorm:"index"` //创建者ID
	Sort           int64
	Ctype          int64   `xorm:"index"`
	Title          string  `xorm:"index"`
	Content        string  `xorm:"text"`
	Attachment     string  `xorm:"text"`
	Created        int64   `xorm:"created index"`
	Updated        int64   `xorm:"updated"`
	Hotness        float64 `xorm:"index"`
	Hotup          int64   `xorm:"index"`
	Hotdown        int64   `xorm:"index"`
	Hotscore       int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote        int64   `xorm:"index"` //Hotup  +	Hotdown
	Views          int64
	Author         string `xorm:"index"` //这里指本用户组创建者
	UserTime       int64
	UserCount      int64
	UserLastUserId int64
}

type Timeline struct {
	Id                int64
	Cid               int64 `xorm:"index"`
	Nid               int64 `xorm:"index"`
	Uid               int64 `xorm:"index"`
	Sort              int64
	Ctype             int64   `xorm:"index"`
	Title             string  `xorm:"index"`
	Content           string  `xorm:"text"`
	Attachment        string  `xorm:"text"`
	Created           int64   `xorm:"created index"`
	Updated           int64   `xorm:"updated"`
	Hotness           float64 `xorm:"index"`
	Hotup             int64   `xorm:"index"`
	Hotdown           int64   `xorm:"index"`
	Hotscore          int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64   `xorm:"index"` //Hotup  +		Hotdown
	Views             int64   `xorm:"index"`
	Author            string  `xorm:"index"`
	AuthorSignature   string  `xorm:"index"`
	Category          string  `xorm:"index"`
	Node              string  `xorm:"index"`
	ReplyTime         int64
	ReplyCount        int64 `xorm:"index"`
	ReplyLastUserId   int64
	ReplyLastUsername string
	ReplyLastNickname string
}

type ReportMark struct {
	Id      int64
	Uid     int64  `xorm:"index"` //举报人
	Rid     int64  `xorm:"index"` //被举报项id
	Ctype   int64  `xorm:"index"` // 1为topic类型  -1是comment类型
	Title   string `xorm:"index"` //举报标题
	Content string `xorm:"text"`  //举报内容
}

func ConDb() (*xorm.Engine, error) {
	switch {
	case DataType == "memory":
		return xorm.NewEngine("tidb", "memory://tidb/tidb")

	case DataType == "goleveldb":
		if DBConnect != "" {
			return xorm.NewEngine("tidb", DBConnect)
		}
		return xorm.NewEngine("tidb", "goleveldb://"+helper.FileStorageDir+"data/tidb/tidb")

	case DataType == "boltdb":
		if DBConnect != "" {
			return xorm.NewEngine("tidb", DBConnect)
		}
		return xorm.NewEngine("tidb", "boltdb://"+helper.FileStorageDir+"data/tidb/tidb")
	case DataType == "sqlite":
		if DBConnect != "" {
			return xorm.NewEngine("sqlite3", DBConnect)
		}
		return xorm.NewEngine("sqlite3", helper.FileStorageDir+"data/sqlite.db")

	case DataType == "mysql":
		return xorm.NewEngine("mysql", DBConnect)
		//return xorm.NewEngine("mysql", "root:YouPass@/db?charset=utf8")

	case DataType == "postgres":
		return xorm.NewEngine("postgres", DBConnect)
		//return xorm.NewEngine("postgres", "user=postgres password=jn!@#$%^&* dbname=pgsql sslmode=disable")

		// "user=postgres password=jn!@#$%^&* dbname=yougam sslmode=disable maxcons=10 persist=true"
		//return xorm.NewEngine("postgres", "host=192.168.1.113 user=postgres password=jn!@#$%^&* dbname=yougam sslmode=disable")
		//return xorm.NewEngine("postgres", "host=127.0.0.1 port=6432 user=postgres password=jn!@#$%^&* dbname=yougam sslmode=disable")
	}
	return nil, errors.New("Unknown database type..")
}

func SetEngine() (*xorm.Engine, error) {
	var _error error
	if Engine, _error = ConDb(); _error != nil {
		return nil, fmt.Errorf("Fail to connect to database: %s", _error.Error())
	} else {
		Engine.SetTableMapper(core.NewPrefixMapper(core.GonicMapper{}, helper.DBTablePrefix))
		Engine.SetColumnMapper(core.GonicMapper{})

		cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 10240)
		Engine.SetDefaultCacher(cacher)

		logPath := path.Join(helper.FileStorageDir+"logs", "xorm.log")
		os.MkdirAll(path.Dir(logPath), os.ModePerm)
		f, err := os.Create(logPath)
		if err != nil {
			return Engine, fmt.Errorf("Fail to create xorm.log: %s", err.Error())
		}

		Engine.SetLogger(xorm.NewSimpleLogger(f))
		Engine.ShowSQL(false)

		if location, err := time.LoadLocation("Asia/Shanghai"); err == nil {
			Engine.TZLocation = location
		}

		return Engine, err
	}
}

func NewEngine() error {
	var _error error
	Engine, _error = SetEngine()
	return _error
}

func init() {

	var _error error
	if Engine, _error = SetEngine(); _error != nil {
		log.Fatal("yougam.models.init() errors:", _error.Error())
	}

	if _error = createTables(Engine); _error != nil {
		log.Fatal("Fail to creatTables errors:", _error.Error())
	}

	initData()

}

func Ping() error {
	return Engine.Ping()
}

func createTables(Engine *xorm.Engine) error {
	return Engine.Sync2(&User{}, &Balance{}, &Message{}, &HistoryMessage{}, &Friend{}, &Category{}, &Node{}, &Topic{}, &Page{}, &Reply{}, &Attachment{}, &Link{}, &Notification{}, &UserMark{}, &TopicMark{}, &ReplyMark{}, &NodeMark{}, &ReportMark{}, &IssueMark{})
}

func initData() {
	//用户等级划分：正数是普通用户，负数是管理员各种等级划分，为0则尚未注册
	if usr, err := GetUserByRole(-1000); usr == nil && err == nil {
		if row, err := AddUser("root@yougam.com", "root", "root", "root", helper.EncryptHash("rootpass", nil), "", "", "", 1, -1000); err == nil && row > 0 {
			log.Println("Default Email:root@yougam.com ,Username:root ,Password:rootpass")

			if usr, err := GetUserByRole(-1000); usr != nil && err == nil {
				SetAmountByUid(usr.Id, 2, 2000, "注册收益2000金币")
			}

		} else {
			log.Println("create root got errors:", err)
		}

	}

	if cats, err := GetCategories(0, 0, "id"); cats == nil && err == nil {
		log.Println(AddCategory("默认分类", "默认分类内容简介", "", 0))
	}

	fmt.Println("-----------------------------------------------------------")
	fmt.Println("The yougam system has started!")
}

func Counts() (categories int, nodes int, topics int, users int, replys int) {

	var err error
	var cnt int64
	if cnt, err = Engine.Count(new(Category)); err != nil {
		log.Println(err)
	}
	categories = int(cnt)

	if cnt, err = Engine.Count(new(Node)); err != nil {
		log.Println(err)
	}
	nodes = int(cnt)

	if cnt, err = Engine.Count(new(Topic)); err != nil {
		log.Println(err)
	}
	topics = int(cnt)

	if cnt, err = Engine.Count(new(User)); err != nil {
		log.Println(err)
	}
	users = int(cnt)

	if cnt, err = GetReplysByPid4Count(0, 0, 0, 0); err != nil {
		log.Println(err)
	}
	replys = int(cnt)
	return categories, nodes, topics, users, replys
}

func SearchTopic(content string, offset int, limit int, field string) (*[]*Topic, error) {
	//排序首先是热值优先，然后是时间优先。
	if len(content) > 0 {

		keyword := "%" + content + "%"

		tps := new([]*Topic)

		err := Engine.Where("title like ? or content like ?", keyword, keyword).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
		return tps, err
	}
	return nil, errors.New("搜索内容为空!")
}

func SearchSubject(content string, offset int, limit int, field string) (*[]*Topic, error) {
	//排序首先是热值优先，然后是时间优先。
	if len(content) > 0 {

		keyword := "%" + content + "%"

		tps := new([]*Topic)

		err := Engine.Table("topic").Where("topic.pid=0 and (topic.title like ? or topic.content like ?)", keyword, keyword).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
		return tps, err
	}
	return nil, errors.New("搜索内容为空!")
}

func SearchSubjectJoinUser(content string, offset int, limit int, field string) (*[]*Topicjuser, error) {
	//排序首先是热值优先，然后是时间优先。
	if len(content) > 0 {

		keyword := "%" + content + "%"

		tps := new([]*Topicjuser)

		err := Engine.Table("topic").Where("topic.pid=0 and (topic.title like ? or topic.content like ?)", keyword, keyword).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		return tps, err
	}
	return nil, errors.New("搜索内容为空!")
}

//SearchUser 搜索用户
func SearchUser(content string, offset int, limit int, field string) (*[]*User, error) {
	//排序首先是热值优先，然后是时间优先。
	if len(content) > 0 {

		keyword := "%" + content + "%"

		usr := new([]*User)

		err := Engine.Where("username like ?", keyword).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(usr)
		return usr, err
	}
	return nil, errors.New("搜索用户为空!")
}

//SearchNode 搜索节点
func SearchNode(content string, offset int, limit int, field string) (*[]*Node, error) {
	//排序首先是热值优先，然后是时间优先。
	if len(content) > 0 {

		keyword := "%" + content + "%"

		nds := new([]*Node)

		err := Engine.Where("title like ? or content like ?", keyword, keyword).Limit(limit, offset).Desc(field, "views", "topic_count", "created").Find(nds)
		return nds, err
	}
	return nil, errors.New("搜索内容为空!")
}

//SearchCategory 搜索分类
func SearchCategory(content string, offset int, limit int, field string) (*[]*Category, error) {
	//排序首先是热值优先，然后是时间优先。
	if len(content) > 0 {

		keyword := "%" + content + "%"

		cats := new([]*Category)

		err := Engine.Where("title like ? or content like ?", keyword, keyword).Limit(limit, offset).Desc(field, "views", "node_count", "created").Find(cats)
		return cats, err
	}
	return nil, errors.New("搜索内容为空!")
}

//SetReportMark 设置举报标记
func SetReportMark(uid int64, rid int64, ctype int64, content string) (int64, error) {
	rptm := &ReportMark{Uid: uid, Rid: rid, Ctype: ctype, Content: content}
	rows, err := Engine.Insert(rptm)
	return rows, err
}

func IsReportMark(uid int64, rid int64, ctype int64) bool {

	rptm := &ReportMark{}

	if has, err := Engine.Where("uid=? and rid=? and ctype=?", uid, rid, ctype).Get(rptm); err != nil {
		fmt.Println(err)
		return false
	} else {
		if has {
			if rptm.Uid == uid {
				return true
			} else {
				return false
			}

		} else {
			return false
		}
	}

}

func GetReport(rid int64) (*ReportMark, error) {

	rm := &ReportMark{}

	has, err := Engine.Id(rid).Get(rm)
	if has {
		return rm, err
	} else {

		return nil, err
	}
}

func GetReports(offset int, limit int, field string, ctype int64) (*[]*ReportMark, error) {

	rms := new([]*ReportMark)

	if ctype == 0 { //同时查询两种类型举报项
		err := Engine.Limit(limit, offset).Desc(field).Find(rms)
		return rms, err

	} else if (ctype == 1) || (ctype == -1) { //1为topic类型 -1是comment类型 按指定类型查询举报项目

		err := Engine.Where("ctype=?", ctype).Limit(limit, offset).Desc(field).Find(rms)
		return rms, err
	} else {

		return nil, nil
	}
}

func GetNodeCountByPid(cid int64) int64 {
	n, _ := Engine.Where("pid=?", cid).Count(&Node{Pid: cid})
	return n
}

func AddTimeline(title string, content string, cid int64, nid int64, uid int64, author string, author_signature string) (int64, error) {

	id, err := Engine.Insert(&Timeline{Cid: cid, Nid: nid, Uid: uid, Title: title, Content: content, Author: author, AuthorSignature: author_signature, Created: time.Now().Unix()})

	return id, err
}

func DelTimeline(lid int64) error {
	if row, err := Engine.Id(lid).Delete(new(Timeline)); err != nil || row == 0 {
		return errors.New("删除时光记录错误!")
	} else {
		return nil
	}

}

func GetTimeline(lid int64) (*Timeline, error) {
	tl := new(Timeline)
	_, err := Engine.Where("id=?", lid).Get(tl)

	return tl, err
}

func GetTimelines(offset int, limit int, path string, uid int64) (*[]*Timeline, error) {
	tls := new([]*Timeline)
	err := errors.New("")
	if uid == 0 {
		err = Engine.Limit(limit, offset).Desc(path).Find(tls)
	} else {
		if err = Engine.Where("uid=?", uid).Limit(limit, offset).Desc(path).Find(tls); err != nil {
			err = Engine.Where("uid=?", uid).NoCache().Limit(limit, offset).Desc(path).Find(tls)
		}
	}
	return tls, err
}
