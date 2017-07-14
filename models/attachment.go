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

type Attachment struct {
	Id            int64
	Pid           int64   `xorm:"index"` //为0代表没有上级，本身就是顶层附件 大于0则本身是子附件，而该数字代表上级附件的id
	Cid           int64   `xorm:"index"`
	Nid           int64   `xorm:"index"` //nodeid
	Uid           int64   `xorm:"index"`
	Sort          int64   `xorm:"index"` //排序字段 需要手动对附件排序置顶等操作时使用
	Ctype         int64   `xorm:"index"` //ctype作用在于区分附件的类型
	Content       string  `xorm:"text"`  //资源URL
	Tags          string  `xorm:"index"`
	Created       int64   `xorm:"created index"`
	Updated       int64   `xorm:"updated"`
	Hotness       float64 `xorm:"index"`
	Confidence    float64 `xorm:"index"` //信任度数值
	Hotup         int64   `xorm:"index"`
	Hotdown       int64   `xorm:"index"`
	Hotscore      int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote       int64   `xorm:"index"` //Hotup  + Hotdown
	Views         int64
	Author        string `xorm:"index"`
	Category      string `xorm:"index"`
	Node          string `xorm:"index"` //nodename
	FavoriteCount int64
	Latitude      float64 `xorm:"index"`   //纬度
	Longitude     float64 `xorm:"index"`   //经度
	Version       int64   `xorm:"version"` //乐观锁
}

func AddAttachment(content string, pid, cid, nid, uid int64) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	attachmentid := int64(0)
	{
		cat := &Category{}
		if cid > 0 {
			if has, err := sess.Where("id=?", cid).Get(cat); (err != nil) || (has == false) {
				sess.Rollback()
				return -1, err
			}
		}

		nd := &Node{}
		if nid > 0 {
			if has, err := sess.Where("id=?", nid).Get(nd); (err != nil) || (has == false) {
				sess.Rollback()
				return -1, err
			}
		}

		usr := &User{}
		if uid > 0 {
			if has, err := sess.Where("id=?", uid).Get(usr); (err != nil) || (has == false) {
				sess.Rollback()
				return -1, err
			}
		}

		obj := new(Attachment)
		obj.Pid = pid
		obj.Cid = cat.Id
		obj.Nid = nid
		obj.Uid = uid
		obj.Content = content
		obj.Author = usr.Username
		obj.Category = cat.Title
		obj.Node = nd.Title
		obj.Created = time.Now().Unix()

		if row, err := sess.Insert(obj); (err == nil) && (row > 0) {
			if obj.Uid > 0 {
				n, _ := sess.Where("uid=?", obj.Uid).Count(&Attachment{})
				_u := map[string]interface{}{"attachment_time": time.Now().Unix(), "attachment_count": n, "attachment_last_nid": obj.Nid, "attachment_last_node": obj.Node}
				if row, err := sess.Table(&User{}).Where("id=?", obj.Uid).Update(&_u); err != nil || row <= 0 {

					return -1, errors.New(fmt.Sprint("AddAttachment更新user表附件相关信息时:出现", row, "行影响,错误:", err))

				}
			}
			attachmentid = obj.Id
		} else {
			sess.Rollback()
			return -1, err
		}

	}

	// 提交事务
	return attachmentid, sess.Commit()
}

func SetAttachment(tid int64, obj *Attachment) (int64, error) {
	obj.Id = tid
	return Engine.Insert(obj)
}

func GetAttachmentsByHotnessNodes(nodelimit int, attachmentlimit int) []*[]*Attachment {
	//找出最热的节点:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级附件
	nds, _ := GetNodes(0, nodelimit, "hotness")
	attachments := make([]*[]*Attachment, 0)

	if len(*nds) > 0 {
		i := 0
		for _, v := range *nds {
			i = i + 1
			objs := GetAttachmentsByNid(v.Id, 0, attachmentlimit, 0, "views")
			if len(*objs) != 0 {
				attachments = append(attachments, objs)
			}
			if i == len(*nds)-1 {
				break
			}
		}
	}

	return attachments

}

