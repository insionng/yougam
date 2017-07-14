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

type Reply struct {
	Id                int64
	Pid               int64 `xorm:"index"` //上级Reply id
	Uid               int64 `xorm:"index"`
	Tid               int64 `xorm:"index"`
	Sort              int64
	Ctype             int64   `xorm:"index"`
	Content           string  `xorm:"text"`
	Tailinfo          int64   `xorm:"index"`
	Attachment        string  `xorm:"text"`
	Avatar            string  `xorm:"index"` //200x200
	AvatarLarge       string  `xorm:"index"` //100x100
	AvatarMedium      string  `xorm:"index"` //48x48
	AvatarSmall       string  `xorm:"index"` //32x32
	Created           int64   `xorm:"created index"`
	Updated           int64   `xorm:"updated"`
	Hotness           float64 `xorm:"index"`
	Confidence        float64 `xorm:"index"`
	Hotup             int64   `xorm:"index"`
	Hotdown           int64   `xorm:"index"`
	Hotscore          int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64   `xorm:"index"` //Hotup  + Hotdown
	Views             int64
	Author            string `xorm:"index"`
	AuthorSignature   string `xorm:"index"`
	Email             string `xorm:"index"`
	Website           string `xorm:"index"`
	ReplyTime         int64
	ReplyCount        int64
	ReplyLastUserId   int64
	ReplyLastUsername string
	ReplyLastNickname string
	ReplyLastTopic    string  //最后回复的话题之标题
	Latitude          float64 `xorm:"index"`
	Longitude         float64 `xorm:"index"`
}

type ReplyMark struct {
	Id  int64
	Uid int64 `xorm:"index"`
	Rid int64 `xorm:"index"` //Reply id
}

type Replyjuser struct {
	Reply `xorm:"extends"`
	User  `xorm:"extends"`
}

type Replyjtopic struct {
	Reply `xorm:"extends"`
	Topic `xorm:"extends"`
}

func GetReplysByPid4Count(pid int64, offset int, limit int, ctype int64) (int64, error) {

	if ctype != 0 {
		return Engine.Where("(pid=? or id=?) and ctype=?", pid, pid, ctype).Limit(limit, offset).Count(&Reply{})
	} else {
		return Engine.Where("pid=? or id=?", pid, pid).Limit(limit, offset).Count(&Reply{})
	}

}

func SetReplyMark(uid int64, rid int64) (int64, error) {

	rpm := &ReplyMark{Uid: uid, Rid: rid}
	rows, err := Engine.Insert(rpm)
	return rows, err
}

func DelReplyMark(uid int64, rid int64) {
	rpm := new(ReplyMark)
	Engine.Where("uid=? and rid=?", uid, rid).Delete(rpm)

}

func ReplyMarkCount(rid int64) (int64, error) {
	return Engine.Where("rid=?", rid).Count(&ReplyMark{})
}

func ReplyMarkCountByUid(uid int64) (int64, error) {
	return Engine.Where("uid=?", uid).Count(&ReplyMark{})
}

func IsReplyMark(uid int64, rid int64) bool {

	rpm := new(ReplyMark)

	if has, err := Engine.Where("uid=? and rid=?", uid, rid).Get(rpm); err != nil {
		return false
	} else {
		if has {
			return (rpm.Uid == uid)
		} else {
			return false
		}
	}

}

func GetReplysByTidUsernameJoinTopic(tid int64, author string, ctype int64, offset int, limit int, field string) *[]*Replyjtopic {

	rp := new([]*Replyjtopic)

	//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
	if tid == 0 {
		if ctype != 0 {
			Engine.Table("reply").Where("reply.ctype=? and reply.author=?", ctype, author).Limit(limit, offset).Desc("reply."+field).Join("LEFT", "topic", "reply.tid = topic.id").Find(rp)
		} else {
			Engine.Table("reply").Where("reply.author=?", author).Limit(limit, offset).Desc("reply."+field).Join("LEFT", "topic", "reply.tid = topic.id").Find(rp)
		}

	} else {

		if ctype == 0 {
			Engine.Table("reply").Where("reply.tid=? and reply.author=?", tid, author).Limit(limit, offset).Desc("reply."+field).Join("LEFT", "topic", "reply.tid = topic.id").Find(rp)

		} else {

			Engine.Table("reply").Where("reply.ctype=? and reply.tid=? and reply.author=?", ctype, tid, author).Limit(limit, offset).Desc("reply."+field).Join("LEFT", "topic", "reply.tid = topic.id").Find(rp)
		}
	}
	return rp
}

func GetSubReplys(pid int64, offset int, limit int, field string) (*[]*Reply, error) {
	rps := &[]*Reply{}
	err := Engine.Where("pid=?", pid).Limit(limit, offset).Desc(field).Find(rps)
	return rps, err
}

