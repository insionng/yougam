package models

import (
	"errors"
	"fmt"
	"time"
)

type Notification struct {
	Id           int64
	Uid          int64  `xorm:"index"`
	Tid          int64  `xorm:"index"`
	Rid          int64  `xorm:"index"`
	Ctype        int64  `xorm:"index"` //0普通通知 1通知作者 -1忽略
	Subject      string `xorm:"text"`
	Reply        string `xorm:"text"`
	Author       string `xorm:"index"`
	Avatar       string `xorm:"index"` //200x200
	AvatarLarge  string `xorm:"index"` //100x100
	AvatarMedium string `xorm:"index"` //48x48
	AvatarSmall  string `xorm:"index"` //32x32
	Created      int64  `xorm:"created index"`
}

type Notificationjuser struct {
	Notification `xorm:"extends"`
	User         `xorm:"extends"`
}

type Notificationjtopic struct {
	Notification `xorm:"extends"`
	Topic        `xorm:"extends"`
}

func GetNotificationsByUid(uid int64, ctype int64, offset int, limit int, field string) *[]Notification {
	rp := &[]Notification{}
	if field == "asc" {

		if uid == 0 {
			if ctype != 0 {
				Engine.Table("notification").Where("notification.ctype=?", ctype).Limit(limit, offset).Asc("notification.id").Find(rp)
			} else {
				Engine.Table("notification").Limit(limit, offset).Asc("notification.id").Find(rp)
			}

		} else {

			if ctype == 0 {
				Engine.Table("notification").Where("notification.uid=?", uid).Limit(limit, offset).Asc("notification.id").Find(rp)

			} else {

				Engine.Table("notification").Where("notification.ctype=? and notification.uid=?", ctype, uid).Limit(limit, offset).Asc("notification.id").Find(rp)
			}
		}
	} else {

		if uid == 0 {
			if ctype != 0 {
				Engine.Table("notification").Where("notification.ctype=?", ctype).Limit(limit, offset).Desc("notification." + field).Find(rp)
			} else {
				Engine.Table("notification").Limit(limit, offset).Desc("notification." + field).Find(rp)
			}

		} else {

			if ctype == 0 {
				//Engine.Where("(ctype=-1 or ctype=1) and uid=?", uid).Limit(limit, offset).Desc(field).Find(rp)
				Engine.Table("notification").Where("notification.uid=?", uid).Limit(limit, offset).Desc("notification." + field).Find(rp)

			} else {

				Engine.Table("notification").Where("notification.ctype=? and notification.uid=?", ctype, uid).Limit(limit, offset).Desc("notification." + field).Find(rp)
			}
		}
	}
	return rp
}

func GetNotificationsByUidJoinUser(uid int64, ctype int64, offset int, limit int, field string) *[]Notificationjuser {
	rp := &[]Notificationjuser{}
	if field == "asc" {

		if uid == 0 {
			if ctype != 0 {
				Engine.Table("notification").Where("notification.ctype=?", ctype).Limit(limit, offset).Asc("notification.id").Join("LEFT", "user", "user.username = notification.author").Find(rp)
			} else {
				Engine.Table("notification").Limit(limit, offset).Asc("notification.id").Join("LEFT", "user", "user.username = notification.author").Find(rp)
			}

		} else {

			if ctype == 0 {
				Engine.Table("notification").Where("notification.uid=?", uid).Limit(limit, offset).Asc("notification.id").Join("LEFT", "user", "user.username = notification.author").Find(rp)

			} else {

				Engine.Table("notification").Where("notification.ctype=? and notification.uid=?", ctype, uid).Limit(limit, offset).Asc("notification.id").Join("LEFT", "user", "user.username = notification.author").Find(rp)
			}
		}
	} else {

		if uid == 0 {
			if ctype != 0 {
				Engine.Table("notification").Where("notification.ctype=?", ctype).Limit(limit, offset).Desc("notification."+field).Join("LEFT", "user", "user.username = notification.author").Find(rp)
			} else {
				Engine.Table("notification").Limit(limit, offset).Desc("notification."+field).Join("LEFT", "user", "user.username = notification.author").Find(rp)
			}

		} else {

			if ctype == 0 {
				//Engine.Where("(ctype=-1 or ctype=1) and uid=?", uid).Limit(limit, offset).Desc(field).Find(rp)
				Engine.Table("notification").Where("notification.uid=?", uid).Limit(limit, offset).Desc("notification."+field).Join("LEFT", "user", "user.username = notification.author").Find(rp)

			} else {

				Engine.Table("notification").Where("notification.ctype=? and notification.uid=?", ctype, uid).Limit(limit, offset).Desc("notification."+field).Join("LEFT", "user", "user.username = notification.author").Find(rp)
			}
		}
	}
	return rp
}