func GetAttachmentsByScoreNodes(nodelimit int, attachmentlimit int) []*[]*Attachment {
	//找出最热的节点:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级附件
	nds, _ := GetNodes(0, nodelimit, "hotscore")
	attachments := make([]*[]*Attachment, 0)

	if len(*nds) > 0 {
		i := 0
		for _, v := range *nds {
			i = i + 1
			objs := GetAttachmentsByNid(v.Id, 0, attachmentlimit, 0, "views")
			if len(*objs) != 0 {
				attachments = append(attachments, objs)
			}
			if i == len(*nds) {
				break
			}
		}
	}

	return attachments

}

func GetAttachmentsByHotnessCategory(catlimit int, attachmentlimit int) []*[]*Attachment {
	//找出最热的分类:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级附件
	cats, _ := GetCategories(0, catlimit, "hotness")
	attachments := make([]*[]*Attachment, 0)

	if len(*cats) > 0 {
		i := 0
		for _, v := range *cats {
			i = i + 1
			//(cid int64, offset int, limit int, ctype int64, field string)
			objs := GetAttachmentsByCid(v.Id, 0, attachmentlimit, 0, "views")
			if len(*objs) != 0 {
				attachments = append(attachments, objs)
			}
			if i == len(*cats) {
				break
			}
		}
	}

	return attachments

}

func GetAttachment(id int64) (*Attachment, error) {
	obj := new(Attachment)
	has, err := Engine.Id(id).Get(obj)
	if has {
		return obj, err
	} else {
		return nil, err
	}
}

func GetAttachments(offset int, limit int, field string) (*[]*Attachment, error) {
	objs := new([]*Attachment)
	err := Engine.Limit(limit, offset).Desc(field).Find(objs)
	return objs, err
}

func GetSubAttachments(pid int64, offset int, limit int, field string) (*[]*Attachment, error) {
	objs := new([]*Attachment)
	err := Engine.Where("pid=?", pid).Limit(limit, offset).Desc(field).Find(objs)
	return objs, err
}

func GetAttachmentsCount(offset int, limit int) (int64, error) {
	total, err := Engine.Limit(limit, offset).Count(&Attachment{})
	return total, err
}

func GetAttachmentsByPid4Count(pid int64, offset int, limit int, ctype int64) (int64, error) {

	if ctype != 0 {
		return Engine.Where("(pid=? or id=?) and ctype=?", pid, pid, ctype).Limit(limit, offset).Count(&Attachment{})
	} else {
		return Engine.Where("pid=? or id=?", pid, pid).Limit(limit, offset).Count(&Attachment{})
	}

}

func GetAttachmentsCountByNode(node string, offset int, limit int) (int64, error) {
	total, err := Engine.Where("node=?", node).Limit(limit, offset).Count(&Attachment{})
	return total, err
}

func GetAttachmentsByCategoryCount(category string, offset int, limit int, field string) (int64, error) {
	total, err := Engine.Where("category=?", category).Limit(limit, offset).Count(&Attachment{})
	return total, err
}

//GetAttachmentsByCid大数据下如出现性能问题 可以使用 GetAttachmentsByCidOnBetween
func GetAttachmentsByCid(cid int64, offset int, limit int, ctype int64, field string) *[]*Attachment {
	//排序首先是热值优先，然后是时间优先。
	objs := new([]*Attachment)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Limit(limit, offset).Asc("id").Find(objs)
		} else {
			Engine.Where("cid=?", cid).Limit(limit, offset).Asc("id").Find(objs)

		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Desc(field).Limit(limit, offset).Find(objs)

		} else {
			if cid == 0 {
				Engine.Desc(field).Limit(limit, offset).Find(objs)
			} else {
				Engine.Where("cid=?", cid).Desc(field).Limit(limit, offset).Find(objs)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
		} else {
			if cid == 0 {
				Engine.Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			} else {
				Engine.Where("cid=?", cid).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			}
		}

	}
	return objs
}

