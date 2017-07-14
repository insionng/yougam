package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRReadReplyHandler(self *makross.Context) error {
	

	TplNames := "root/update_reply"
	self.Set("catpage", "RReadReplyHandler")

	switch cmid := self.Param("cmid").MustInt64(); {
	//单独模式
	case cmid > 0:
		{
			if rp, err := models.GetReply(cmid); rp != nil && err == nil {
				self.Set("reply", *rp)

			} else {
				self.Flash.Error(err.Error())
				return self.Render(TplNames)

			}
		}

	//列表模式
	case cmid <= 0:
		{
			TplNames = "root/reply_table"
			offset := self.Args("offset").MustInt64()
			limit := self.Args("limit").MustInt64()
			field := self.FormValue("field")
			ctype := self.Args("ctype").MustInt64()
			pid := self.Args("pid").MustInt64()

			if limit == 0 {
				limit = 1000 //默认限制显示最近1000条,需要显示全部请在提交请求的时候设置limit字段为-1
			}

			if field == "" {
				field = "id"
			}

			if rps := models.GetReplysByPid(pid, ctype, int(offset), int(limit), field); rps != nil {
				self.Set("replys", rps)
				return self.Render(TplNames)

			}
		}
	}

	return self.Render(TplNames)
}

func PostRReadReplyHandler(self *makross.Context) error {
	

	delrowids := self.FormValue("delrowids")
	iserr := false
	//删除评论
	if delrowids != "" {

		delids := helper.Split(delrowids, ",")
		for _, delid := range delids {
			rid, _ := strconv.ParseInt(delid, 10, 0)

			if rid > 0 {
				if rp, err := models.GetReply(rid); rp != nil && err == nil {

					if rp.Tid > 0 {

						if rps := models.GetReplysByTid(rp.Tid, 0, 0, 0, "id"); rps != nil {
							for _, g := range *rps {
								models.DelReply(g.Id)
							}

							self.Flash.Success("删除 Reply id:" + delrowids + "成功!")

						}

					}
				} else {

					if e := models.DelReply(rid); e != nil {
						iserr = true
						self.Flash.Error("删除 Reply id:" + strconv.FormatInt(rid, 10) + "出现错误 " + fmt.Sprintf("%s", e) + "!")

					} else {
						self.Flash.Success("删除 Reply id:" + strconv.FormatInt(rid, 10) + "成功!")

					}
				}
			}
		}

		if iserr == false {
			self.Flash.Success("删除 Reply id:" + delrowids + "成功!")
		}
	} else {
		self.Flash.Error("非法 Reply id!")

	}

	return self.Redirect("/root/read/reply/")

}
