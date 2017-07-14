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

type Link struct {
	Id             int64
	Pid            int64 `xorm:"index"` //小于等于0 代表链接本身是顶层链接, pid大于0 代表属于某节点id的下级 本系统围绕node设计,链接亦可以是节点的下级
	Uid            int64 `xorm:"index"`
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
	Hotvote        int64   `xorm:"index"` //Hotup  + 	Hotdown
	Views          int64
	Author         string `xorm:"index"` //这里指本链接创建者
	Parent         string `xorm:"index"` //父级节点名称
	NodeTime       int64
	NodeCount      int64
	NodeLastUserId int64
}

func AddLink(title string, content string, attachment string, nid int64, parent string) (int64, error) {
	cat := &Link{Pid: nid, Parent: parent, Title: title, Content: content, Attachment: attachment, Created: time.Now().Unix()}

	if _, err := Engine.Insert(cat); err == nil {

		return cat.Id, err
	} else {
		return -1, err
	}

}

func GetLinks(offset int, limit int, field string) (*[]*Link, error) {
	linkz := new([]*Link)
	err := Engine.Table("link").Limit(limit, offset).Desc(field).Find(linkz)
	return linkz, err
}

func GetLinksByNodeCount(offset int, limit int, nodecount int64, field string) (*[]*Link, error) {
	linkz := new([]*Link)
	err := Engine.Table("link").Where("node_count > ?", nodecount).Limit(limit, offset).Desc(field).Find(linkz)
	return linkz, err
}

func GetLink(id int64) (*Link, error) {

	cat := &Link{}
	has, err := Engine.Id(id).Get(cat)

	if has {
		return cat, err
	} else {

		return nil, err
	}
}

func GetLinkByTitle(title string) (*Link, error) {
	cat := &Link{}
	cat.Title = title
	has, err := Engine.Get(cat)
	if has {
		return cat, err
	} else {

		return nil, err
	}
}

func PutLink(cid int64, cat *Link) (int64, error) {
	//覆盖式更新
	row, err := Engine.Update(cat, &Link{Id: cid})
	return row, err

}

//map[string]interface{}{"ctype": ctype}
func UpdateLink(cid int64, catmap *map[string]interface{}) error {
	cat := &Link{}
	if row, err := Engine.Table(cat).Where("id=?", cid).Update(catmap); err != nil || row == 0 {
		return errors.New(fmt.Sprint("UpdateLink row:", row, "出现错误:", err))
	} else {
		return nil
	}

}

func DelLink(id int64, uid int64, role int64) error {
	allow := false
	if role < 0 {
		allow = true
	}

	link := &Link{}

	if has, err := Engine.Id(id).Get(link); has == true && err == nil {

		if link.Uid == uid || allow {
			//检查附件字段并尝试删除文件
			if len(link.Attachment) > 0 {

				if p := helper.URL2local(link.Attachment); helper.Exist(p) {
					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误日志，但不要反回错误，以免陷入死循环无法删掉
							fmt.Println("ROOT DEL LINK Attachment, LINK ID:", id, ",ERR:", err)
						}
					} else { //检查用户对本地文件的所有权
						if helper.VerifyUserfile(p, strconv.FormatInt(uid, 10)) {
							if err := os.Remove(p); err != nil {
								fmt.Println("DEL LINK Attachment, LINK ID:", id, ",ERR:", err)
							}
						}
					}

				}
			}

			//检查内容字段并尝试删除文件
			if len(link.Content) > 0 {
				//若内容中存在图片则开始尝试删除图片
				delfiles_local := []string{}

				if m, n := helper.GetImages(link.Content); n > 0 {

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
			if link.Id == id {

				if row, err := Engine.Id(id).Delete(new(Link)); err != nil || row == 0 {
					return errors.New(fmt.Sprint("删除链接错误!", row, err))
				} else {
					return nil
				}

			}
		}
		return errors.New("你无权删除此链接:" + strconv.FormatInt(id, 10))
	}
	return errors.New("无法删除不存在的LINK ID:" + strconv.FormatInt(id, 10))
}

func GetLinkCountByPid(pid int64) int64 {
	n, _ := Engine.Where("pid=?", pid).Count(&Link{Pid: pid})
	return n
}