func GetNotificationsByUidUsernameJoinTopic(uid int64, author string, ctype int64, offset int, limit int, field string) *[]Notificationjtopic {
	rp := &[]Notificationjtopic{}

	if uid == 0 {
		if ctype != 0 {
			Engine.Table("notification").Where("notification.ctype=? and notification.author=?", ctype, author).Limit(limit, offset).Desc("notification."+field).Join("LEFT", "topic", "notification.uid = topic.id").Find(rp)
		} else {
			Engine.Table("notification").Where("notification.author=?", author).Limit(limit, offset).Desc("notification."+field).Join("LEFT", "topic", "notification.uid = topic.id").Find(rp)
		}

	} else {

		if ctype == 0 {
			Engine.Table("notification").Where("notification.uid=? and notification.author=?", uid, author).Limit(limit, offset).Desc("notification."+field).Join("LEFT", "topic", "notification.uid = topic.id").Find(rp)

		} else {

			Engine.Table("notification").Where("notification.ctype=? and notification.uid=? and notification.author=?", ctype, uid, author).Limit(limit, offset).Desc("notification."+field).Join("LEFT", "topic", "notification.uid = topic.id").Find(rp)
		}
	}
	return rp
}

func PostNotification(tid int64, rp *Notification) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	notificationid := int64(0)
	if _, err := sess.Insert(rp); err == nil && rp != nil {
		if rp.Tid > 0 {
			n, _ := sess.Where("tid=?", rp.Tid).Count(&Notification{})

			if row, err := sess.Table(&Topic{}).Where("id=?", rp.Tid).Update(&map[string]interface{}{"author": rp.Author, "notification_time": time.Now().Unix(), "notification_count": n, "notification_last_user_id": rp.Uid}); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PostNotification更新topic表相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		if rp.Uid > 0 {
			n, _ := sess.Where("uid=?", rp.Uid).Count(&Notification{})

			if row, err := sess.Table(&User{}).Where("id=?", rp.Uid).Update(&map[string]interface{}{"notification_count": n}); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PostNotification更新user表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		notificationid = rp.Id
	} else {
		sess.Rollback()
		return -1, err
	}

	// 提交事务
	return notificationid, sess.Commit()
}

func PutNotification(rid int64, rp *Notification) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	notificationid := int64(0)
	if row, err := sess.Update(rp, &Notification{Id: rid}); err != nil || row <= 0 {
		sess.Rollback()
		return -1, err
	} else {
		if rp.Uid > 0 {
			n, _ := sess.Where("uid=?", rp.Uid).Count(&Notification{})

			if row, err := sess.Table(&User{}).Where("id=?", rp.Uid).Update(&map[string]interface{}{"notification_count": n}); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PutNotification更新user表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}
	}

	// 提交事务
	return notificationid, sess.Commit()
}

func GetAllNotification() *[]Notification {
	rps := &[]Notification{}
	Engine.Desc("id").Find(rps)
	return rps
}

func GetNotification(id int64) (*Notification, error) {

	rpy := &Notification{}
	has, err := Engine.Id(id).Get(rpy)
	if has {
		return rpy, err
	} else {

		return nil, err
	}
}

func GetNotificationsByTid(tid int64, ctype int64, offset int, limit int, field string) *[]Notification {
	rp := &[]Notification{}

	//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
	if tid == 0 {
		if ctype != 0 {
			Engine.Where("ctype=?", ctype).Limit(limit, offset).Desc(field).Find(rp)
		} else {
			Engine.Limit(limit, offset).Desc(field).Find(rp)
		}

	} else {

		if ctype == 0 {
			//Engine.Where("(ctype=-1 or ctype=1) and tid=?", tid).Limit(limit, offset).Desc(field).Find(rp)
			Engine.Where("tid=?", tid).Limit(limit, offset).Desc(field).Find(rp)

		} else {

			Engine.Where("ctype=? and tid=?", ctype, tid).Limit(limit, offset).Desc(field).Find(rp)
		}
	}
	return rp
}

func GetNotificationsByTidUid(tid int64, uid int64, ctype int64, offset int, limit int, field string) *[]Notification {
	rp := &[]Notification{}

	//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
	if tid == 0 {
		if ctype != 0 {
			Engine.Where("ctype=? and uid=?", ctype, uid).Limit(limit, offset).Desc(field).Find(rp)
		} else {
			Engine.Where("uid=?", uid).Limit(limit, offset).Desc(field).Find(rp)
		}

	} else {

		if ctype == 0 {
			//Engine.Where("(ctype=-1 or ctype=1) and tid=?", tid).Limit(limit, offset).Desc(field).Find(rp)
			Engine.Where("tid=? and uid=?", tid, uid).Limit(limit, offset).Desc(field).Find(rp)

		} else {

			Engine.Where("ctype=? and tid=? and uid=?", ctype, tid, uid).Limit(limit, offset).Desc(field).Find(rp)
		}
	}
	return rp
}