func GetAttachmentsByCidOnBetween(cid int64, startid int64, endid int64, offset int, limit int, ctype int64, field string) *[]*Attachment {
	objs := new([]*Attachment)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=? and id>? and id<?", cid, ctype, startid-1, endid+1).Limit(limit, offset).Asc("id").Find(objs)
		} else {
			if cid == 0 {
				Engine.Where("id>? and id<?", startid-1, endid+1).Limit(limit, offset).Asc("id").Find(objs)
			} else {
				Engine.Where("cid=? and id>? and id<?", cid, startid-1, endid+1).Limit(limit, offset).Asc("id").Find(objs)
			}
		}
	default: //Desc
		if ctype != 0 {
			Engine.Where("cid=? and ctype=? and id>? and id<?", cid, ctype, startid-1, endid+1).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
		} else {
			if cid == 0 {
				Engine.Where("id>? and id<?", startid-1, endid+1).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			} else {
				Engine.Where("cid=? and id>? and id<?", cid, startid-1, endid+1).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			}
		}
	}

	return objs
}

func GetAttachmentsByCategory(category string, offset int, limit int, ctype int64, field string) *[]*Attachment {
	//排序首先是热值优先，然后是时间优先。
	objs := new([]*Attachment)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Limit(limit, offset).Asc("id").Find(objs)
		} else {
			Engine.Where("category=?", category).Limit(limit, offset).Asc("id").Find(objs)

		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Desc(field).Limit(limit, offset).Find(objs)

		} else {
			if category == "" {
				Engine.Desc(field).Limit(limit, offset).Find(objs)
			} else {
				Engine.Where("category=?", category).Desc(field).Limit(limit, offset).Find(objs)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
		} else {
			if category == "" {
				Engine.Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			} else {
				Engine.Where("category=?", category).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			}
		}

	}
	return objs
}

//GetAttachmentsByUid不区分父子附件
func GetAttachmentsByUid(uid int64, offset int, limit int, ctype int64, field string) *[]*Attachment {
	//排序首先是热值优先，然后是时间优先。
	objs := new([]*Attachment)

	switch {
	case field == "asc":
		if uid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(field).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("uid=? and ctype=?", uid, ctype).Limit(limit, offset).Asc("id").Find(objs)
			} else {
				Engine.Where("uid=?", uid).Limit(limit, offset).Asc("id").Find(objs)
			}
		}
	default:
		if uid <= 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(field).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("uid=? and ctype=?", uid, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			} else {
				Engine.Where("uid=?", uid).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			}
		}
	}
	return objs
}

func GetAttachmentsByNid(nodeid int64, offset int, limit int, ctype int64, field string) *[]*Attachment {
	//排序首先是热值优先，然后是时间优先。
	objs := new([]*Attachment)

	switch {
	case field == "asc":
		if nodeid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(field).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("nid=? and ctype=?", nodeid, ctype).Limit(limit, offset).Asc("id").Find(objs)
			} else {
				Engine.Where("nid=?", nodeid).Limit(limit, offset).Asc("id").Find(objs)
			}
		}
	default:
		if nodeid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(field).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("nid=? and ctype=?", nodeid, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			} else {
				Engine.Where("nid=?", nodeid).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(objs)
			}
		}
	}
	return objs
}

