package models

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/insionng/yougam/helper"
)

type Topic struct {
	Id                int64
	Pid               int64   `xorm:"index"` //为0代表没有上级，本身就是顶层话题 大于0则本身是子话题，而该数字代表上级话题的id
	Cid               int64   `xorm:"index"`
	Nid               int64   `xorm:"index"` //nodeid
	Uid               int64   `xorm:"index"`
	Sort              int64   `xorm:"index"` //排序字段 需要手动对话题排序置顶等操作时使用
	Ctype             int64   `xorm:"index"` //ctype作用在于区分话题的类型
	Title             string  `xorm:"index"`
	Excerpt           string  `xorm:"index"` //摘录
	Content           string  `xorm:"text"`
	Tailinfo          string  `xorm:"index"` //尾巴信息 附带内容，譬如回复可见类型的话题时作为存储评论者UID集合
	Attachment        string  `xorm:"text"`  //附件 JSON
	Thumbnails        string  `xorm:"index"` //Original remote file
	ThumbnailsLarge   string  `xorm:"index"` //200x300
	ThumbnailsMedium  string  `xorm:"index"` //200x150
	ThumbnailsSmall   string  `xorm:"index"` //70x70
	Avatar            string  `xorm:"index"` //200x200
	AvatarLarge       string  `xorm:"index"` //100x100
	AvatarMedium      string  `xorm:"index"` //48x48
	AvatarSmall       string  `xorm:"index"` //32x32
	Tags              string  `xorm:"index"`
	Created           int64   `xorm:"created index"`
	Updated           int64   `xorm:"updated"`
	Hotness           float64 `xorm:"index"`
	Confidence        float64 `xorm:"index"` //信任度数值
	Hotup             int64   `xorm:"index"`
	Hotdown           int64   `xorm:"index"`
	Hotscore          int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64   `xorm:"index"` //Hotup  + Hotdown
	Views             int64
	Author            string `xorm:"index"`
	Template          string `xorm:"index"`
	Category          string `xorm:"index"`
	Node              string `xorm:"index"` //nodename
	ReplyTime         int64
	ReplyCount        int64
	ReplyLastUserId   int64
	ReplyLastUsername string
	ReplyLastNickname string
	FavoriteCount     int64
	Latitude          float64 `xorm:"index"`   //纬度
	Longitude         float64 `xorm:"index"`   //经度
	Version           int64   `xorm:"version"` //乐观锁
}

type TopicMark struct {
	Id  int64
	Uid int64 `xorm:"index"`
	Cid int64 `xorm:"index"` //该话题所属于的分类
	Tid int64 `xorm:"index"` //Topic id
}

type Topicjuser struct {
	Topic `xorm:"extends"`
	User  `xorm:"extends"`
}

type Topicjtopicmark struct {
	Topic     `xorm:"extends"`
	TopicMark `xorm:"extends"`
}

type Topicjtopicmarkjuser struct {
	Topic     `xorm:"extends"`
	TopicMark `xorm:"extends"`
	User      `xorm:"extends"`
}

/*
func SetTopicMark(uid int64, tid int64) (int64, error) {

	tpm := new(TopicMark)
	tpm.Uid = uid
	tpm.Tid = tid
	rows, err := Engine.Insert(tpm)
	return rows, err
}
*/

func SetTopicMark(uid, cid, tid int64) (int64, error) {
	qm := &TopicMark{Uid: uid, Cid: cid, Tid: tid}
	return Engine.Insert(qm)
}

/*
func DelTopicMark(uid int64, tid int64) {
	tpm := new(TopicMark)
	Engine.Where("uid=? and tid=?", uid, tid).Delete(tpm)
}
*/
func DelTopicMark(uid, tid int64) (int64, error) {
	return Engine.Where("uid=? and tid=?", uid, tid).Delete(new(TopicMark))
}

func TopicMarkCount(tid int64) (int64, error) {
	return Engine.Where("tid=?", tid).Count(&TopicMark{})
}

func TopicMarkCountByUid(uid int64) (int64, error) {
	return Engine.Where("uid=?", uid).Count(&TopicMark{})
}

func IsTopicMark(uid int64, tid int64) bool {

	tpm := new(TopicMark)

	if has, err := Engine.Where("uid=? and tid=?", uid, tid).Get(tpm); err != nil {
		return false
	} else {
		if has {
			return (tpm.Uid == uid)
		} else {
			return false
		}
	}

}

func GetTopicMarksViaUid(offset, limit int, uid int64, field string) (*[]*TopicMark, error) {
	objs := new([]*TopicMark)
	var err error
	if len(field) > 0 {
		err = Engine.Where("uid=?", uid).Limit(limit, offset).Desc(field).Find(objs)
	} else {
		err = Engine.Where("uid=?", uid).Limit(limit, offset).Find(objs)
	}
	return objs, err
}

func GetTopicMarksViaUidWithCid(offset, limit int, uid, cid int64, field string) (*[]*TopicMark, error) {
	objs := new([]*TopicMark)
	var err error
	if len(field) > 0 {
		err = Engine.Where("uid = ? and cid = ?", uid, cid).Limit(limit, offset).Desc(field).Find(objs)
	} else {
		err = Engine.Where("uid = ? and cid = ?", uid, cid).Limit(limit, offset).Find(objs)
	}
	return objs, err
}