func PostReply(tid int64, rp *Reply) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	replyid := int64(0)
	if _, err := sess.Insert(rp); err == nil && rp != nil {
		if rp.Tid > 0 {
			n, _ := sess.Where("tid=?", rp.Tid).Count(&Reply{})

			if row, err := sess.Table(&Topic{}).Where("id=?", rp.Tid).Update(&map[string]interface{}{"author": rp.Author, "reply_time": time.Now().Unix(), "reply_count": n, "reply_last_user_id": rp.Uid}); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PostReply更新topic表相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		if rp.Uid > 0 {
			n, _ := sess.Where("uid=?", rp.Uid).Count(&Reply{})

			if row, err := sess.Table(&User{}).Where("id=?", rp.Uid).Update(&map[string]interface{}{"reply_time": time.Now().Unix(), "reply_count": n, "reply_last_tid": rp.Tid, "reply_last_topic": rp.ReplyLastTopic}); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PostReply更新user表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		replyid = rp.Id
	} else {
		sess.Rollback()
		return -1, err
	}

	// 提交事务
	return replyid, sess.Commit()
}

func PutReply(rid int64, rp *Reply) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	replyid := int64(0)
	if row, err := sess.Update(rp, &Reply{Id: rid}); err != nil || row <= 0 {
		sess.Rollback()
		return -1, err
	} else {
		if rp.Uid > 0 {
			n, _ := sess.Where("uid=?", rp.Uid).Count(&Reply{})

			if row, err := sess.Table(&User{}).Where("id=?", rp.Uid).Update(&map[string]interface{}{"reply_time": time.Now().Unix(), "reply_count": n, "reply_last_tid": rp.Tid, "reply_last_topic": rp.ReplyLastTopic}); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PutReply更新user表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}
	}

	// 提交事务
	return replyid, sess.Commit()
}

func GetAllReply() *[]*Reply {
	rps := &[]*Reply{}
	Engine.Desc("id").Find(rps)
	return rps
}

func GetReply(id int64) (*Reply, error) {

	rpy := &Reply{}
	has, err := Engine.Id(id).Get(rpy)
	if has {
		return rpy, err
	} else {

		return nil, err
	}
}

func GetReplysByPid(pid int64, ctype int64, offset int, limit int, field string) *[]*Reply {
	rp := &[]*Reply{}

	//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
	if pid == 0 {
		if ctype == 0 {
			Engine.Limit(limit, offset).Desc(field).Find(rp)
		} else {
			Engine.Where("ctype=?", ctype).Limit(limit, offset).Desc(field).Find(rp)
		}
	} else {

		if ctype == 0 {
			//Engine.Where("(ctype=-1 or ctype=1) and pid=?", pid).Limit(limit, offset).Desc(field).Find(rp)
			Engine.Where("pid=?", pid).Limit(limit, offset).Desc(field).Find(rp)

		} else {

			Engine.Where("ctype=? and pid=?", ctype, pid).Limit(limit, offset).Desc(field).Find(rp)
		}
	}
	return rp
}

func GetReplysByTid(tid int64, ctype int64, offset int, limit int, field string) *[]*Reply {
	rp := &[]*Reply{}
	if field == "asc" {

		//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
		if tid == 0 {
			if ctype != 0 {
				Engine.Table("reply").Where("reply.ctype=?", ctype).Limit(limit, offset).Asc("reply.id").Find(rp)
			} else {
				Engine.Table("reply").Limit(limit, offset).Asc("reply.id").Find(rp)
			}

		} else {

			if ctype == 0 {
				Engine.Table("reply").Where("reply.tid=?", tid).Limit(limit, offset).Asc("reply.id").Find(rp)

			} else {

				Engine.Table("reply").Where("reply.ctype=? and reply.tid=?", ctype, tid).Limit(limit, offset).Asc("reply.id").Find(rp)
			}
		}
	} else {

		//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
		if tid == 0 {
			if ctype != 0 {
				Engine.Table("reply").Where("reply.ctype=?", ctype).Limit(limit, offset).Desc("reply." + field).Find(rp)
			} else {
				Engine.Table("reply").Limit(limit, offset).Desc("reply." + field).Find(rp)
			}

		} else {

			if ctype == 0 {
				//Engine.Where("(ctype=-1 or ctype=1) and tid=?", tid).Limit(limit, offset).Desc(field).Find(rp)
				Engine.Table("reply").Where("reply.tid=?", tid).Limit(limit, offset).Desc("reply." + field).Find(rp)

			} else {

				Engine.Table("reply").Where("reply.ctype=? and reply.tid=?", ctype, tid).Limit(limit, offset).Desc("reply." + field).Find(rp)
			}
		}
	}
	return rp

}

