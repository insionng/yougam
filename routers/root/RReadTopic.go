package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRReadTopicHandler(self *makross.Context) error {
	

	TplNames := "root/create_topic"
	self.Set("catpage", "RReadTopicHandler")

	switch tid := self.Param("tid").MustInt64(); {
	//单独模式
	case tid > 0:
		{
			TplNames = "root/read_topic"
			if tp, err := models.GetTopic(tid); tp != nil && err == nil {
				self.Set("topic", *tp)

				if nodes, err := models.GetNodes(0, 0, "id"); nodes != nil && err == nil {
					self.Set("nodes", *nodes)
				}

			} else {
				self.Flash.Error(fmt.Sprint(err))
				return self.Render(TplNames)

			}
		}

	//列表模式
	case tid <= 0:
		{
			TplNames = "root/topic_table"
			self.Set("catpage", "RReadTopicHandlerlist")

			offset := self.Args("offset").MustInt64()
			limit := self.Args("limit").MustInt64()
			field := self.FormValue("field")

			if limit == 0 {
				limit = 1000 //默认限制显示最近1000条,需要显示全部请在提交请求的时候设置limit字段为-1
			}

			if field == "" {
				field = "id"
			}

			if pid := self.Param("pid").MustInt64(); pid > 0 {

				if tps := models.GetTopicsByPid(pid, int(offset), int(limit), 0, field); tps != nil && (len(*tps) > 0) {
					self.Set("catpage", "RReadTopicHandler")
					self.Set("topics", *tps)
				}

			} else {

				if tps := models.GetSubjectsByNid(0, int(offset), int(limit), 0, field); tps != nil {
					self.Set("topics", *tps)
				}
			}

		}
	}

	return self.Render(TplNames)
}

func PostRReadTopicHandler(self *makross.Context) error {
	

	delrowids := self.FormValue("delrowids")

	//删除话题
	if len(delrowids) > 0 {

		uid := int64(0)
		role := int64(-1000)

		delids := helper.Split(delrowids, ",")
		for _, delid := range delids {
			tid, _ := strconv.ParseInt(delid, 10, 0)

			if tid > 0 {
				if tp, err := models.GetTopic(tid); tp != nil && err == nil {
					if tp.Pid == 0 {

						if tps := models.GetTopicsByPid(tid, 0, 0, 0, "id"); tps != nil {
							for _, g := range *tps {
								models.DelTopic(g.Id, uid, role)

							}

							self.Flash.Success("删除 Topic id:" + delrowids + "成功!")

						}
					}
				} else {

					if e := models.DelTopic(tid, uid, role); e != nil {
						self.Flash.Error(fmt.Sprintf("删除 Topic id:", strconv.FormatInt(tid, 10), "出现错误 ", e, "!"))

					} else {
						self.Flash.Success("删除 Topic id:" + strconv.FormatInt(tid, 10) + "成功!")

					}
				}
			}
		}

	} else {
		self.Flash.Error("非法 Topic id!")

	}

	return self.Redirect("/root/read/topic/")

}