func AddTopic(title, content, avatar, avatarLarge, avatarMedium, avatarSmall string, pid, cid, nid, uid int64) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	topicid := int64(0)
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

		tp := new(Topic)
		tp.Pid = pid
		tp.Cid = cat.Id
		tp.Nid = nid
		tp.Uid = uid
		tp.Title = title
		tp.Content = content
		tp.Author = usr.Username
		if pid == 0 { //若果本身是父级话题则处理
			tp.Avatar = avatar
			tp.AvatarLarge = avatarLarge
			tp.AvatarMedium = avatarMedium
			tp.AvatarSmall = avatarSmall
		}
		tp.Category = cat.Title
		tp.Node = nd.Title
		tp.Created = time.Now().Unix()

		if row, err := sess.Insert(tp); (err == nil) && (row > 0) {
			if tp.Uid > 0 {
				n, _ := sess.Where("uid=?", tp.Uid).Count(&Topic{})
				_u := map[string]interface{}{"topic_time": time.Now().Unix(), "topic_count": n, "topic_last_nid": tp.Nid, "topic_last_node": tp.Node}
				if row, err := sess.Table(&User{}).Where("id=?", tp.Uid).Update(&_u); err != nil || row <= 0 {

					return -1, errors.New(fmt.Sprint("AddTopic更新user表话题相关信息时:出现", row, "行影响,错误:", err))

				}
			}
			topicid = tp.Id
		} else {
			sess.Rollback()
			return -1, err
		}

	}

	// 提交事务
	return topicid, sess.Commit()
}

func SetTopic(tid int64, tp *Topic) (int64, error) {
	tp.Id = tid
	return Engine.Insert(tp)
}

func GetTopicsByHotnessNodes(nodelimit int, topiclimit int) []*[]*Topic {
	//找出最热的节点:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级话题
	nds, _ := GetNodes(0, nodelimit, "hotness")
	topics := make([]*[]*Topic, 0)

	if len(*nds) > 0 {
		i := 0
		for _, v := range *nds {
			i = i + 1
			tps := GetTopicsByNid(v.Id, 0, topiclimit, 0, "views")
			if len(*tps) != 0 {
				topics = append(topics, tps)
			}
			if i == len(*nds)-1 {
				break
			}
		}
	}

	return topics

}

func GetTopicsByScoreNodes(nodelimit int, topiclimit int) []*[]*Topic {
	//找出最热的节点:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级话题
	nds, _ := GetNodes(0, nodelimit, "hotscore")
	topics := make([]*[]*Topic, 0)

	if len(*nds) > 0 {
		i := 0
		for _, v := range *nds {
			i = i + 1
			tps := GetTopicsByNid(v.Id, 0, topiclimit, 0, "views")
			if len(*tps) != 0 {
				topics = append(topics, tps)
			}
			if i == len(*nds) {
				break
			}
		}
	}

	return topics

}

func GetTopicsByHotnessCategory(catlimit int, topiclimit int) []*[]*Topic {
	//找出最热的分类:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级话题
	cats, _ := GetCategories(0, catlimit, "hotness")
	topics := make([]*[]*Topic, 0)

	if len(*cats) > 0 {
		i := 0
		for _, v := range *cats {
			i = i + 1
			//(cid int64, offset int, limit int, ctype int64, field string)
			tps := GetTopicsByCid(v.Id, 0, topiclimit, 0, "views")
			if len(*tps) != 0 {
				topics = append(topics, tps)
			}
			if i == len(*cats) {
				break
			}
		}
	}

	return topics

}

func GetTopic(id int64) (*Topic, error) {
	tp := new(Topic)
	has, err := Engine.Id(id).Get(tp)
	if has {
		return tp, err
	} else {
		return nil, err
	}
}

func GetTopics(offset int, limit int, field string) (*[]*Topic, error) {
	tps := new([]*Topic)
	err := Engine.Limit(limit, offset).Desc(field).Find(tps)
	return tps, err
}

func GetSubTopics(pid int64, offset int, limit int, field string) (*[]*Topic, error) {
	tps := new([]*Topic)
	err := Engine.Where("pid=?", pid).Limit(limit, offset).Desc(field).Find(tps)
	return tps, err
}

func GetTopicsCount(offset int, limit int) (int64, error) {
	total, err := Engine.Limit(limit, offset).Count(&Topic{})
	return total, err
}

func GetTopicsByPid4Count(pid int64, offset int, limit int, ctype int64) (int64, error) {

	if ctype != 0 {
		return Engine.Where("(pid=? or id=?) and ctype=?", pid, pid, ctype).Limit(limit, offset).Count(&Topic{})
	} else {
		return Engine.Where("pid=? or id=?", pid, pid).Limit(limit, offset).Count(&Topic{})
	}

}

func GetTopicsCountByNode(node string, offset int, limit int) (int64, error) {
	total, err := Engine.Where("node=?", node).Limit(limit, offset).Count(&Topic{})
	return total, err
}

func GetTopicsByCategoryCount(category string, offset int, limit int, field string) (int64, error) {
	total, err := Engine.Where("category=?", category).Limit(limit, offset).Count(&Topic{})
	return total, err
}

//GetTopicsByCid大数据下如出现性能问题 可以使用 GetTopicsByCidOnBetween
func GetTopicsByCid(cid int64, offset int, limit int, ctype int64, field string) *[]*Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]*Topic)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Limit(limit, offset).Asc("id").Find(tps)
		} else {
			Engine.Where("cid=?", cid).Limit(limit, offset).Asc("id").Find(tps)

		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Desc(field).Limit(limit, offset).Find(tps)

		} else {
			if cid == 0 {
				Engine.Desc(field).Limit(limit, offset).Find(tps)
			} else {
				Engine.Where("cid=?", cid).Desc(field).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
		} else {
			if cid == 0 {
				Engine.Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("cid=?", cid).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			}
		}

	}
	return tps
}

