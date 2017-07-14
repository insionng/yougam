package models

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/insionng/yougam/helper"
)

type User struct {
	Id                int64
	Pid               int64  `xorm:"index"` //用在归属地 归属学校 归属组织 等方面
	Group             string `xorm:"index"` //小组名称集合
	Email             string `xorm:"index"`
	Password          string `xorm:"index" json:"-"` //禁止Json输出Password字段
	Username          string `xorm:"index"`
	Nickname          string `xorm:"index"`
	Realname          string `xorm:"index"`
	Content           string `xorm:"text"`
	Avatar            string `xorm:"index"` //200x200
	AvatarLarge       string `xorm:"index"` //100x100
	AvatarMedium      string `xorm:"index"` //48x48
	AvatarSmall       string `xorm:"index"` //32x32
	Birth             int64
	Province          string
	City              string
	Occupation        string //职业
	Company           string
	Address           string
	Postcode          string
	Mobile            string `xorm:"index"`
	Website           string
	Gender            int64 `xorm:"index"` // 1==男 -1==女
	Age               int64
	School            string
	Weight            int64
	Height            int64
	ZodiacSign        string //星座
	Qq                string
	Weixin            string
	WeixinOpenId      string `xorm:"index"` //微信移动端OpenId
	Weibo             string
	Ctype             int64   `xorm:"index"`
	Role              int64   `xorm:"index"` //角色属性同时也是权限属性 类似linux权限的形式 ，譬如：664/755/777，TODO
	Created           int64   `xorm:"created index"`
	Updated           int64   `xorm:"updated"`
	Hotness           float64 `xorm:"index"`
	Confidence        float64 `xorm:"index"` //信任度数值
	Hotup             int64   `xorm:"index"`
	Hotdown           int64   `xorm:"index"`
	Hotscore          int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64   `xorm:"index"` //Hotup  + 	Hotdown
	Views             int64
	LastSigninTime    int64
	LastSigninIp      string `xorm:"index"`
	SigninCount       int64  `xorm:"index"`
	NodeTime          int64
	NodeCount         int64
	NodeLastTid       int64
	NodeLastTopic     string
	TopicTime         int64
	TopicCount        int64
	TopicLastNid      int64
	TopicLastNode     string
	ReplyTime         int64
	ReplyCount        int64
	ReplyLastTid      int64
	ReplyLastTopic    string //topic title
	FavoriteCount     int64
	NotificationCount int64
	Balance           int64
	//Version           int64 `xorm:"version"` //乐观锁
}

type UserMark struct {
	Id   int64
	Uid  int64 `xorm:"index"` //投票本人
	User int64 `xorm:"index"` //user id，被投票人
}

func SetUserMark(uid int64, UserID int64) (int64, error) {
	um := new(UserMark)
	um.Uid = uid
	um.User = UserID
	rows, err := Engine.Insert(um)
	return rows, err
}

func IsUserMark(uid int64, UserID int64) bool {

	um := new(UserMark)

	if has, err := Engine.Where("uid=? and user=?", uid, UserID).Get(um); err != nil {
		return false
	} else {
		if has {
			return (um.Uid == uid)
		} else {
			return false
		}
	}

}

func GetUsersCount(offset int, limit int) (int64, error) {
	total, err := Engine.Limit(limit, offset).Count(&User{})
	return total, err
}

func HasUser(username string) (int64, bool) {
	user := new(User)
	if has, err := Engine.Where("username = ?", username).Get(user); err != nil {
		return -1, false
	} else {
		if has {
			return user.Id, true
		}
		return -1, false
	}
}

func SetUser(uid int64, usr *User) (int64, error) {
	usr.Id = uid
	return Engine.Insert(usr)
}

func PutUser(uid int64, usr *User) (int64, error) {
	row, err := Engine.Update(usr, &User{Id: uid})
	return row, err
}

func SetBalanceForUser(uid, balance int64) error {
	/*
		user, e := GetUser(uid)
		if e != nil {
			return e
		}
	*/
	user := &User{}
	user.Balance = balance

	if affected, err := Engine.Where("id=?", uid).Cols("balance").Update(user); err != nil || affected <= 0 {
		return errors.New(fmt.Sprintf("SetBalanceForUser() Update() Errors:%v", err))
	}

	return nil
}

func GetUserByRole(role int64) (*User, error) {
	user := &User{}
	if has, err := Engine.Where("role=?", role).Get(user); has {
		return user, err
	} else {
		return nil, err
	}
}

func GetUsers(offset int, limit int, field string) (*[]*User, error) {
	users := &[]*User{}
	err := Engine.Limit(limit, offset).Desc(field).Find(users)
	return users, err
}

func GetUsersByRole(role int64, offset int, limit int, field string) (*[]*User, error) {
	users := &[]*User{}
	err := Engine.Where("role=?", role).Limit(limit, offset).Desc(field).Desc("hotness").Find(users)
	return users, err
}

func GetUsersOnHotness(offset int, limit int, field string) (*[]*User, error) {
	users := &[]*User{}
	err := Engine.Limit(limit, offset).Desc(field).Desc("hotness").Find(users)
	return users, err
}

