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

type Category struct {
	Id             int64
	Pid            int64 `xorm:"index"` //小于等于0 代表分类本身是顶层分类, pid大于0 代表属于某分类id的下级
	Nid            int64 `xorm:"index"` //小于等于0 代表本分类不属于某节点,大于0 代表属于某节点id的下级 本系统围绕node设计,分类亦可以是节点的下级
	Uid            int64 `xorm:"index"`
	Sort           int64
	Ctype          int64   `xorm:"index"`
	Title          string  `xorm:"index"`
	Content        string  `xorm:"text"`
	Attachment     string  `xorm:"text"`
	Created        int64   `xorm:"created index"`
	Updated        int64   `xorm:"updated"`
	Hotness        float64 `xorm:"index"`
	Confidence     float64 `xorm:"index"` //信任度数值
	Hotup          int64   `xorm:"index"`
	Hotdown        int64   `xorm:"index"`
	Hotscore       int64   `xorm:"index"` //Hotup  -	Hotdown
	Hotvote        int64   `xorm:"index"` //Hotup  + 	Hotdown
	Views          int64
	Author         string `xorm:"index"` //这里指本分类创建者
	Parent         string `xorm:"index"` //父级分类名称
	NodeTime       int64
	NodeCount      int64
	NodeLastUserId int64
	Template       string `xorm:"index"`
}

func GetCategoriesCount(offset int, limit int) (int64, error) {
	total, err := Engine.Limit(limit, offset).Count(&Category{})
	return total, err
}

func SetCategory(cid int64, cat *Category) (int64, error) {
	cat.Id = cid
	return Engine.Insert(cat)
}

func PostCategory(cat *Category) (int64, error) {
	return Engine.Insert(cat)
}

func AddCategory(title string, content string, attachment string, nid int64) (int64, error) {
	cat := &Category{Pid: nid, Title: title, Content: content, Attachment: attachment, Created: time.Now().Unix()}
	if _, err := Engine.Insert(cat); err == nil {
		return cat.Id, err
	} else {
		return -1, err
	}

}

func HasCategory(title string) (int64, bool) {
	cat := new(Category)
	if has, err := Engine.Where("title = ?", title).Get(cat); err != nil {
		return -1, false
	} else {
		if has {
			return cat.Id, true
		}
		return -1, false
	}
}

func GetCategories(offset int, limit int, field string) (*[]*Category, error) {
	catz := new([]*Category)
	err := Engine.Table("category").Limit(limit, offset).Desc(field).Find(catz)
	return catz, err
}

func GetCategoriesByNodeCount(offset int, limit int, nodecount int64, field string) (*[]*Category, error) {
	catz := new([]*Category)
	err := Engine.Table("category").Where("node_count > ?", nodecount).Limit(limit, offset).Desc(field).Find(catz)
	return catz, err
}