func GetSubjectsByCategory4Count(category string, offset int, limit int, ctype int64) (int64, error) {

	if ctype != 0 {
		return Engine.Where("pid=0 and category=? and ctype=?", category, ctype).Limit(limit, offset).Count(&Topic{})
	} else {
		if category == "" {
			return Engine.Where("pid=0").Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0 and category=?", category).Limit(limit, offset).Count(&Topic{})
		}
	}
}

func GetSubjectsByCid4Count(cid int64, offset int, limit int, ctype int64) (int64, error) {

	if ctype != 0 {
		return Engine.Where("pid=0 and cid=? and ctype=?", cid, ctype).Limit(limit, offset).Count(&Topic{})
	} else {
		if cid == 0 {
			return Engine.Where("pid=0").Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0 and cid=?", cid).Limit(limit, offset).Count(&Topic{})
		}
	}
}

func GetSubjectsByCid(cid int64, offset int, limit int, ctype int64, field string) *[]*Topic {

	tps := new([]*Topic)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.cid=? and topic.ctype=?", cid, ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
		} else {
			if cid == 0 {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Asc("topic.id").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.cid=?", cid).Limit(limit, offset).Asc("topic.id").Find(tps)
			}
		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.cid=? and topic.ctype=?", cid, ctype).Desc("topic."+field).Limit(limit, offset).Find(tps)

		} else {
			if cid == 0 {
				Engine.Table("topic").Where("topic.pid=0").Desc("topic."+field).Limit(limit, offset).Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.cid=?", cid).Desc("topic."+field).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.cid=? and topic.ctype=?", cid, ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
		} else {
			if cid == 0 {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.cid=?", cid).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
		}

	}
	return tps
}

func GetSubjectsByCidJoinUser(cid int64, offset int, limit int, ctype int64, field string) *[]*Topicjuser {

	tps := new([]*Topicjuser)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.cid=? and topic.ctype=?", cid, ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		} else {
			if cid == 0 {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.cid=?", cid).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.cid=? and topic.ctype=?", cid, ctype).Desc("topic."+field).Limit(limit, offset).Join("LEFT", "user", "user.id = topic.uid").Find(tps)

		} else {
			if cid == 0 {
				Engine.Table("topic").Where("topic.pid=0").Desc("topic."+field).Limit(limit, offset).Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.cid=?", cid).Desc("topic."+field).Limit(limit, offset).Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.cid=? and topic.ctype=?", cid, ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		} else {
			if cid == 0 {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.cid=?", cid).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}

	}
	return tps
}

func GetSubjectsByCategory(category string, offset int, limit int, ctype int64, field string) *[]*Topic {

	tps := new([]*Topic)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.category=? and topic.ctype=?", category, ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
		} else {
			Engine.Table("topic").Where("topic.pid=0 and topic.category=?", category).Limit(limit, offset).Asc("topic.id").Find(tps)

		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.category=? and topic.ctype=?", category, ctype).Desc("topic."+field).Limit(limit, offset).Find(tps)

		} else {
			if category == "" {
				Engine.Table("topic").Where("topic.pid=0").Desc("topic."+field).Limit(limit, offset).Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.category=?", category).Desc("topic."+field).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.category=? and topic.ctype=?", category, ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
		} else {
			if category == "" {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.category=?", category).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
		}

	}
	return tps
}

func GetSubjectsByCategoryJoinUser(category string, offset int, limit int, ctype int64, field string) *[]*Topicjuser {

	tps := new([]*Topicjuser)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.category=? and topic.ctype=?", category, ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		} else {
			Engine.Table("topic").Where("topic.pid=0 and topic.category=?", category).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)

		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.category=? and topic.ctype=?", category, ctype).Desc("topic."+field).Limit(limit, offset).Join("LEFT", "user", "user.id = topic.uid").Find(tps)

		} else {
			if category == "" {
				Engine.Table("topic").Where("topic.pid=0").Desc("topic."+field).Limit(limit, offset).Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.category=?", category).Desc("topic."+field).Limit(limit, offset).Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Table("topic").Where("topic.pid=0 and topic.category=? and topic.ctype=?", category, ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		} else {
			if category == "" {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.category=?", category).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}

	}
	return tps
}

func GetSubjectsByUid(uid int64, offset int, limit int, ctype int64, field string) *[]*Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]*Topic)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Where("pid=0 and uid=? and ctype=?", uid, ctype).Limit(limit, offset).Asc("id").Find(tps)
		} else {
			Engine.Where("pid=0 and uid=?", uid).Limit(limit, offset).Asc("id").Find(tps)

		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Where("pid=0 and uid=? and ctype=?", uid, ctype).Desc(field).Limit(limit, offset).Find(tps)

		} else {
			if uid == 0 {
				Engine.Where("pid=0").Desc(field).Limit(limit, offset).Find(tps)
			} else {
				Engine.Where("pid=0 and uid=?", uid).Desc(field).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("pid=0 and uid=? and ctype=?", uid, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
		} else {
			if uid == 0 {
				Engine.Where("pid=0").Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("pid=0 and uid=?", uid).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			}
		}

	}
	return tps
}

func GetSubjectsByUsername(author string, offset int, limit int, ctype int64, field string) *[]*Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]*Topic)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Where("pid=0 and author=? and ctype=?", author, ctype).Limit(limit, offset).Asc("id").Find(tps)
		} else {
			Engine.Where("pid=0 and author=?", author).Limit(limit, offset).Asc("id").Find(tps)

		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Where("pid=0 and author=? and ctype=?", author, ctype).Desc(field).Limit(limit, offset).Find(tps)

		} else {
			if len(author) == 0 {
				Engine.Where("pid=0").Desc(field).Limit(limit, offset).Find(tps)
			} else {
				Engine.Where("pid=0 and author=?", author).Desc(field).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("pid=0 and author=? and ctype=?", author, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
		} else {
			if len(author) == 0 {
				Engine.Where("pid=0").Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("pid=0 and author=?", author).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			}
		}

	}
	return tps
}

