package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
)

func testCascade(engine *xorm.Engine, t *testing.T) {
	cascadeGet(engine, t)
	UserTest1(engine, t)
}

func cascadeGet(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{Uid: 11}

	has, err := engine.Get(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if has {
		fmt.Println(user)
	} else {
		fmt.Println("no record id is 2")
	}
}

type Users struct {
	Uid      string `xorm:"notnull pk UUID"`
	UserName string `xorm:"notnull unique VARCHAR(30)"`
	NickName string `xorm:"notnull VARCHAR(30)"`
	Password string `xorm:"notnull VARCHAR(44)"`
	Email    string `xorm:"notnull unique VARCHAR(80)"`

	Profile *Profile `xorm:"profile_id UUID"`

	CreatedAt   time.Time `xorm:"notnull TIMESTAMPZ created"`
	UpdatedAt   time.Time `xorm:"notnull TIMESTAMPZ updated"`
	CreatedUser *Users    `xorm:"created_user_id UUID"`
	UpdatedUser *Users    `xorm:"update_user_id UUID"`
}

type Profile struct {
	Uid string `xorm:"notnull pk UUID"`

	User *Users `xorm:"user_id varchar(50)"`

	CreatedAt   time.Time `xorm:"notnull TIMESTAMPZ created"`
	UpdatedAt   time.Time `xorm:"notnull TIMESTAMPZ updated"`
	CreatedUser *Users    `xorm:"created_user_id varchar(50)"`
	UpdatedUser *Users    `xorm:"update_user_id varchar(50)"`
}

func UserTest1(engine *xorm.Engine, t *testing.T) {
	table1, table2 := new(Users), new(Profile)
	engine.DropTables(table1, table2)
	err := engine.Sync2(table1, table2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	profile := &Profile{
		Uid:         "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		User:        nil,
		CreatedUser: nil,
		UpdatedUser: nil,
	}

	users := &Users{
		Uid:      "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		UserName: "sfafdafds",
		NickName: "sfsfds",
		Password: "ssss",
		Email:    "llll@123.com",
		Profile: &Profile{
			Uid: profile.Uid,
		},
		CreatedUser: nil,
		UpdatedUser: nil,
	}

	_, err = engine.Insert(users)
	if err != nil {
		t.Fatal("insert UUID PK users failed:", err)
	}

	profile.CreatedUser = users

	_, err = engine.Insert(profile)
	if err != nil {
		t.Fatal("insert UUID PK profile failed:", err)
	}

	var newUsers Users
	has, err := engine.Id(users.Uid).Get(&newUsers)
	if err != nil {
		t.Fatal("get UUID pk users failed:", err)
	}
	if !has {
		t.Fatal("get UUID pk users failed: get none")
	}
	fmt.Println(newUsers.Profile, "==", profile)
	if newUsers.Profile.Uid != profile.Uid {
		t.Fatal("should equal profile.uid")
	}
}