func GetUsersOnConfidence(offset int, limit int, field string) (*[]*User, error) {
	users := &[]*User{}
	err := Engine.Limit(limit, offset).Desc(field).Desc("confidence").Find(users)
	return users, err
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	if has, err := Engine.Where("username=?", username).Get(user); has {
		return user, err
	} else {
		return nil, err
	}
}

func GetUserByMobile(mobile string) (*User, error) {
	user := &User{}
	if has, err := Engine.Where("mobile=?", mobile).Get(user); has {
		return user, err
	} else {
		return nil, err
	}
}

func GetUserByNickname(nickname string) (*User, error) {

	user := &User{Nickname: nickname}
	has, err := Engine.Get(user)
	if has {
		return user, err
	} else {

		return nil, err
	}
}

//返回值尽量返回指针 不然会出现诡异的问题
func GetUserByEmail(email string) (*User, error) {

	var user = &User{Email: email}
	has, err := Engine.Get(user)
	if has {
		return user, err
	} else {

		return nil, err
	}
}

func SetUserNewpassword(username string, email string, oldpassword string, newpassword string) (bool, error) {

	user := &User{}
	if has, _ := Engine.Where("username=? and email=?", username, email).Get(user); has {

		if helper.ValidateHash(user.Password, oldpassword) {
			user.Password = helper.EncryptHash(newpassword, nil)
			if row, err := PutUser(user.Id, user); row > 0 && err == nil {
				return true, err

			} else {
				return false, errors.New("更新密码失败!")
			}
		} else {
			return false, errors.New("校验密码失败!")

		}
	} else {
		return false, errors.New("用户不存在!")
	}

}

func GetUser(id int64) (*User, error) {

	user := new(User)
	has, err := Engine.Id(id).Get(user)
	if has {
		return user, err
	} else {

		return nil, err
	}
}

func AddUser(email string, username string, nickname string, realname string, password string, group string, content string, mobile string, gender int64, role int64) (int64, error) {

	usr := new(User)
	usr.Email = email
	usr.Password = password
	usr.Username = username
	usr.Nickname = nickname
	usr.Realname = realname
	usr.Group = group
	usr.Content = content
	usr.Mobile = mobile
	usr.Gender = gender
	usr.Role = role
	usr.Created = time.Now().Unix()
	usr.LastSigninTime = time.Now().Unix()

	//if row, err := CacheInsert(Engine, usr); err == nil && row > 0 {

	if row, err := Engine.Insert(usr); err == nil && row > 0 {

		//if usr, err := GetUserByEmail(email); usr != nil && err == nil {

		return usr.Id, err
		//} else {

		//	return -1, err
		//}
	} else {

		return -1, err
	}

}

func PostUser(m *User) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	Userid := int64(0)
	if _, err := sess.Insert(m); err == nil && m != nil {
		Userid = m.Id
	} else {
		// 发生错误时进行回滚
		sess.Rollback()
		return -1, err
	}

	// 提交事务
	return Userid, sess.Commit()
}

// DelUser id：被删除用户之ID，uid：当前用户之id，role：当前用户之角色
func DelUser(id int64, uid int64, role int64) error {
	allow := false
	if role < 0 {
		allow = true
	}

	user := &User{}

	if has, err := Engine.Id(id).Get(user); has == true && err == nil {

		if user.Id == uid || allow {
			//检查附件字段并尝试删除文件
			if user.Avatar != "" {

				if p := helper.URL2local(user.Avatar); helper.Exist(p) {
					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误日志，但不要反回错误，以免陷入死循环无法删掉
							log.Println("ROOT DEL User Avatar, User ID:", id, ",ERR:", err)
						}
					} else { //检查用户对本地文件的所有权
						if helper.VerifyUserfile(p, strconv.FormatInt(uid, 10)) {
							if err := os.Remove(p); err != nil {
								log.Println("DEL User Avatar, User ID:", id, ",ERR:", err)
							}
						}
					}

				}
			}

			//不管实际路径中是否存在文件均删除该数据库记录，以免数据库记录陷入死循环无法删掉
			if user.Id == id {

				if row, err := Engine.Id(id).Delete(new(User)); err != nil || row == 0 {
					return errors.New(fmt.Sprint("删除用户错误!", row, err)) //错误还要我自己构造?!
				} else {
					return nil
				}

			}
		}
		return errors.New("你无权删除此用户:" + strconv.FormatInt(id, 10))
	}
	return errors.New("无法删除不存在的NODE ID:" + strconv.FormatInt(id, 10))
}

func PutSignin2User(uid, LastSigninTime, SigninCount int64, LastSigninIp string) (int64, error) {
	return Engine.Table(&User{}).Where("id=?", uid).Update(&map[string]interface{}{
		"last_signin_time": LastSigninTime,
		"last_signin_ip":   LastSigninIp,
		"signin_count":     SigninCount,
	})
}

func PutSignout2User(uid, LastSigninTime int64, LastSigninIp string) (int64, error) {
	return Engine.Table(&User{}).Where("id=?", uid).Update(&map[string]interface{}{
		"last_signin_time": LastSigninTime,
		"last_signin_ip":   LastSigninIp,
	})
}
