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

type Node struct {
	Id              int64
	Pid             int64 `xorm:"index"` //pid为0 代表节点本身是顶层节点, pid大于0 代表上级节点id , 编写逻辑的时候注意判断此pid不能等于自身id
	Cid             int64 `xorm:"index"` //所属分类的id
	Uid             int64 `xorm:"index"`
	Sort            int64
	Ctype           int64   `xorm:"index"`
	Title           string  `xorm:"index"`
	Content         string  `xorm:"text"`
	Attachment      string  `xorm:"text"`
	Created         int64   `xorm:"created index"`
	Updated         int64   `xorm:"updated"`
	Hotness         float64 `xorm:"index"`
	Confidence      float64 `xorm:"index"` //信任度数值
	Hotup           int64   `xorm:"index"`
	Hotdown         int64   `xorm:"index"`
	Hotscore        int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote         int64   `xorm:"index"` //Hotup  + 	Hotdown
	Views           int64
	Author          string `xorm:"index"` //节点的创建者
	Parent          string `xorm:"index"` //父级节点名称
	Category        string `xorm:"index"` //所属分类标题
	Tid             int64  //最后一次发布的topic的Tid
	Topic           int64  //最后一次发布的topic的标题
	TopicTime       int64
	TopicCount      int64
	TopicLastUserId int64
	FavoriteCount   int64  `xorm:"index"`
	Template        string `xorm:"index"`
}

type NodeMark struct {
	Id  int64
	Uid int64 `xorm:"index"`
	Nid int64 `xorm:"index"` //node id
}

type NodeBookmark struct {
	Id     int64
	Userid int64 `xorm:"index"`
	Nodeid int64 `xorm:"index"`
}

func HasNode(title string) (int64, bool) {

	nd := new(Node)
	if has, err := Engine.Where("title = ?", title).Get(nd); err != nil {
		return -1, false
	} else {
		if has {
			return nd.Id, true
		}
		return -1, false
	}
}

//返回所有存在的节点
func AllExistingNodes(offset int, limit int, field string) (*[]*Node, error) {
	nds := new([]*Node)
	err := Engine.Where("cid>=?", -2).Limit(limit, offset).Desc(field).Find(nds)
	return nds, err
}

func AvailableNodes(offset int, limit int, field string) (*[]*Node, error) {
	nds := new([]*Node)
	err := Engine.Where("cid>=?", -1).Limit(limit, offset).Desc(field).Find(nds)
	return nds, err
}

func NodesOfNodes(offset int, limit int, field string) (*[]*Node, error) {
	nds := new([]*Node)
	err := Engine.Where("cid>?", 0).Limit(limit, offset).Desc(field).Find(nds)
	return nds, err
}

func NodesOfNavor(offset int, limit int, field string) (*[]*Node, error) {
	nds := new([]*Node)
	err := Engine.Where("cid<=? and cid>=?", 0, -1).Limit(limit, offset).Desc(field).Find(nds)
	return nds, err
}

func SetNodeMark(uid int64, nid int64) (int64, error) {
	nm := new(NodeMark)
	nm = &NodeMark{Uid: uid, Nid: nid}
	rows, err := Engine.Insert(nm)
	return rows, err
}

func IsNodeMark(uid int64, nid int64) bool {

	nm := new(NodeMark)

	if has, err := Engine.Where("uid=? and nid=?", uid, nid).Get(nm); err != nil {
		return false
	} else {
		if has {
			if nm.Uid == uid {
				return true
			} else {
				return false
			}

		} else {
			return false
		}
	}

}

