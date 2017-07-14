package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRReadUserHandler(self *makross.Context) error {
	

	TplNames := ""
	self.Set("catpage", "RReadUserHandler")
	keyword := self.Param("keyword").String()
	if keyword == "" {
		keyword = self.FormValue("keyword")
	}

	switch uid := self.Param("uid").MustInt64(); {
	//单独模式
	case uid > 0:
		{
			TplNames = "root/create_user"

			if usr, err := models.GetUser(uid); usr != nil && err == nil {
				self.Set("user", *usr)

			} else {
				self.Flash.Error(err.Error())
				return self.Render(TplNames)

			}
		}
	//列表模式
	case uid <= 0:
		{
			TplNames = "root/user_table"
			offset := self.Args("offset").MustInt64()
			limit := self.Args("limit").MustInt64()
			field := self.FormValue("field")

			if limit == 0 {
				limit = 1000 //默认限制显示最近1000条,需要显示全部请在提交请求的时候设置limit字段为-1
			}

			if field == "" {
				field = "id"
			}

			if keyword != "" {
				if usrs, err := models.SearchUser(keyword, int(offset), int(limit), field); err == nil && usrs != nil {
					self.Set("users", *usrs)
					self.Set("keyword", keyword)
				}

			} else {

				if usrs, err := models.GetUsersOnHotness(int(offset), int(limit), field); err == nil && usrs != nil {
					self.Set("users", *usrs)
				}
			}

		}
	}

	return self.Render(TplNames)

}

func PostRReadUserHandler(self *makross.Context) error {
	

	delrowids := self.FormValue("delrowids")

	//删除用户
	if len(delrowids) > 0 {

		uid := int64(0)
		role := int64(-1000)
		iserr := false
		delids := helper.Split(delrowids, ",")
		for _, delid := range delids {
			userid, _ := strconv.ParseInt(delid, 10, 0)

			if e := models.DelUser(userid, uid, role); e != nil {
				iserr = true
				self.Flash.Error("删除 User id:" + strconv.FormatInt(userid, 10) + "出现错误 " + fmt.Sprintf("%s", e) + "!")

			} else {
				self.Flash.Success("删除 User id:" + strconv.FormatInt(userid, 10) + "成功!")

			}

		}
		if iserr == false {
			self.Flash.Success("删除 User id:" + delrowids + "成功!")

		}
	} else {
		self.Flash.Error("非法 User id!")

	}

	return self.Redirect("/root/read/user/")

}