func GetReplysByTidJoinUser(tid int64, ctype int64, offset int, limit int, field string) *[]*Replyjuser {
	rp := &[]*Replyjuser{}
	if field == "asc" {

		//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
		if tid == 0 {
			if ctype != 0 {
				Engine.Table("reply").Where("reply.ctype=?", ctype).Limit(limit, offset).Asc("reply.id").Join("LEFT", "user", "user.id = reply.uid").Find(rp)
			} else {
				Engine.Table("reply").Limit(limit, offset).Asc("reply.id").Join("LEFT", "user", "user.id = reply.uid").Find(rp)
			}

		} else {

			if ctype == 0 {
				Engine.Table("reply").Where("reply.tid=?", tid).Limit(limit, offset).Asc("reply.id").Join("LEFT", "user", "user.id = reply.uid").Find(rp)

			} else {

				Engine.Table("reply").Where("reply.ctype=? and reply.tid=?", ctype, tid).Limit(limit, offset).Asc("reply.id").Join("LEFT", "user", "user.id = reply.uid").Find(rp)
			}
		}
	} else {

		//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
		if tid == 0 {
			if ctype != 0 {
				Engine.Table("reply").Where("reply.ctype=?", ctype).Limit(limit, offset).Desc("reply."+field).Join("LEFT", "user", "user.id = reply.uid").Find(rp)
			} else {
				Engine.Table("reply").Limit(limit, offset).Desc("reply."+field).Join("LEFT", "user", "user.id = reply.uid").Find(rp)
			}

		} else {

			if ctype == 0 {
				//Engine.Where("(ctype=-1 or ctype=1) and tid=?", tid).Limit(limit, offset).Desc(field).Find(rp)
				Engine.Table("reply").Where("reply.tid=?", tid).Limit(limit, offset).Desc("reply."+field).Join("LEFT", "user", "user.id = reply.uid").Find(rp)

			} else {

				Engine.Table("reply").Where("reply.ctype=? and reply.tid=?", ctype, tid).Limit(limit, offset).Desc("reply."+field).Join("LEFT", "user", "user.id = reply.uid").Find(rp)
			}
		}
	}
	return rp

}

