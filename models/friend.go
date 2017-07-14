package models

import (
	"time"
)

//好友关联表
//包含  用户 ID，好友 ID
type Friend struct {
	Id      int64
	Uid     int64  `xorm:"index"`
	Fid     int64  `xorm:"index"`
	Content string `xorm:"index"`         //申请留言
	Group   string `xorm:"index"`         //好友分组
	Accept  int64  `xorm:"index"`         //是否通过好友申请，即是大于0就是所有好友 3为特别关注  2为普通 1为待处理 0为缺省值 -1为拒绝  -2黑名单
	Created int64  `xorm:"created index"` //何时相交
}

type Friendjuser struct {
	Friend `xorm:"extends",json:"Friend"`
	User   `xorm:"extends",json:"User"`
}

func DelFriend(uid int64, fid int64) (int64, error) {
	Engine.Where("uid=? and fid=?", fid, uid).Delete(new(Friend))
	return Engine.Where("uid=? and fid=?", uid, fid).Delete(new(Friend))
}

func SetFriend(uid, fid, accept int64, content, group string) (int64, error) {
	DelFriend(uid, fid)

	f := new(Friend)
	f.Uid = uid
	f.Fid = fid
	f.Content = content
	f.Group = group
	f.Accept = accept
	f.Created = time.Now().Unix()
	rows, err := Engine.Insert(f)
	//2rd
	f.Id = 0
	f.Uid = fid
	f.Fid = uid
	rows, err = Engine.Insert(f)
	return rows, err
}

func SetFriendTo(fromUid, toFid, accept int64, content, group string) (int64, error) {

	f := new(Friend)
	f.Uid = toFid
	f.Fid = fromUid
	f.Content = content
	f.Group = group
	f.Accept = accept
	f.Created = time.Now().Unix()
	rows, err := Engine.Insert(f)
	return rows, err
}

func GetRelationship(uid, fid int64) *Friend {
	f := new(Friend)
	if has, err := Engine.Where("uid=? and fid=?", uid, fid).Get(f); err != nil {
		return nil
	} else {
		if has {
			return f
		} else {
			return nil
		}
	}
}

func IsFriend(uid, fid int64) bool {

	f := new(Friend)
	if has, err := Engine.Where("uid=? and fid=? and accept>=?", uid, fid, 2).Get(f); err != nil {
		return false
	} else {
		if has {
			return true
		} else {
			return false
		}
	}

}

func GetFriends(offset int, limit int, field string) (*[]*Friend, error) {
	friends := new([]*Friend)
	err := Engine.Limit(limit, offset).Desc(field).Find(friends)
	return friends, err

}

func GetFriendsByUid(uid int64, offset int, limit int, group string, field string) *[]*Friend {
	fr := new([]*Friend)
	if field == "asc" {

		if uid == 0 {
			if len(group) > 0 {
				Engine.Table("friend").Where("friend.group=?", group).And("friend.accept>=?", 1).Limit(limit, offset).Asc("friend.id").Find(fr)
			} else {
				Engine.Table("friend").Where("friend.accept>=?", 1).Limit(limit, offset).Asc("friend.id").Find(fr)
			}

		} else {

			if len(group) == 0 {
				Engine.Table("friend").Where("friend.uid=?", uid).And("friend.accept>=?", 1).Limit(limit, offset).Asc("friend.id").Find(fr)

			} else {

				Engine.Table("friend").Where("friend.group=? and friend.uid=?", group, uid).And("friend.accept>=?", 1).Limit(limit, offset).Asc("friend.id").Find(fr)
			}
		}
	} else {

		if uid == 0 {
			if len(group) > 0 {
				Engine.Table("friend").Where("friend.group=? and friend.accept>=?", group, 1).Limit(limit, offset).Desc("friend." + field).Find(fr)
			} else {
				Engine.Table("friend").Where("friend.accept>=?", 1).Limit(limit, offset).Desc("friend." + field).Find(fr)
			}

		} else {

			if len(group) == 0 {
				Engine.Table("friend").Where("friend.uid=? and friend.accept>=?", uid, 1).Limit(limit, offset).Desc("friend." + field).Find(fr)

			} else {

				Engine.Table("friend").Where("friend.group=? and friend.uid=? and friend.accept>=?", group, uid, 1).Limit(limit, offset).Desc("friend." + field).Find(fr)
			}
		}
	}
	return fr
}

func GetFriendsByUidJoinUser(uid int64, offset int, limit int, group string, field string) *[]*Friendjuser {
	fr := new([]*Friendjuser)
	if field == "asc" {

		if uid == 0 {
			if len(group) > 0 {
				Engine.Table("friend").Where("friend.group=?", group).And("friend.accept>=?", 1).Limit(limit, offset).Asc("friend.id").Join("LEFT", "user", "user.id = friend.fid").Find(fr)
			} else {
				Engine.Table("friend").Where("friend.accept>=?", 1).Limit(limit, offset).Asc("friend.id").Join("LEFT", "user", "user.id = friend.fid").Find(fr)
			}

		} else {

			if len(group) == 0 {
				Engine.Table("friend").Where("friend.uid=?", uid).And("friend.accept>=?", 1).Limit(limit, offset).Asc("friend.id").Join("LEFT", "user", "user.id = friend.fid").Find(fr)

			} else {

				Engine.Table("friend").Where("friend.group=? and friend.uid=?", group, uid).And("friend.accept>=?", 1).Limit(limit, offset).Asc("friend.id").Join("LEFT", "user", "user.id = friend.fid").Find(fr)
			}
		}
	} else {

		if uid == 0 {
			if len(group) > 0 {
				Engine.Table("friend").Where("friend.group=? and friend.accept>=?", group, 1).Limit(limit, offset).Desc("friend."+field).Join("LEFT", "user", "user.id = friend.fid").Find(fr)
			} else {
				Engine.Table("friend").Where("friend.accept>=?", 1).Limit(limit, offset).Desc("friend."+field).Join("LEFT", "user", "user.id = friend.fid").Find(fr)
			}

		} else {

			if len(group) == 0 {
				Engine.Table("friend").Where("friend.uid=? and friend.accept>=?", uid, 1).Limit(limit, offset).Desc("friend."+field).Join("LEFT", "user", "user.id = friend.fid").Find(fr)

			} else {

				Engine.Table("friend").Where("friend.group=? and friend.uid=? and friend.accept>=?", group, uid, 1).Limit(limit, offset).Desc("friend."+field).Join("LEFT", "user", "user.id = friend.fid").Find(fr)
			}
		}
	}
	return fr
}
