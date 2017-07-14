package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
)

type tempUser2 struct {
	TempUser   tempUser `xorm:"extends"`
	Departname string
}

type tempUser3 struct {
	Temp       *tempUser `xorm:"extends"`
	Departname string
}

type tempUser4 struct {
	TempUser2 tempUser2 `xorm:"extends"`
}

type UserAndDetail struct {
	Userinfo   `xorm:"extends"`
	Userdetail `xorm:"extends"`
}

func testExtends(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&tempUser2{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&tempUser2{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu := &tempUser2{tempUser{0, "extends"}, "dev depart"}
	_, err = engine.Insert(tu)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu2 := &tempUser2{}
	_, err = engine.Get(tu2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu3 := &tempUser2{tempUser{0, "extends update"}, ""}
	_, err = engine.Id(tu2.TempUser.Id).Update(tu3)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.DropTables(&tempUser4{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&tempUser4{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu8 := &tempUser4{tempUser2{tempUser{0, "extends"}, "dev depart"}}
	_, err = engine.Insert(tu8)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu9 := &tempUser4{}
	_, err = engine.Get(tu9)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if tu9.TempUser2.TempUser.Username != tu8.TempUser2.TempUser.Username || tu9.TempUser2.Departname != tu8.TempUser2.Departname {
		err = errors.New(fmt.Sprintln("not equal for", tu8, tu9))
		t.Error(err)
		panic(err)
	}

	tu10 := &tempUser4{tempUser2{tempUser{0, "extends update"}, ""}}
	_, err = engine.Id(tu9.TempUser2.TempUser.Id).Update(tu10)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.DropTables(&tempUser3{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&tempUser3{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu4 := &tempUser3{&tempUser{0, "extends"}, "dev depart"}
	_, err = engine.Insert(tu4)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu5 := &tempUser3{}
	_, err = engine.Get(tu5)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if tu5.Temp == nil {
		err = errors.New("error get data extends")
		t.Error(err)
		panic(err)
	}
	if tu5.Temp.Id != 1 || tu5.Temp.Username != "extends" ||
		tu5.Departname != "dev depart" {
		err = errors.New("error get data extends")
		t.Error(err)
		panic(err)
	}

	tu6 := &tempUser3{&tempUser{0, "extends update"}, ""}
	_, err = engine.Id(tu5.Temp.Id).Update(tu6)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	users := make([]tempUser3, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users) != 1 {
		err = errors.New("error get data not 1")
		t.Error(err)
		panic(err)
	}

	var info UserAndDetail
	qt := engine.Quote
	engine.Update(&Userinfo{Detail: Userdetail{Id: 1}})
	ui := engine.TableMapper.Obj2Table("Userinfo")
	ud := engine.TableMapper.Obj2Table("Userdetail")
	uiid := engine.TableMapper.Obj2Table("Id")
	udid := "detail_id"
	sql := fmt.Sprintf("select * from %s, %s where %s.%s = %s.%s",
		qt(ui), qt(ud), qt(ui), qt(udid), qt(ud), qt(uiid))
	b, err := engine.Sql(sql).Get(&info)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !b {
		err = errors.New("should has lest one record")
		t.Error(err)
		panic(err)
	}
	if info.Userinfo.Uid == 0 || info.Userdetail.Id == 0 {
		err = errors.New("all of the id should has value")
		t.Error(err)
		panic(err)
	}
	fmt.Println(info)

	fmt.Println("----join--info2")
	var info2 UserAndDetail
	b, err = engine.Table(&Userinfo{}).Join("LEFT", qt(ud), qt(ui)+"."+qt("detail_id")+" = "+qt(ud)+"."+qt(uiid)).Get(&info2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !b {
		err = errors.New("should has lest one record")
		t.Error(err)
		panic(err)
	}
	if info2.Userinfo.Uid == 0 || info2.Userdetail.Id == 0 {
		err = errors.New("all of the id should has value")
		t.Error(err)
		panic(err)
	}
	fmt.Println(info2)

	fmt.Println("----join--infos2")
	var infos2 = make([]UserAndDetail, 0)
	err = engine.Table(&Userinfo{}).Join("LEFT", qt(ud), qt(ui)+"."+qt("detail_id")+" = "+qt(ud)+"."+qt(uiid)).Find(&infos2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(infos2)

	//testExtends2(engine, t)
	//testExtends3(engine, t)
	//testExtends4(engine, t)
}

type MessageBase struct {
	Id     int64 `xorm:"int(11) pk autoincr"`
	TypeId int64 `xorm:"int(11) notnull"`
}

type Message struct {
	MessageBase `xorm:"extends"`
	Title       string    `xorm:"varchar(100) notnull"`
	Content     string    `xorm:"text notnull"`
	Uid         int64     `xorm:"int(11) notnull"`
	ToUid       int64     `xorm:"int(11) notnull"`
	CreateTime  time.Time `xorm:"datetime notnull created"`
}

type MessageUser struct {
	Id   int64
	Name string
}

type MessageType struct {
	Id   int64
	Name string
}

type MessageExtend3 struct {
	Message  `xorm:"extends"`
	Sender   MessageUser `xorm:"extends"`
	Receiver MessageUser `xorm:"extends"`
	Type     MessageType `xorm:"extends"`
}

type MessageExtend4 struct {
	Message     `xorm:"extends"`
	MessageUser `xorm:"extends"`
	MessageType `xorm:"extends"`
}

func testExtends2(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&Message{}, &MessageUser{}, &MessageType{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Message{}, &MessageUser{}, &MessageType{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var sender = MessageUser{Name: "sender"}
	var receiver = MessageUser{Name: "receiver"}
	var msgtype = MessageType{Name: "type"}
	_, err = engine.Insert(&sender, &receiver, &msgtype)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	msg := Message{
		MessageBase: MessageBase{
			Id: msgtype.Id,
		},
		Title:   "test",
		Content: "test",
		Uid:     sender.Id,
		ToUid:   receiver.Id,
	}
	_, err = engine.Insert(&msg)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var mapper = engine.TableMapper.Obj2Table
	userTableName := mapper("MessageUser")
	typeTableName := mapper("MessageType")
	msgTableName := mapper("Message")

	list := make([]Message, 0)
	err = engine.Table(msgTableName).Join("LEFT", []string{userTableName, "sender"}, "`sender`.`"+mapper("Id")+"`=`"+msgTableName+"`.`"+mapper("Uid")+"`").
		Join("LEFT", []string{userTableName, "receiver"}, "`receiver`.`"+mapper("Id")+"`=`"+msgTableName+"`.`"+mapper("ToUid")+"`").
		Join("LEFT", []string{typeTableName, "type"}, "`type`.`"+mapper("Id")+"`=`"+msgTableName+"`.`"+mapper("Id")+"`").
		Find(&list)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(list) != 1 {
		err = errors.New(fmt.Sprintln("should have 1 message, got", len(list)))
		t.Error(err)
		panic(err)
	}

	if list[0].Id != msg.Id {
		err = errors.New(fmt.Sprintln("should message equal", list[0], msg))
		t.Error(err)
		panic(err)
	}
}

func testExtends3(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&Message{}, &MessageUser{}, &MessageType{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Message{}, &MessageUser{}, &MessageType{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var sender = MessageUser{Name: "sender"}
	var receiver = MessageUser{Name: "receiver"}
	var msgtype = MessageType{Name: "type"}
	_, err = engine.Insert(&sender, &receiver, &msgtype)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	msg := Message{
		MessageBase: MessageBase{
			Id: msgtype.Id,
		},
		Title:   "test",
		Content: "test",
		Uid:     sender.Id,
		ToUid:   receiver.Id,
	}
	_, err = engine.Insert(&msg)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var mapper = engine.TableMapper.Obj2Table
	userTableName := mapper("MessageUser")
	typeTableName := mapper("MessageType")
	msgTableName := mapper("Message")

	list := make([]MessageExtend3, 0)
	err = engine.Table(msgTableName).Join("LEFT", []string{userTableName, "sender"}, "`sender`.`"+mapper("Id")+"`=`"+msgTableName+"`.`"+mapper("Uid")+"`").
		Join("LEFT", []string{userTableName, "receiver"}, "`receiver`.`"+mapper("Id")+"`=`"+msgTableName+"`.`"+mapper("ToUid")+"`").
		Join("LEFT", []string{typeTableName, "type"}, "`type`.`"+mapper("Id")+"`=`"+msgTableName+"`.`"+mapper("Id")+"`").
		Find(&list)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(list) != 1 {
		err = errors.New(fmt.Sprintln("should have 1 message, got", len(list)))
		t.Error(err)
		panic(err)
	}

	if list[0].Message.Id != msg.Id {
		err = errors.New(fmt.Sprintln("should message equal", list[0].Message, msg))
		t.Error(err)
		panic(err)
	}

	if list[0].Sender.Id != sender.Id || list[0].Sender.Name != sender.Name {
		err = errors.New(fmt.Sprintln("should sender equal", list[0].Sender, sender))
		t.Error(err)
		panic(err)
	}

	if list[0].Receiver.Id != receiver.Id || list[0].Receiver.Name != receiver.Name {
		err = errors.New(fmt.Sprintln("should receiver equal", list[0].Receiver, receiver))
		t.Error(err)
		panic(err)
	}

	if list[0].Type.Id != msgtype.Id || list[0].Type.Name != msgtype.Name {
		err = errors.New(fmt.Sprintln("should msgtype equal", list[0].Type, msgtype))
		t.Error(err)
		panic(err)
	}
}

func testExtends4(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&Message{}, &MessageUser{}, &MessageType{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Message{}, &MessageUser{}, &MessageType{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var sender = MessageUser{Name: "sender"}
	var msgtype = MessageType{Name: "type"}
	_, err = engine.Insert(&sender, &msgtype)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	msg := Message{
		MessageBase: MessageBase{
			Id: msgtype.Id,
		},
		Title:   "test",
		Content: "test",
		Uid:     sender.Id,
	}
	_, err = engine.Insert(&msg)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var mapper = engine.TableMapper.Obj2Table
	userTableName := mapper("MessageUser")
	typeTableName := mapper("MessageType")
	msgTableName := mapper("Message")

	list := make([]MessageExtend4, 0)
	err = engine.Table(msgTableName).Join("LEFT", userTableName, "`"+userTableName+"`.`"+mapper("Id")+"`=`"+msgTableName+"`.`"+mapper("Uid")+"`").
		Join("LEFT", typeTableName, "`"+typeTableName+"`.`"+mapper("Id")+"`=`"+msgTableName+"`.`"+mapper("Id")+"`").
		Find(&list)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(list) != 1 {
		err = errors.New(fmt.Sprintln("should have 1 message, got", len(list)))
		t.Error(err)
		panic(err)
	}

	if list[0].Message.Id != msg.Id {
		err = errors.New(fmt.Sprintln("should message equal", list[0].Message, msg))
		t.Error(err)
		panic(err)
	}

	if list[0].MessageUser.Id != sender.Id || list[0].MessageUser.Name != sender.Name {
		err = errors.New(fmt.Sprintln("should sender equal", list[0].MessageUser, sender))
		t.Error(err)
		panic(err)
	}

	if list[0].MessageType.Id != msgtype.Id || list[0].MessageType.Name != msgtype.Name {
		err = errors.New(fmt.Sprintln("should msgtype equal", list[0].MessageType, msgtype))
		t.Error(err)
		panic(err)
	}
}
