package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRUpdateReplyHandler(self *makross.Context) error {
	

	self.Set("catpage", "RUpdateReplyHandler")
	if rid := self.Param("rid").MustInt64(); rid > 0 {

		if rp, err := models.GetReply(rid); rp != nil && err == nil {
			self.Set("reply", *rp)
		} else {
			self.Flash.Error(err.Error())
			return err
		}
	}

	TplNames := "root/update_reply"
	return self.Render(TplNames)

}

func PostRUpdateReplyHandler(self *makross.Context) error {
	

	self.Set("catpage", "RUpdateReplyHandler")
	TplNames := "root/update_reply"
	rid := self.Param("rid").MustInt64()
	images := self.FormValue("images")

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))

	if rid > 0 {
		if rp, err := models.GetReply(rid); rp != nil && err == nil {

			rp.Content = content
			if len(images) > 0 {
				rp.Attachment = images
			}
			if _, err := models.PutReply(rid, rp); err != nil {
				self.Flash.Error(fmt.Sprint("更新回复出现错误:", err.Error()))
			} else {
				self.Flash.Success("更新回复成功!")

			}

		}

		return self.Redirect("/root/read/reply/" + strconv.FormatInt(rid, 10) + "/")

	} else {
		self.Flash.Error("更新回复出现错误:不存在该回复!")
		return self.Redirect("/root/read/reply/")

	}
	return self.Render(TplNames)

}