func GetAttachmentsByPid(pid int64, offset int, limit int, ctype int64, field string) *[]*Attachment {

	objs := new([]*Attachment)

	switch {
	case field == "asc":
		{
			if pid > 0 { //即只查询单个附件的父级与子级，此时sort字段无效，所以无须作为条件参与排序
				if ctype != 0 {
					Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?) and attachment.ctype=?", pid, pid, ctype).Limit(limit, offset).Asc("attachment.id").Find(objs)
				} else {
					Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?)", pid, pid).Limit(limit, offset).Asc("attachment.id").Find(objs)
				}
			} else {
				if ctype != 0 {
					Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?) and attachment.ctype=? and attachment.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Asc("attachment.id").Find(objs)
				} else {
					Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?) and attachment.sort<=?", pid, pid, 0).Limit(limit, offset).Asc("attachment.id").Find(objs)
				}
			}

		}
	case field == "cold":
		{
			if ctype != 0 {
				Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?) and attachment.ctype=? and attachment.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Asc("attachment.views").Find(objs)
			} else {
				Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?) and attachment.sort<=?", pid, pid, 0).Limit(limit, offset).Asc("attachment.views").Find(objs)
			}
		}
	case field == "sort":
		{ //Desc("attachment.sort", "attachment.confidence") ，这段即是"attachment.sort"的优先级较"attachment.confidence"要高
			if ctype != 0 {
				Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?) and attachment.ctype=? and attachment.sort>?", pid, pid, ctype, 0).Limit(limit, offset).Desc("attachment.sort", "attachment.confidence").Find(objs)
			} else {
				Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?) and attachment.sort>?", pid, pid, 0).Limit(limit, offset).Desc("attachment.sort", "attachment.confidence").Find(objs)
			}
		}
	default:
		{
			if field == "desc" {
				field = "id"
			}
			if ctype != 0 {
				Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?) and attachment.ctype=? and attachment.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Desc("attachment."+field, "attachment.views", "attachment.reply_count", "attachment.created").Find(objs)
			} else {
				Engine.Table("attachment").Where("(attachment.pid=? or attachment.id=?)", pid, pid).And("attachment.sort<=?", 0).Limit(limit, offset).Desc("attachment."+field, "attachment.views", "attachment.reply_count", "attachment.created").Find(objs)
			}
		}
	}
	return objs
}

func GetAttachmentsByPidSinceCreated(pid int64, offset int, limit int, ctype int64, field string, since int64) *[]*Attachment {

	switch objs := new([]*Attachment); {
	case field == "asc":
		{
			if ctype != 0 {
				Engine.Table("attachment").Where("attachment.created>? and (attachment.pid=? or attachment.id=?) and attachment.ctype=? and attachment.sort<=?", since, pid, pid, ctype, 0).Limit(limit, offset).Asc("attachment.id").Find(objs)
			} else {
				Engine.Table("attachment").Where("attachment.created>? and (attachment.pid=? or attachment.id=?) and attachment.sort<=?", since, pid, pid, 0).Limit(limit, offset).Asc("attachment.id").Find(objs)
			}
			return objs
		}
	default:
		{
			if ctype != 0 {
				Engine.Table("attachment").Where("attachment.created>? and (attachment.pid=? or attachment.id=?) and attachment.ctype=? and attachment.sort<=?", since, pid, pid, ctype, 0).Limit(limit, offset).Desc("attachment."+field, "attachment.views", "attachment.reply_count", "attachment.created").Find(objs)
			} else {
				Engine.Table("attachment").Where("attachment.created>? and (attachment.pid=? or attachment.id=?) and attachment.sort<=?", since, pid, pid, 0).Limit(limit, offset).Desc("attachment."+field, "attachment.views", "attachment.reply_count", "attachment.created").Find(objs)
			}
			return objs
		}
	}

}

func GetAttachmentsByNode(node string, offset int, limit int, field string) (*[]*Attachment, error) {
	objs := new([]*Attachment)
	err := Engine.Where("node=?", node).Limit(limit, offset).Desc(field).Find(objs)
	return objs, err
}