func GetReplysByTidUid(tid int64, uid int64, ctype int64, offset int, limit int, field string) *[]*Reply {
	rp := &[]*Reply{}

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

func GetReplysByTidUsername(tid int64, author string, ctype int64, offset int, limit int, field string) *[]*Reply {
	rp := &[]*Reply{}

	//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
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

func SetReplyCountByPid(qid int64) (int64, error) {
	n, _ := Engine.Where("pid=?", qid).Count(&Reply{Pid: qid})

	qs := &Topic{}
	qs.ReplyCount = n
	affected, err := Engine.Where("id=?", qid).Cols("reply_count").Update(qs)
	return affected, err
}

func SetReplyCountByTid(tid int64) (int64, error) {
	n, _ := Engine.Where("tid=?", tid).Count(&Reply{Tid: tid})

	qs := &Topic{}
	qs.ReplyCount = n
	affected, err := Engine.Where("id=?", tid).Cols("reply_count").Update(qs)
	return affected, err
}

func GetReplyCountByPid(pid int64) int64 {
	n, _ := Engine.Where("pid=?", pid).Count(&Reply{Pid: pid})
	return n
}

func GetReplyCountByTid(tid int64) int64 {
	n, _ := Engine.Where("tid=?", tid).Count(&Reply{Tid: tid})
	return n + 1
}

func SetReply(rid int64, rp *Reply) (int64, error) {
	rp.Id = rid
	return Engine.Insert(rp)
}

func AddReply(tid, pid, uid, ctype int64, content, attachment, author, avatar, avatarLarge, avatarMedium, avatarSmall, author_signature, nickname, email, website string) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	// 执行事务
	replyid := int64(0)
	{
		rp := new(Reply)
		rp.Tid = tid
		rp.Pid = pid
		rp.Uid = uid
		rp.Ctype = ctype
		rp.Content = content
		rp.Attachment = attachment
		rp.Author = author
		rp.AuthorSignature = author_signature
		rp.Email = email
		rp.Website = website
		rp.Avatar = avatar
		rp.AvatarLarge = avatarLarge
		rp.AvatarMedium = avatarMedium
		rp.AvatarSmall = avatarSmall
		rp.Created = time.Now().Unix()
		rp.ReplyTime = time.Now().Unix()

		if _, err := sess.Insert(rp); err == nil {
			//更新话题中的回应相关记录

			tp := new(Topic)
			if has, err := sess.Id(tid).Get(tp); has {
				tp.ReplyTime = time.Now().Unix()
				ReplyCount, _ := sess.Where("tid=?", tid).Count(&Reply{Tid: tid})
				tp.ReplyCount = ReplyCount + 1
				tp.ReplyLastUserId = uid
				tp.ReplyLastUsername = author
				tp.ReplyLastNickname = nickname
				//if row, err := sess.Table(&Topic{}).Where("id=?", tid).Update(&map[string]interface{}{"reply_time": time.Now().Unix(), "reply_count": ReplyCount, "reply_last_user_id": uid, "reply_last_username": author, "reply_last_nickname": nickname}); err != nil || row <= 0 {
				if row, err := sess.Table(&Topic{}).Where("id=?", tid).Update(tp); err != nil || row <= 0 {
					sess.Rollback()
					return -1, errors.New(fmt.Sprint("AddReply #", rp.Id, "更新topic表相关信息出现错误,row:", row, "错误:", err))
				}
			} else {
				sess.Rollback()
				return -1, err
			}

			if rp.Uid > 0 {
				n, _ := sess.Where("uid=?", rp.Uid).Count(&Reply{})

				if row, err := sess.Table(&User{}).Where("id=?", rp.Uid).Update(&map[string]interface{}{"reply_time": time.Now().Unix(), "reply_count": n, "reply_last_tid": rp.Tid, "reply_last_topic": rp.ReplyLastTopic}); err != nil || row == 0 {
					sess.Rollback()
					return -1, errors.New(fmt.Sprint("AddReply #", rp.Id, "更新user表相关信息出现错误,row:", row, "错误:", err))
				}
			}

			replyid = rp.Id
		}

	}

	// 提交事务
	return replyid, sess.Commit()
}

func SetReplyContentByRid(rid int64, Content string) error {
	if row, err := Engine.Table(&Reply{}).Where("id=?", rid).Update(&map[string]interface{}{"content": Content}); err != nil || row == 0 {
		log.Println("SetReplyContentByRid  row:", row, "SetReplyContentByRid出现错误:", err)
		return err
	} else {
		return nil
	}

}

func DelReply(rid int64) error {
	if row, err := Engine.Id(rid).Delete(new(Reply)); err != nil || row == 0 {
		log.Println("row:", row, "删除回应错误:", err)
		return errors.New("删除回应错误!")
	} else {
		return nil
	}

}

func DelReplyByRole(rid int64, uid int64, role int64) error {
	allow := bool(false)
	if anz, err := GetReply(rid); err == nil && anz != nil {
		if anz.Uid == uid {
			allow = true
		} else if role < 0 {
			allow = true
		}
		if allow {
			if row, err := Engine.Id(rid).Delete(new(Reply)); err != nil || row == 0 {
				log.Println("row:", row, "删除回复发生错误:", err)
				return errors.New("删除回复发生错误!")
			} else {
				return nil
			}
		} else {
			return errors.New("你没有权限删除回复!")
		}

	} else {
		return errors.New("没有办法删除根本不存在的回复!")
	}

}

func DelReplysByPid(pid int64) error {
	rpy := &[]*Reply{}
	if err := Engine.Where("pid=?", pid).Find(rpy); err == nil && rpy != nil {
		for _, v := range *rpy {
			if err := DelReplyByRole(v.Id, v.Uid, -1000); err != nil {
				log.Println("DelReplyByRole:", err)
			}
		}
		return nil
	} else {
		return err
	}

}

func DelReplysByTid(tid int64) error {

	rpy := &[]*Reply{}
	if err := Engine.Where("tid=?", tid).Find(rpy); err == nil && rpy != nil {
		for _, v := range *rpy {
			if err := DelReplyByRole(v.Id, v.Uid, -1000); err != nil {
				log.Println("DelReplyByRole:", err)
			}

			//检查附件字段并尝试删除文件
			if len(v.Attachment) > 0 {
				allow := true //设为允许管理员权限 暂时不细分
				if p := helper.URL2local(v.Attachment); helper.Exist(p) {
					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误，但不要反回错误，以免陷入死循环无法删掉
							log.Println("ROOT DEL Reply Attachment, Reply ID:", v.Id, ",ERR:", err)
						}
					} else { //检查用户对文件的所有权
						if helper.VerifyUserfile(p, strconv.FormatInt(v.Uid, 10)) {
							if err := os.Remove(p); err != nil {
								log.Println("DEL Reply Attachment, Reply ID:", v.Id, ",ERR:", err)
							}
						}
					}
				} /*else {

					////删除七牛云存储中的图片

					rsCli := rs.New(nil)

					//2014-8-21-134506AB4F24A56EF60543.jpg,2014-8-21-134506AB4F24A56EF60543.jpg
					imgkeys := strings.Split(v.Attachment, ",")
					for _, v := range imgkeys {
						tmpkey := strings.Split(v, ".")
						delkey := tmpkey[0] //key是32位  不含后缀
						rsCli.Delete(nil, BUCKET, delkey)
					}

				}
				*/

			}
		}
		return nil
	} else {
		return err
	}

}