func GetSubjectsCountByUsername(author string, offset int, limit int, ctype int64) (int64, error) {

	if ctype != 0 {
		return Engine.Where("pid=0 and author=? and ctype=?", author, ctype).Limit(limit, offset).Count(&Topic{})
	} else {
		if author == "" {
			return Engine.Where("pid=0").Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0 and author=?", author).Limit(limit, offset).Count(&Topic{})
		}
	}

}

func GetTopicsByCidOnBetween(cid int64, startid int64, endid int64, offset int, limit int, ctype int64, field string) *[]*Topic {
	tps := new([]*Topic)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=? and id>? and id<?", cid, ctype, startid-1, endid+1).Limit(limit, offset).Asc("id").Find(tps)
		} else {
			if cid == 0 {
				Engine.Where("id>? and id<?", startid-1, endid+1).Limit(limit, offset).Asc("id").Find(tps)
			} else {
				Engine.Where("cid=? and id>? and id<?", cid, startid-1, endid+1).Limit(limit, offset).Asc("id").Find(tps)
			}
		}
	default: //Desc
		if ctype != 0 {
			Engine.Where("cid=? and ctype=? and id>? and id<?", cid, ctype, startid-1, endid+1).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
		} else {
			if cid == 0 {
				Engine.Where("id>? and id<?", startid-1, endid+1).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("cid=? and id>? and id<?", cid, startid-1, endid+1).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			}
		}
	}

	return tps
}

func GetTopicsByCategory(category string, offset int, limit int, ctype int64, field string) *[]*Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]*Topic)
	switch {
	case field == "asc":
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Limit(limit, offset).Asc("id").Find(tps)
		} else {
			Engine.Where("category=?", category).Limit(limit, offset).Asc("id").Find(tps)

		}
	case field == "views" || field == "reply_count":
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Desc(field).Limit(limit, offset).Find(tps)

		} else {
			if category == "" {
				Engine.Desc(field).Limit(limit, offset).Find(tps)
			} else {
				Engine.Where("category=?", category).Desc(field).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
		} else {
			if category == "" {
				Engine.Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("category=?", category).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			}
		}

	}
	return tps
}

//GetTopicsByUid不区分父子话题
func GetTopicsByUid(uid int64, offset int, limit int, ctype int64, field string) *[]*Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]*Topic)

	switch {
	case field == "asc":
		if uid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(field).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("uid=? and ctype=?", uid, ctype).Limit(limit, offset).Asc("id").Find(tps)
			} else {
				Engine.Where("uid=?", uid).Limit(limit, offset).Asc("id").Find(tps)
			}
		}
	default:
		if uid <= 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(field).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("uid=? and ctype=?", uid, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("uid=?", uid).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			}
		}
	}
	return tps
}

func GetTopicsByNid(nodeid int64, offset int, limit int, ctype int64, field string) *[]*Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]*Topic)

	switch {
	case field == "asc":
		if nodeid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(field).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("nid=? and ctype=?", nodeid, ctype).Limit(limit, offset).Asc("id").Find(tps)
			} else {
				Engine.Where("nid=?", nodeid).Limit(limit, offset).Asc("id").Find(tps)
			}
		}
	default:
		if nodeid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(field).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("nid=? and ctype=?", nodeid, ctype).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("nid=?", nodeid).Limit(limit, offset).Desc(field, "views", "reply_count", "created").Find(tps)
			}
		}
	}
	return tps
}

func GetTopicsByPid(pid int64, offset int, limit int, ctype int64, field string) *[]*Topic {

	tps := new([]*Topic)
	if len(field) == 0 {
		field = "id"
	}
	if limit == 0 {
		switch {
		case field == "asc":
			{
				if pid > 0 { //即只查询单个话题的父级与子级，此时sort字段无效，所以无须作为条件参与排序
					if ctype != 0 {
						Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=?", pid, pid, ctype).Asc("topic.id").Find(tps)
					} else {
						Engine.Table("topic").Where("(topic.pid=? or topic.id=?)", pid, pid).Asc("topic.id").Find(tps)
					}
				} else {
					if ctype != 0 {
						Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Asc("topic.id").Find(tps)
					} else {
						Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort<=?", pid, pid, 0).Asc("topic.id").Find(tps)
					}
				}

			}
		case field == "cold":
			{
				if ctype != 0 {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Asc("topic.views").Find(tps)
				} else {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort<=?", pid, pid, 0).Asc("topic.views").Find(tps)
				}
			}
		case field == "sort":
			{ //Desc("topic.sort", "topic.confidence") ，这段即是"topic.sort"的优先级较"topic.confidence"要高
				if ctype != 0 {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort>?", pid, pid, ctype, 0).Desc("topic.sort", "topic.confidence").Find(tps)
				} else {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort>?", pid, pid, 0).Desc("topic.sort", "topic.confidence").Find(tps)
				}
			}
		default:
			{
				if field == "desc" {
					field = "id"
				}
				if ctype != 0 {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
				} else {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?)", pid, pid).And("topic.sort<=?", 0).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
				}
			}
		}
	} else {
		switch {
		case field == "asc":
			{
				if pid > 0 { //即只查询单个话题的父级与子级，此时sort字段无效，所以无须作为条件参与排序
					if ctype != 0 {
						Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=?", pid, pid, ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
					} else {
						Engine.Table("topic").Where("(topic.pid=? or topic.id=?)", pid, pid).Limit(limit, offset).Asc("topic.id").Find(tps)
					}
				} else {
					if ctype != 0 {
						Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Asc("topic.id").Find(tps)
					} else {
						Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort<=?", pid, pid, 0).Limit(limit, offset).Asc("topic.id").Find(tps)
					}
				}

			}
		case field == "cold":
			{
				if ctype != 0 {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Asc("topic.views").Find(tps)
				} else {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort<=?", pid, pid, 0).Limit(limit, offset).Asc("topic.views").Find(tps)
				}
			}
		case field == "sort":
			{ //Desc("topic.sort", "topic.confidence") ，这段即是"topic.sort"的优先级较"topic.confidence"要高
				if ctype != 0 {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort>?", pid, pid, ctype, 0).Limit(limit, offset).Desc("topic.sort", "topic.confidence").Find(tps)
				} else {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort>?", pid, pid, 0).Limit(limit, offset).Desc("topic.sort", "topic.confidence").Find(tps)
				}
			}
		default:
			{
				if field == "desc" {
					field = "id"
				}
				if ctype != 0 {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
				} else {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?)", pid, pid).And("topic.sort<=?", 0).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
				}
			}
		}
	}
	return tps
}