func GetNodesViaPid(pid int64, offset int, limit int, ctype int64, field string) *[]*Node {
	var objs = new([]*Node)
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

func GetNodesByCtypeWithNid(offset int, limit int, ctype int64, nid int64, field string) (*[]*Node, error) {
	objs := new([]*Node)
	err := Engine.Table("node").Where("ctype = ? and nid = ?", ctype, nid).Limit(limit, offset).Desc(field).Find(objs)
	return objs, err
}

func GetNodesByCtype(offset int, limit int, ctype int64, field string) (*[]*Node, error) {
	catz := new([]*Node)
	err := Engine.Table("node").Where("ctype = ?", ctype).Limit(limit, offset).Desc(field).Find(catz)
	return catz, err
}

/*
func GetNodesByCid(cid int64, offset int, limit int, field string) (*[]*Node, error) {
	nds := new([]*Node)
	err := Engine.Where("cid=?", cid).Limit(limit, offset).Desc(field).Find(nds)
	return nds, err
}
*/
func GetNodesByCid(cid int64, offset int, limit int, field string) (*[]*Node, error) {
	nds := new([]*Node)
	var err error
	if cid != 0 {
		/*
			if cid == -2 { //cid为-2时节点对前端隐藏展示
				err = Engine.Where("cid>?", cid).Limit(limit, offset).Desc(field).Find(nds)
			} else {
				err = Engine.Where("cid=?", cid).Limit(limit, offset).Desc(field).Find(nds)
			}
		*/
		err = Engine.Where("cid=?", cid).Limit(limit, offset).Desc(field).Find(nds)

	} else {
		err = Engine.Limit(limit, offset).Desc(field).Find(nds)
	}
	return nds, err
}

func GetNodesByCtypeWithPid(offset int, limit int, ctype int64, pid int64, field string) (*[]*Node, error) {
	catz := new([]*Node)
	err := Engine.Table("node").Where("ctype = ? and pid = ?", ctype, pid).Limit(limit, offset).Desc(field).Find(catz)
	return catz, err
}

func AddNode(title string, content string, attachment string, nid int64, cid int64, uid int64) (int64, error) {

	nd := &Node{Pid: nid, Cid: cid, Uid: uid, Title: title, Content: content, Attachment: attachment, Created: time.Now().Unix()}
	if _, err := Engine.Insert(nd); err == nil {
		return nd.Id, err
	} else {
		return -1, err
	}

}

func SetNode(nid int64, nd *Node) (int64, error) {
	nd.Id = nid
	return Engine.Insert(nd)
}

func GetNode(id int64) (*Node, error) {

	nd := &Node{}
	has, err := Engine.Id(id).Get(nd)
	if has {
		return nd, err
	} else {

		return nil, err
	}

}

func GetNodeByTitle(title string) (*Node, error) {
	nd := &Node{}
	nd.Title = title
	has, err := Engine.Get(nd)
	if has {
		return nd, err
	} else {

		return nil, err
	}
}

func GetNodesByCtypeWithCid(offset int, limit int, ctype int64, cid int64, field string) (*[]*Node, error) {
	catz := new([]*Node)
	err := Engine.Table("node").Where("ctype = ? and cid = ?", ctype, cid).Limit(limit, offset).Desc(field).Find(catz)
	return catz, err
}

func GetNodesByPid(pid int64, offset int, limit int, ctype int64, field string) *[]*Node {

	nds := new([]*Node)

	switch {
	case field == "asc":
		{
			if ctype != 0 {
				Engine.Where("(pid=? or id=?) and ctype=?", pid, pid, ctype).Limit(limit, offset).Asc("id").Find(nds)
			} else {
				Engine.Where("pid=? or id=?", pid, pid).Limit(limit, offset).Asc("id").Find(nds)
			}
		}
	default:
		{
			if ctype != 0 {
				Engine.Where("(pid=? or id=?) and ctype=?", pid, pid, ctype).Limit(limit, offset).Desc(field).Find(nds)
			} else {
				Engine.Where("pid=? or id=?", pid, pid).Limit(limit, offset).Desc(field).Find(nds)
			}
		}
	}
	return nds
}

func GetNodes(offset int, limit int, field string) (*[]*Node, error) {
	nds := new([]*Node)
	err := Engine.Limit(limit, offset).Desc(field).Find(nds)
	return nds, err
}

func PutNode(nid int64, nd *Node) (int64, error) {
	//覆盖式更新
	sess := Engine.NewSession()
	defer sess.Close()
	// 启动事务
	if err := sess.Begin(); err != nil {
		return -1, err
	}

	//执行事务
	nodeid := int64(0)
	if row, err := sess.Update(nd, &Node{Id: nid}); err != nil || row <= 0 {
		sess.Rollback()
		return -1, err
	} else {
		if nd.Uid > 0 {
			n, _ := sess.Where("uid=?", nd.Uid).Count(&Node{})
			_u := map[string]interface{}{"node_time": time.Now().Unix(), "node_count": n, "node_last_tid": nd.Tid, "node_last_topic": nd.Topic}
			if row, err := sess.Table(&User{}).Where("id=?", nd.Uid).Update(&_u); err != nil || row <= 0 {
				sess.Rollback()
				return -1, errors.New(fmt.Sprint("PutNode更新user表话题相关信息时,执行:", row, "行变更,出现错误:", err))
			}
		}

		nodeid = nd.Id
	}

	// 提交事务
	return nodeid, sess.Commit()
}

//map[string]interface{}{"ctype": ctype}
func UpdateNode(nid int64, nodemap *map[string]interface{}) error {
	nd := &Node{}
	if row, err := Engine.Table(nd).Where("id=?", nid).Update(nodemap); err != nil || row == 0 {
		log.Println("UpdateNode  row:::", row, "UpdateNode出现错误:", err)
		return err
	} else {
		return nil
	}

}

func DelNode(id int64, uid int64, role int64) error {
	allow := false
	if role < 0 {
		allow = true
	}

	node := &Node{}

	if has, err := Engine.Id(id).Get(node); has == true && err == nil {

		if node.Uid == uid || allow {
			//检查附件字段并尝试删除文件
			if len(node.Attachment) > 0 {

				if p := helper.URL2local(node.Attachment); helper.Exist(p) {
					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误日志，但不要反回错误，以免陷入死循环无法删掉
							fmt.Println("ROOT DEL NODE Attachment, NODE ID:", id, ",ERR:", err)
						}
					} else { //检查用户对本地文件的所有权
						if helper.VerifyUserfile(p, strconv.FormatInt(uid, 10)) {
							if err := os.Remove(p); err != nil {
								fmt.Println("DEL NODE Attachment, NODE ID:", id, ",ERR:", err)
							}
						}
					}

				}
			}

			//检查内容字段并尝试删除文件
			if len(node.Content) > 0 {
				//若内容中存在图片则开始尝试删除图片
				delfiles_local := []string{}

				if m, n := helper.GetImages(node.Content); n > 0 {

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
							//isThumbnails := bool(true) //false代表不是缩略图 true代表是缩略图
							/*
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

							} else { //检查用户对文件的所有权
								if helper.VerifyUserfile(p, strconv.FormatInt(uid, 10)) {
									if err := os.Remove(p); err != nil {
										log.Println("#", k, ",DEL FILE ERROR:", err)
									}

								}
							}

						}
					}
				}

			}
			//不管实际路径中是否存在文件均删除该数据库记录，以免数据库记录陷入死循环无法删掉
			if node.Id == id {

				if row, err := Engine.Id(id).Delete(new(Node)); err != nil || row == 0 {
					return errors.New(fmt.Sprint("删除话题错误!", row, err)) //错误还要我自己构造?!
				} else {
					return nil
				}

			}
		}
		return errors.New("你无权删除此话题:" + strconv.FormatInt(id, 10))
	}
	return errors.New("无法删除不存在的NODE ID:" + strconv.FormatInt(id, 10))
}

func GetNodesCount(offset int, limit int) (int64, error) {
	total, err := Engine.Limit(limit, offset).Count(&Node{})
	return total, err
}

func PostNode(node *Node) (int64, error) {
	return Engine.Insert(node)
}