//发布附件  返回 附件id,错误
func PostAttachment(obj *Attachment) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	attachmentid := int64(0)
	if _, err := sess.Insert(obj); err == nil && obj != nil {

		if obj.Nid > 0 {
			n, _ := sess.Where("nid=?", obj.Nid).Count(&Attachment{})
			_n := map[string]interface{}{"author": obj.Author, "attachment_time": time.Now().Unix(), "attachment_count": n, "attachment_last_user_id": obj.Uid}
			if row, err := sess.Table(&Node{}).Where("id=?", obj.Nid).Update(&_n); err != nil || row <= 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PostAttachment更新node表附件相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		if obj.Uid > 0 {
			n, _ := sess.Where("uid=?", obj.Uid).Count(&Attachment{})
			_u := map[string]interface{}{"attachment_time": time.Now().Unix(), "attachment_count": n, "attachment_last_nid": obj.Nid, "attachment_last_node": obj.Node}
			if row, err := sess.Table(&User{}).Where("id=?", obj.Uid).Update(&_u); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PostAttachment更新user表附件相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		attachmentid = obj.Id
	} else {
		// 发生错误时进行回滚
		sess.Rollback()
		return -1, err
	}

	// 提交事务
	return attachmentid, sess.Commit()
}

func PutAttachment(tid int64, obj *Attachment) (int64, error) {
	//覆盖式更新
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	attachmentid := int64(0)
	if row, err := sess.Update(obj, &Attachment{Id: tid}); err != nil || row <= 0 {
		sess.Rollback()
		return -1, err
	} else {
		if obj.Uid > 0 {
			n, _ := sess.Where("uid=?", obj.Uid).Count(&Attachment{})
			_u := map[string]interface{}{"attachment_time": time.Now().Unix(), "attachment_count": n, "attachment_last_nid": obj.Nid, "attachment_last_node": obj.Node}
			if row, err := sess.Table(&User{}).Where("id=?", obj.Uid).Update(&_u); err != nil || row <= 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PutAttachment更新user表附件相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		attachmentid = obj.Id
	}

	// 提交事务
	return attachmentid, sess.Commit()
}

func PutAttachmentViaVersion(tid int64, attachment *Attachment) (int64, error) {
	/*
		//乐观锁目前不支持map 以下句式无效！
			return Engine.Table(&Attachment{}).Where("id=?", tid).Update(&map[string]interface{}{
				"views":   views,
				"version": version,
			})
	*/
	//return Engine.Table(&Attachment{}).Where("id=?", tid).Update(attachment)
	return Engine.Id(tid).Update(attachment)
}

func PutViews2AttachmentViaVersion(tid int64, attachment *Attachment) (int64, error) {
	/*
		return Engine.Table(&Attachment{}).Where("id=?", tid).Update(&map[string]interface{}{
			"views": views,
		})
	*/
	return Engine.Id(tid).Cols("views").Update(attachment)
}

//func PutSort2Attachment(tid int64, sort int64) (int64, error) {
func PutSort2AttachmentViaVersion(tid int64, attachment *Attachment) (int64, error) {
	/*
		return Engine.Table(&Attachment{}).Where("id=?", tid).Update(&map[string]interface{}{
			"sort": sort,
		})
	*/
	return Engine.Id(tid).Cols("sort").Update(attachment)
}

func DelAttachment(id int64, uid int64, role int64) error {

	allow := false
	if role < 0 {
		allow = true
	}

	attachment := new(Attachment)

	if has, err := Engine.Id(id).Get(attachment); has == true && err == nil {

		if attachment.Uid == uid || allow {
			//检查附件字段并尝试删除文件
			if len(attachment.Content) > 0 {

				if p := helper.URL2local(attachment.Content); helper.Exist(p) {

					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误，但不要反回错误，以免陷入死循环无法删掉
							log.Println("ROOT DEL ATTACHMENT Content, ATTACHMENT ID:", id, ",ERR:", err)
						}
					} else { //检查用户对文件的所有权
						if helper.VerifyUserfile(p, strconv.Itoa(int(uid))) {
							if err := os.Remove(p); err != nil {
								log.Println("DEL ATTACHMENT Content, ATTACHMENT ID:", id, ",ERR:", err)
							}
						}
					}
				} /*else {

					////删除七牛云存储中的图片
					rsCli := rs.New(nil)

					//2014-8-21-134506AB4F24A56EF60543.jpg,2014-8-21-134506AB4F24A56EF60543.jpg
					imgkeys := strings.Split(attachment.Attachment, ",")
					for _, v := range imgkeys {
						rsCli.Delete(nil, BUCKET, strings.Split(v, ".")[0]) //key是32位  不含后缀
					}

				}
				*/
			}

			//检查内容字段并尝试删除文件
			if len(attachment.Content) > 0 {
				//若内容中存在图片则开始尝试删除图片
				delfiles_local := []string{}

				if m, n := helper.GetImages(attachment.Content); n > 0 {

					for _, v := range m {
						if helper.IsLocal(v) {
							delfiles_local = append(delfiles_local, v)
							//如果本地同时也存在banner缓存文件,则加入旧图集合中,等待后面一次性删除
							if p := helper.URL2local(helper.SetSuffix(v, "_banner.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.URL2local(helper.SetSuffix(v, "_large.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.URL2local(helper.SetSuffix(v, "_medium.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.URL2local(helper.SetSuffix(v, "_small.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
						}
					}
					for k, v := range delfiles_local {
						if p := helper.URL2local(v); helper.Exist(p) { //如若文件存在,则处理,否则忽略
							//先行判断是否缩略图  如果不是则执行删除image表记录的操作 因为缩略图是没有存到image表记录里面的
							/*
								isThumbnails := bool(true) //false代表不是缩略图 true代表是缩略图
								if (!strings.HasSuffix(v, "_large.jpg")) &&
									(!strings.HasSuffix(v, "_medium.jpg")) &&
									(!strings.HasSuffix(v, "_small.jpg")) {
									isThumbnails = false

								}
							*/

							//验证是否管理员权限
							if allow {
								if err := os.Remove(p); err != nil {
									log.Println("#", k, ",ROOT DEL FILE ERROR:", err)
								}

								//删除image表中已经被删除文件的记录
								/*
									if !isThumbnails {
										if e := DelImageByLocation(v); e != nil {
											fmt.Println("DelImageByLocation删除未使用文件", v, "的数据记录时候出现错误:", e)
										}
									}
								*/
							} else { //检查用户对文件的所有权
								if helper.VerifyUserfile(p, strconv.Itoa(int(uid))) {
									if err := os.Remove(p); err != nil {
										log.Println("#", k, ",DEL FILE ERROR:", err)
									}

									//删除image表中已经被删除文件的记录
									/*
										if !isThumbnails {
											if e := DelImageByLocation(v); e != nil {
												fmt.Println("v:", v)
												fmt.Println("DelImageByLocation删除未使用文件", v, "的数据记录时候出现错误:", e)
											}
										}
									*/
								}
							}

						}
					}
				}
			}

			//不管实际路径中是否存在文件均删除该数据库记录，以免数据库记录陷入死循环无法删掉
			if attachment.Id == id {

				if row, err := Engine.Id(id).Delete(new(Attachment)); err != nil || row == 0 {
					log.Println("row:", row, "删除附件错误:", err)
					return errors.New("E,删除附件错误!")
				}
				return nil
			}

		}
		return errors.New("你无权删除此附件:" + strconv.FormatInt(id, 10))
	}
	return errors.New("无法删除不存在的ATTACHMENT ID:" + strconv.FormatInt(id, 10))
}

func GetAttachmentCountByNid(nid int64) int64 {
	n, _ := Engine.Where("nid=?", nid).Count(&Attachment{Nid: nid})
	return n
}