func GetTopicsViaPid(pid int64, offset int, limit int, ctype int64, field string) *[]*Topic {
	var objs = new([]*Topic)
	if pid <= 0 {
		switch {
		case field == "asc":
			{
				if ctype != 0 {
					Engine.Where("pid <= ? and ctype = ?", 0, ctype).Limit(limit, offset).Asc("id").Find(objs)
				} else {
					Engine.Where("pid <= ?", 0).Limit(limit, offset).Asc("id").Find(objs)
				}
			}
		default:
			{
				if ctype != 0 {
					Engine.Where("pid <= ? and ctype = ?", 0, ctype).Limit(limit, offset).Desc(field).Find(objs)
				} else {
					Engine.Where("pid <= ?", 0).Limit(limit, offset).Desc(field).Find(objs)
				}
			}
		}
	} else {
		switch {
		case field == "asc":
			{
				if ctype != 0 {
					Engine.Where("pid=? and ctype=?", pid, ctype).Limit(limit, offset).Asc("id").Find(objs)
				} else {
					Engine.Where("pid=?", pid).Limit(limit, offset).Asc("id").Find(objs)
				}
			}
		default:
			{
				if ctype != 0 {
					Engine.Where("pid=? and ctype=?", pid, ctype).Limit(limit, offset).Desc(field).Find(objs)
				} else {
					Engine.Where("pid=?", pid).Limit(limit, offset).Desc(field).Find(objs)
				}
			}
		}
	}
	return objs
}