func GetNotificationsByTidUsername(tid int64, author string, ctype int64, offset int, limit int, field string) *[]Notification {
	rp := &[]Notification{}

	if tid == 0 {
		if ctype != 0 {
			Engine.Where("ctype=? and author=?", ctype, author).Limit(limit, offset).Desc(field).Find(rp)
		} else {
			Engine.Where("author=?", author).Limit(limit, offset).Desc(field).Find(rp)
		}

	} else {

		if ctype == 0 {
			//Engine.Where("(ctype=-1 or ctype=1) and tid=?", tid).Limit(limit, offset).Desc(field).Find(rp)
			Engine.Where("tid=? and author=?", tid, author).Limit(limit, offset).Desc(field).Find(rp)

		} else {

			Engine.Where("ctype=? and tid=? and author=?", ctype, tid, author).Limit(limit, offset).Desc(field).Find(rp)
		}
	}
	return rp
}

func SetNotificationCountByUid(uid int64) (int64, error) {
	n, _ := Engine.Where("uid=?", uid).Count(&Notification{Uid: uid})

	qs := &User{}
	qs.NotificationCount = n
	affected, err := Engine.Where("id=?", uid).Cols("notification_count").Update(qs)
	return affected, err
}

func SetNotificationCount(uid int64, count int64) (int64, error) {

	qs := &User{}
	qs.NotificationCount = count
	affected, err := Engine.Where("id=?", uid).Cols("notification_count").Update(qs)
	return affected, err
}

func GetNotificationCountByUid(uid int64) int64 {
	n, _ := Engine.Where("uid=?", uid).Count(&Notification{Uid: uid})
	return n
}

func AddNotification(tid, rid, uid, ctype int64, subject, reply, author, avatar, avatarLarge, avatarMedium, avatarSmall string) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	// 执行事务
	notificationid := int64(0)
	{
		nf := &Notification{}
		nf.Tid = tid
		nf.Rid = rid
		nf.Uid = uid
		nf.Ctype = ctype
		nf.Subject = subject
		nf.Reply = reply
		nf.Author = author
		nf.Avatar = avatar
		nf.AvatarLarge = avatarLarge
		nf.AvatarMedium = avatarMedium
		nf.AvatarSmall = avatarSmall
		nf.Created = time.Now().Unix()

		if _, err := sess.Insert(nf); err == nil {
			//更新用户中的通知数量记录

			if nf.Uid > 0 {
				n, _ := sess.Where("uid=?", nf.Uid).Count(&Notification{})
				if affected, err := sess.Table(&User{}).Where("id=?", nf.Uid).Update(&map[string]interface{}{"notification_count": n}); err != nil || affected <= 0 {
					sess.Rollback()
					return -1, errors.New(fmt.Sprint("AddNotification #", nf.Id, "更新user表相关信息出现错误,affected:", affected, "错误:", err))
				}
			}

			notificationid = nf.Id
		} else {
			return -1, err
		}

	}

	// 提交事务
	return notificationid, sess.Commit()
}

func SetNotificationContentByRid(rid int64, Content string) error {
	if row, err := Engine.Table(&Notification{}).Where("id=?", rid).Update(&map[string]interface{}{"content": Content}); err != nil || row == 0 {
		return err
	} else {
		return nil
	}

}

func DelNotification(rid int64) error {
	if row, err := Engine.Id(rid).Delete(new(Notification)); err != nil || row == 0 {

		return errors.New("删除通知错误!")
	} else {
		return nil
	}

}

func DelNotificationByRole(rid int64, uid int64, role int64) error {
	allow := bool(false)
	if anz, err := GetNotification(rid); err == nil && anz != nil {
		if anz.Uid == uid {
			allow = true
		} else if role < 0 {
			allow = true
		}
		if allow {
			if row, err := Engine.Id(rid).Delete(new(Notification)); err != nil || row == 0 {
				return errors.New("删除通知发生错误!")
			} else {
				return nil
			}
		} else {
			return errors.New("你没有权限删除通知!")
		}

	} else {
		return errors.New("没有办法删除根本不存在的通知!")
	}

}

func DelNotificationsByPid(pid int64) error {
	rpy := &[]Notification{}
	if err := Engine.Where("pid=?", pid).Find(rpy); err == nil && rpy != nil {
		for _, v := range *rpy {
			if err := DelNotificationByRole(v.Id, v.Uid, -1000); err != nil {
				fmt.Println("DelNotificationByRole:", err)
			}
		}
		return nil
	} else {
		return err
	}

}

func DelNotificationsByTid(tid int64) error {

	rpy := &[]Notification{}
	if err := Engine.Where("tid=?", tid).Find(rpy); err == nil && rpy != nil {
		for _, v := range *rpy {
			if err := DelNotificationByRole(v.Id, v.Uid, -1000); err != nil {
				fmt.Println("DelNotificationByRole:", err)
			}
		}
		return nil
	} else {
		return err
	}

}
