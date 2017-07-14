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

type Page struct {
	Id         int64
	Pid        int64 `xorm:"index"` //小于等于0 代表页面本身是顶层页面, pid大于0 代表属于某节点id的下级 本系统围绕node设计,页面亦可以是节点的下级
	Uid        int64 `xorm:"index"`
	Sort       int64
	Ctype      int64   `xorm:"index"`
	Title      string  `xorm:"index"`
	Content    string  `xorm:"text"`
	Attachment string  `xorm:"text"`
	Created    int64   `xorm:"created index"`
	Updated    int64   `xorm:"updated"`
	Hotness    float64 `xorm:"index"`
	Hotup      int64   `xorm:"index"`
	Hotdown    int64   `xorm:"index"`
	Hotscore   int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote    int64   `xorm:"index"` //Hotup  + 	Hotdown
	Views      int64
	Author     string `xorm:"index"` //这里指本页面创建者
	Template   string `xorm:"index"`
}

func AddPage(title string, content string, attachment string, nid int64) (int64, error) {
	page := &Page{Pid: nid, Title: title, Content: content, Attachment: attachment, Created: time.Now().Unix()}

	if _, err := Engine.Insert(page); err == nil {

		return page.Id, err
	} else {
		return -1, err
	}

}

func SetPage(pageid int64, p *Page) (int64, error) {
	p.Id = pageid
	return Engine.Insert(p)
}

func GetPages(offset int, limit int, field string) (*[]*Page, error) {
	pagez := new([]*Page)
	err := Engine.Table("page").Limit(limit, offset).Desc(field).Find(pagez)
	return pagez, err
}

func GetPagesByNodeCount(offset int, limit int, nodecount int64, field string) (*[]*Page, error) {
	pagez := new([]*Page)
	err := Engine.Table("page").Where("node_count > ?", nodecount).Limit(limit, offset).Desc(field).Find(pagez)
	return pagez, err

}

func GetPage(id int64) (*Page, error) {

	obj := &Page{}
	has, err := Engine.Id(id).Get(obj)

	if has {
		return obj, err
	} else {

		return nil, err
	}
}

 
func GetPageByTitle(title string) (*Page, error) {
	obj := &Page{}
	obj.Title = title
	has, err := Engine.Get(obj)
	if has {
		return obj, err
	} else {

		return nil, err
	}
}

func PutPage(cid int64, obj *Page) (int64, error) {
	//覆盖式更新
	row, err := Engine.Update(obj, &Page{Id: cid})
	return row, err

}

//map[string]interface{}{"ctype": ctype}
func UpdatePage(cid int64, objmap *map[string]interface{}) error {
	obj := &Page{}
	if row, err := Engine.Table(obj).Where("id=?", cid).Update(objmap); err != nil || row == 0 {
		return errors.New(fmt.Sprint("UpdatePage row:", row, "出现错误:", err))
	} else {
		return nil
	}

}

func DelPage(id int64, uid int64, role int64) error {
	allow := false
	if role < 0 {
		allow = true
	}

	page := &Page{}

	if has, err := Engine.Id(id).Get(page); has == true && err == nil {

		if page.Uid == uid || allow {
			//检查附件字段并尝试删除文件
			if len(page.Attachment) > 0 {

				if p := helper.URL2local(page.Attachment); helper.Exist(p) {
					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误日志，但不要反回错误，以免陷入死循环无法删掉
							log.Println("ROOT DEL PAGE Attachment, PAGE ID:", id, ",ERR:", err)
						}
					} else { //检查用户对本地文件的所有权
						if helper.VerifyUserfile(p, strconv.FormatInt(uid, 10)) {
							if err := os.Remove(p); err != nil {
								log.Println("DEL PAGE Attachment, PAGE ID:", id, ",ERR:", err)
							}
						}
					}

				}
			}

			//检查内容字段并尝试删除文件
			if len(page.Content) > 0 {
				//若内容中存在图片则开始尝试删除图片
				delfiles_local := []string{}

				if m, n := helper.GetImages(page.Content); n > 0 {

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
			if page.Id == id {

				if row, err := Engine.Id(id).Delete(new(Page)); err != nil || row == 0 {
					return fmt.Errorf("E:删除页面错误 %v !", err.Error())
				}
				return nil
			}
		}
		return errors.New("你无权删除此页面:" + strconv.FormatInt(id, 10))
	}
	return errors.New("无法删除不存在的PAGE ID:" + strconv.FormatInt(id, 10))
}

func GetPageCountByPid(pid int64) int64 {
	n, _ := Engine.Where("pid=?", pid).Count(&Page{Pid: pid})
	return n
}