func GetTopicsByPidJoinUser(pid int64, offset int, limit int, ctype int64, field string) *[]*Topicjuser {

	tps := new([]*Topicjuser)

	switch {
	case field == "asc":
		{
			if pid > 0 { //即只查询单个话题的父级与子级，此时sort字段无效，所以无须作为条件参与排序
				if ctype != 0 {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=?", pid, pid, ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
				} else {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?)", pid, pid).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
				}
			} else {
				if ctype != 0 {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
				} else {
					Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort<=?", pid, pid, 0).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
				}
			}

		}
	case field == "cold":
		{
			if ctype != 0 {
				Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Asc("topic.views").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort<=?", pid, pid, 0).Limit(limit, offset).Asc("topic.views").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	case field == "sort":
		{ //Desc("topic.sort", "topic.confidence") ，这段即是"topic.sort"的优先级较"topic.confidence"要高
			if ctype != 0 {
				Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort>?", pid, pid, ctype, 0).Limit(limit, offset).Desc("topic.sort", "topic.confidence").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.sort>?", pid, pid, 0).Limit(limit, offset).Desc("topic.sort", "topic.confidence").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	default:
		{
			if field == "desc" {
				field = "id"
			}
			if ctype != 0 {
				Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", pid, pid, ctype, 0).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("(topic.pid=? or topic.id=?)", pid, pid).And("topic.sort<=?", 0).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	}
	return tps
}

func GetTopicsByPidSinceCreated(pid int64, offset int, limit int, ctype int64, field string, since int64) *[]*Topic {

	switch tps := new([]*Topic); {
	case field == "asc":
		{
			if ctype != 0 {
				Engine.Table("topic").Where("topic.created>? and (topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", since, pid, pid, ctype, 0).Limit(limit, offset).Asc("topic.id").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.created>? and (topic.pid=? or topic.id=?) and topic.sort<=?", since, pid, pid, 0).Limit(limit, offset).Asc("topic.id").Find(tps)
			}
			return tps
		}
	default:
		{
			if ctype != 0 {
				Engine.Table("topic").Where("topic.created>? and (topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", since, pid, pid, ctype, 0).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.created>? and (topic.pid=? or topic.id=?) and topic.sort<=?", since, pid, pid, 0).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
			return tps
		}
	}

}

func GetTopicsByPidJoinUserSinceCreated(pid int64, offset int, limit int, ctype int64, field string, since int64) *[]*Topicjuser {

	switch tps := new([]*Topicjuser); {
	case field == "asc":
		{
			if ctype != 0 {
				Engine.Table("topic").Where("topic.created>? and (topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", since, pid, pid, ctype, 0).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.created>? and (topic.pid=? or topic.id=?) and topic.sort<=?", since, pid, pid, 0).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
			return tps
		}
	default:
		{
			if ctype != 0 {
				Engine.Table("topic").Where("topic.created>? and (topic.pid=? or topic.id=?) and topic.ctype=? and topic.sort<=?", since, pid, pid, ctype, 0).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.created>? and (topic.pid=? or topic.id=?) and topic.sort<=?", since, pid, pid, 0).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
			return tps
		}
	}

}

func GetTopicsByPidJoinTopicmark(pid int64, offset int, limit int, ctype int64, field string) *[]*Topicjtopicmark {

	tps := new([]*Topicjtopicmark)

	switch {
	case field == "asc":
		{
			if ctype != 0 {
				Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=?", pid, pid, ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "topic_mark", "topic_mark.uid = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=? or topic.id=?", pid, pid).Limit(limit, offset).Asc("topic.id").Join("LEFT", "topic_mark", "topic_mark.uid = topic.uid").Find(tps)
			}
		}
	default:
		{
			if ctype != 0 {
				Engine.Table("topic").Where("(topic.pid=? or topic.id=?) and topic.ctype=?", pid, pid, ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "topic_mark", "topic_mark.uid = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=? or topic.id=?", pid, pid).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "topic_mark", "topic_mark.uid = topic.uid").Find(tps)
			}
		}
	}
	return tps
}

func JoinTopicmarkJoinUserForGetTopicsByPid(pid int64, offset int, limit int, field string) *[]*Topicjtopicmarkjuser {
	tps := new([]*Topicjtopicmarkjuser)
	switch {
	case field == "asc":
		{
			Engine.Table("topic_mark").Where("topic_mark.tid=?", pid).Limit(limit, offset).Asc("topic_mark.id").Join("LEFT", "topic", "topic_mark.tid = topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		}
	default:
		{
			Engine.Table("topic_mark").Where("topic_mark.tid=?", pid).Limit(limit, offset).Desc("topic_mark."+field).Join("LEFT", "topic", "topic_mark.tid = topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		}
	}
	return tps
}

func JoinTopicmarkJoinUserForGetTopicsByUid(uid int64, offset int, limit int, field string) *[]*Topicjtopicmarkjuser {
	tps := new([]*Topicjtopicmarkjuser)
	switch {
	case field == "asc":
		{
			Engine.Table("topic_mark").Where("topic_mark.uid=?", uid).Limit(limit, offset).Asc("topic_mark.id").Join("LEFT", "topic", "topic_mark.tid = topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		}
	default:
		{
			Engine.Table("topic_mark").Where("topic_mark.uid=?", uid).Limit(limit, offset).Desc("topic_mark."+field).Join("LEFT", "topic", "topic_mark.tid = topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
		}
	}
	return tps
}

func GetSubjectsByNid4Count(nodeid int64, offset int, limit int, ctype int64) (int64, error) {

	if nodeid == 0 {
		if ctype != 0 {
			return Engine.Where("pid=0 and ctype=?", ctype).Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0").Limit(limit, offset).Count(&Topic{})
		}
	} else {
		if ctype != 0 {
			return Engine.Where("pid=0 and nid=? and ctype=?", nodeid, ctype).Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0 and nid=?", nodeid).Limit(limit, offset).Count(&Topic{})
		}
	}
}

func GetSubjectsByNid(nodeid int64, offset int, limit int, ctype int64, field string) *[]*Topic {

	tps := new([]*Topic)

	switch {
	case field == "asc":
		if nodeid == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Asc("topic.id").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.nid=? and topic.ctype=?", nodeid, ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.nid=?", nodeid).Limit(limit, offset).Asc("topic.id").Find(tps)
			}
		}
	default:
		if nodeid == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.nid=? and topic.ctype=?", nodeid, ctype).Limit(limit, offset).Desc("topic."+field, "topic.topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.nid=?", nodeid).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
		}
	}
	return tps
}

func GetSubjectsByNidJoinUser(nodeid int64, offset int, limit int, ctype int64, field string) *[]*Topicjuser {

	tps := new([]*Topicjuser)

	switch {
	case field == "asc":
		if nodeid == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.nid=? and topic.ctype=?", nodeid, ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.nid=?", nodeid).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	default:
		if nodeid == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.nid=? and topic.ctype=?", nodeid, ctype).Limit(limit, offset).Desc("topic."+field, "topic.topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.nid=?", nodeid).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	}
	return tps
}

func GetTopicsByNode(node string, offset int, limit int, field string) (*[]*Topic, error) {
	tps := new([]*Topic)
	err := Engine.Where("node=?", node).Limit(limit, offset).Desc(field).Find(tps)
	return tps, err
}

func GetSubjectsByNode4Count(node string, offset int, limit int, ctype int64) (int64, error) {

	if node == "" {
		if ctype != 0 {
			return Engine.Where("pid=0 and ctype=?", ctype).Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0").Limit(limit, offset).Count(&Topic{})
		}
	} else {
		if ctype != 0 {
			return Engine.Where("pid=0 and node=? and ctype=?", node, ctype).Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0 and node=?", node).Limit(limit, offset).Count(&Topic{})
		}
	}
}

func GetSubjectsByNode(node string, offset int, limit int, ctype int64, field string) *[]*Topic {

	tps := new([]*Topic)

	switch {
	case field == "asc":
		if len(node) == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Asc("topic.id").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.node=? and topic.ctype=?", node, ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.node=?", node).Limit(limit, offset).Asc("topic.id").Find(tps)
			}
		}
	default:
		if len(node) == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.node=? and topic.ctype=?", node, ctype).Limit(limit, offset).Desc("topic."+field, "topic.topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.node=?", node).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
		}
	}
	return tps
}

func GetSubjectsByNodeJoinUser(node string, offset int, limit int, ctype int64, field string) *[]*Topicjuser {

	tps := new([]*Topicjuser)

	switch {
	case field == "asc":
		if len(node) == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.node=? and topic.ctype=?", node, ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.node=?", node).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	default:
		if len(node) == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.node=? and topic.ctype=?", node, ctype).Limit(limit, offset).Desc("topic."+field, "topic.topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.node=?", node).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	}
	return tps
}

func GetSubjectsByUser4Count(user string, offset int, limit int, ctype int64) (int64, error) {

	if user == "" {
		if ctype != 0 {
			return Engine.Where("pid=0 and ctype=?", ctype).Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0").Limit(limit, offset).Count(&Topic{})
		}
	} else {
		if ctype != 0 {
			return Engine.Where("pid=0 and author=? and ctype=?", user, ctype).Limit(limit, offset).Count(&Topic{})
		} else {
			return Engine.Where("pid=0 and author=?", user).Limit(limit, offset).Count(&Topic{})
		}
	}
}

func GetSubjectsByUser(username string, offset int, limit int, ctype int64, field string) *[]*Topic {

	tps := new([]*Topic)

	switch {
	case field == "asc":
		if len(username) == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Asc("topic.id").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.author=? and topic.ctype=?", username, ctype).Limit(limit, offset).Asc("topic.id").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.author=?", username).Limit(limit, offset).Asc("topic.id").Find(tps)
			}
		}
	default:
		if len(username) == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.author=? and topic.ctype=?", username, ctype).Limit(limit, offset).Desc("topic."+field, "topic.topic.views", "topic.reply_count", "topic.created").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.author=?", username).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Find(tps)
			}
		}
	}
	return tps
}