func GetCategoriesViaPid(pid int64, offset int, limit int, ctype int64, field string) *[]*Category {
	var objs = new([]*Category)
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

func GetCategoriesByPid(pid int64, offset int, limit int, ctype int64, field string) *[]*Category {
	var objs = new([]*Category)
	switch {
	case field == "asc":
		{
			if ctype != 0 {
				Engine.Where("(pid=? or id=?) and ctype=?", pid, pid, ctype).Limit(limit, offset).Asc("id").Find(objs)
			} else {
				Engine.Where("pid=? or id=?", pid, pid).Limit(limit, offset).Asc("id").Find(objs)
			}
		}
	default:
		{
			if ctype != 0 {
				Engine.Where("(pid=? or id=?) and ctype=?", pid, pid, ctype).Limit(limit, offset).Desc(field).Find(objs)
			} else {
				Engine.Where("pid=? or id=?", pid, pid).Limit(limit, offset).Desc(field).Find(objs)
			}
		}
	}
	return objs
}

func GetCategoriesByCtypeWithNid(offset int, limit int, ctype int64, nid int64, field string) (*[]*Category, error) {
	objs := new([]*Category)
	err := Engine.Table("category").Where("ctype = ? and nid = ?", ctype, nid).Limit(limit, offset).Desc(field).Find(objs)
	return objs, err
}

func GetCategoriesByNid(nid int64, offset int, limit int, field string) (*[]*Category, error) {
	nds := new([]*Category)
	err := Engine.Where("nid=?", nid).Limit(limit, offset).Desc(field).Find(nds)
	return nds, err
}

func GetCategoriesByCtype(offset int, limit int, ctype int64, field string) (*[]*Category, error) {
	catz := new([]*Category)
	err := Engine.Table("category").Where("ctype = ?", ctype).Limit(limit, offset).Desc(field).Find(catz)
	return catz, err
}

func GetCategoriesByCtypeWithPid(offset int, limit int, ctype int64, pid int64, field string) (*[]*Category, error) {
	catz := new([]*Category)
	var err error
	if pid <= 0 {
		err = Engine.Table("category").Where("ctype = ? and pid <= ?", ctype, 0).Limit(limit, offset).Desc(field).Find(catz)
	} else {
		err = Engine.Table("category").Where("ctype = ? and pid = ?", ctype, pid).Limit(limit, offset).Desc(field).Find(catz)
	}
	return catz, err
}

func GetCategory(id int64) (*Category, error) {

	cat := &Category{}
	has, err := Engine.Id(id).Get(cat)

	if has {
		return cat, err
	} else {
		return nil, err
	}
}

func GetCategoryByTitle(title string) (*Category, error) {
	var cat = &Category{Title: title}
	has, err := Engine.Get(cat)
	if has {
		return cat, err
	} else {
		return nil, err
	}
}

func PutCategory(cid int64, cat *Category) (int64, error) {
	//覆盖式更新
	row, err := Engine.Update(cat, &Category{Id: cid})
	return row, err
}

//map[string]interface{}{"ctype": ctype}
func UpdateCategory(cid int64, catmap *map[string]interface{}) error {
	cat := &Category{}
	if row, err := Engine.Table(cat).Where("id=?", cid).Update(catmap); (err != nil) || (row <= 0) {
		return errors.New(fmt.Sprint("UpdateCategory row:", row, "出现错误:", err))
	} else {
		return nil
	}

}

func DelCategory(id int64, uid int64, role int64) error {
	allow := false
	if role < 0 {
		allow = true
	}

	category := &Category{}

	if has, err := Engine.Id(id).Get(category); has == true && err == nil {

		if category.Uid == uid || allow {
			//检查附件字段并尝试删除文件
			if len(category.Attachment) > 0 {

				if p := helper.URL2local(category.Attachment); helper.Exist(p) {
					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误日志，但不要反回错误，以免陷入死循环无法删掉
							log.Println("ROOT DEL CATEGORY Attachment, CATEGORY ID:", id, ",ERR:", err)
						}
					} else { //检查用户对本地文件的所有权
						if helper.VerifyUserfile(p, strconv.FormatInt(uid, 10)) {
							if err := os.Remove(p); err != nil {
								log.Println("DEL CATEGORY Attachment, CATEGORY ID:", id, ",ERR:", err)
							}
						}
					}

				}
			}

			//检查内容字段并尝试删除文件
			if len(category.Content) > 0 {
				//若内容中存在图片则开始尝试删除图片
				delfiles_local := []string{}

				if m, n := helper.GetImages(category.Content); n > 0 {

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
			if category.Id == id {

				if row, err := Engine.Id(id).Delete(new(Category)); err != nil || row == 0 {
					return errors.New(fmt.Sprint("删除分类错误!", row, err)) //错误还要我自己构造?!
				} else {
					return nil
				}

			}
		}
		return errors.New("你无权删除此分类:" + strconv.FormatInt(id, 10))
	}
	return errors.New("无法删除不存在的CATEGORY ID:" + strconv.FormatInt(id, 10))
}

func GetCategoryCountByPid(pid int64) int64 {
	n, _ := Engine.Where("pid=?", pid).Count(&Category{Pid: pid})
	return n
}