func GetSubjectsByUserJoinUser(username string, offset int, limit int, ctype int64, field string) *[]*Topicjuser {

	tps := new([]*Topicjuser)

	switch {
	case field == "asc":
		if len(username) == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.author=? and topic.ctype=?", username, ctype).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.author=?", username).Limit(limit, offset).Asc("topic.id").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	default:
		if len(username) == 0 {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.ctype=?", ctype).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0").Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		} else {
			if ctype != 0 {
				Engine.Table("topic").Where("topic.pid=0 and topic.author=? and topic.ctype=?", username, ctype).Limit(limit, offset).Desc("topic."+field, "topic.topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			} else {
				Engine.Table("topic").Where("topic.pid=0 and topic.author=?", username).Limit(limit, offset).Desc("topic."+field, "topic.views", "topic.reply_count", "topic.created").Join("LEFT", "user", "user.id = topic.uid").Find(tps)
			}
		}
	}
	return tps
}

//发布话题  返回 话题id,错误
func PostTopic(tp *Topic) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	topicid := int64(0)
	if _, err := sess.Insert(tp); err == nil && tp != nil {

		if tp.Nid > 0 {
			n, _ := sess.Where("nid=?", tp.Nid).Count(&Topic{})
			_n := map[string]interface{}{"author": tp.Author, "topic_time": time.Now().Unix(), "topic_count": n, "topic_last_user_id": tp.Uid}
			if row, err := sess.Table(&Node{}).Where("id=?", tp.Nid).Update(&_n); err != nil || row <= 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PostTopic更新node表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		if tp.Uid > 0 {
			n, _ := sess.Where("uid=?", tp.Uid).Count(&Topic{})
			_u := map[string]interface{}{"topic_time": time.Now().Unix(), "topic_count": n, "topic_last_nid": tp.Nid, "topic_last_node": tp.Node}
			if row, err := sess.Table(&User{}).Where("id=?", tp.Uid).Update(&_u); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PostTopic更新user表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		topicid = tp.Id
	} else {
		// 发生错误时进行回滚
		sess.Rollback()
		return -1, err
	}

	// 提交事务
	return topicid, sess.Commit()
}

func PutTopic(tid int64, tp *Topic) (int64, error) {
	// 创建 Session 对象
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	topicid := int64(0)
	if _, err := sess.Id(tid).Update(tp); (err == nil) && (tp != nil) {

		if tp.Nid > 0 {
			n, _ := sess.Where("nid=?", tp.Nid).Count(&Topic{})
			_n := map[string]interface{}{"author": tp.Author, "topic_time": time.Now().Unix(), "topic_count": n, "topic_last_user_id": tp.Uid}
			if row, err := sess.Table(&Node{}).Where("id=?", tp.Nid).Update(&_n); err != nil || row <= 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PutTopic更新node表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		if tp.Uid > 0 {
			n, _ := sess.Where("uid=?", tp.Uid).Count(&Topic{})
			_u := map[string]interface{}{"topic_time": time.Now().Unix(), "topic_count": n, "topic_last_nid": tp.Nid, "topic_last_node": tp.Node}
			if row, err := sess.Table(&User{}).Where("id=?", tp.Uid).Update(&_u); err != nil || row == 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PutTopic更新user表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		topicid = tp.Id
	} else {
		// 发生错误时进行回滚
		sess.Rollback()
		return -1, err
	}

	// 提交事务
	return topicid, sess.Commit()
}

/*
func PutTopic(tid int64, tp *Topic) (int64, error) {
	//覆盖式更新
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	topicid := int64(0)
	//	tp.Id = tid
	if row, err := sess.Table(new(Topic)).Id(tid).Update(tp); (err != nil) || (row <= 0) {
		sess.Rollback()
		return -1, err
	} else {
		if tp.Uid > 0 {
			n, _ := sess.Where("uid=?", tp.Uid).Count(&Topic{})
			_u := map[string]interface{}{"topic_time": time.Now().Unix(), "topic_count": n, "topic_last_nid": tp.Nid, "topic_last_node": tp.Node}
			if row, err := sess.Table(&User{}).Where("id=?", tp.Uid).Update(&_u); err != nil || row <= 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PutTopic更新user表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		topicid = tp.Id
	}

	// 提交事务
	return topicid, sess.Commit()
}
*/

func PutTopicViaVersion(tid int64, topic *Topic) (int64, error) {
	/*
		//乐观锁目前不支持map 以下句式无效！
			return Engine.Table(&Topic{}).Where("id=?", tid).Update(&map[string]interface{}{
				"views":   views,
				"version": version,
			})
	*/
	//return Engine.Table(&Topic{}).Where("id=?", tid).Update(topic)
	return Engine.Id(tid).Update(topic)
}

func PutViews2TopicViaVersion(tid int64, topic *Topic) (int64, error) {
	/*
		return Engine.Table(&Topic{}).Where("id=?", tid).Update(&map[string]interface{}{
			"views": views,
		})
	*/
	return Engine.Id(tid).Cols("views").Update(topic)
}

//func PutSort2Topic(tid int64, sort int64) (int64, error) {
func PutSort2TopicViaVersion(tid int64, topic *Topic) (int64, error) {
	/*
		return Engine.Table(&Topic{}).Where("id=?", tid).Update(&map[string]interface{}{
			"sort": sort,
		})
	*/
	return Engine.Id(tid).Cols("sort").Update(topic)
}

/*
func PutTailinfo2Topic(tid int64, Tailinfo string) (int64, error) {
	return Engine.Table(&Topic{}).Where("id=?", tid).Update(&map[string]interface{}{
		"tailinfo": func() string {
			tp, e := GetTopic(tid)
			if (e != nil) || (tp == nil) {
				return ""
			}
			if len(tp.Tailinfo) == 0 {
				tp.Tailinfo = Tailinfo
			} else {
				tp.Tailinfo = tp.Tailinfo + "," + Tailinfo
			}

			return tp.Tailinfo

		},
	})
}
*/

func PutTailinfo2Topic(tid int64, Tailinfo string) (int64, error) {

	tp, e := GetTopic(tid)
	if (e != nil) || (tp == nil) {
		return -1, e
	}

	if len(tp.Tailinfo) == 0 {
		tp.Tailinfo = Tailinfo
	} else {
		tp.Tailinfo = tp.Tailinfo + "," + Tailinfo
	}

	return PutTopic(tid, tp)

}

func IsTailinfoOfUser(tid, uid int64) bool {
	tp, e := GetTopic(tid)
	if (e != nil) || (tp == nil) {
		return false
	}
	if len(tp.Tailinfo) == 0 {
		return false
	} else {
		for _, v := range strings.Split(tp.Tailinfo, ",") {
			if fmt.Sprintf("%v", uid) == v {
				return true
			}
		}
	}
	return false
}

func DelTopicsByPid(id int64, uid int64, role int64) (err error) {
	objs := GetTopicsByPid(id, 0, 0, 0, "id")
	for _, v := range *objs {
		if e := DelTopic(v.Id, uid, role); e != nil {
			err = fmt.Errorf("%v;%v", e, err)
		}
	}
	return
}

func DelTopic(id int64, uid int64, role int64) error {

	allow := false
	if role < 0 {
		allow = true
	}

	topic := new(Topic)

	if has, err := Engine.Id(id).Get(topic); has == true && err == nil {

		if topic.Uid == uid || allow {
			//检查附件字段并尝试删除文件
			if len(topic.Attachment) > 0 {
				attachments := GetAttachmentsByPid(topic.Id, 0, 0, 0, "id")
				for _, v := range *attachments {
					if p := helper.URL2local(v.Content); helper.Exist(p) {
						//验证是否管理员权限
						if allow {
							if err := os.Remove(p); err != nil {
								//可以输出错误，但不要返回错误，以免陷入死循环无法删掉
								log.Println("ROOT DEL TOPIC Attachment, TOPIC ID:", id, ",ERR:", err)
							} else {
								topic.Attachment = "" // topic.Attachment - 1
							}
						} else { //检查用户对文件的所有权
							if helper.VerifyUserfile(p, strconv.Itoa(int(uid))) {
								if err := os.Remove(p); err != nil {
									log.Println("DEL TOPIC Attachment, TOPIC ID:", id, ",ERR:", err)
								} else {
									topic.Attachment = "" //topic.Attachment - 1
								}
							}
						}
					}
				}
			}

			//检查内容字段并尝试删除文件
			if len(topic.Content) > 0 {
				//若内容中存在图片则开始尝试删除图片
				delfiles_local := []string{}

				if m, n := helper.GetImages(topic.Content); n > 0 {

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
			if topic.Id == id {
				if row, err := Engine.Id(id).Delete(new(Topic)); err != nil || row == 0 {
					log.Println("row:", row, "删除话题错误:", err)
					return errors.New("E,删除话题错误!")
				}
				return nil
			}

		}
		return errors.New("你无权删除此话题:" + strconv.FormatInt(id, 10))
	}
	return errors.New("无法删除不存在的TOPIC ID:" + strconv.FormatInt(id, 10))
}

func GetTopicCountByNid(nid int64) int64 {
	n, _ := Engine.Where("nid=?", nid).Count(&Topic{Nid: nid})
	return n
}
