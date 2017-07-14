package root

import (
	"fmt"
	"github.com/insionng/makross"

	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRCreateTopicHandler(self *makross.Context) error {

	self.Set("haspid", false)

	self.Set("catpage", "RCreateTopicHandler")
	if pid := self.Param("pid").MustInt64(); pid > 0 {
		self.Set("haspid", true)
		self.Set("pid", pid)
	}

	if tid := self.Param("tid").MustInt64(); tid > 0 {
		if tp, err := models.GetTopic(tid); tp != nil && err == nil {
			self.Set("topic", *tp)
		} else {
			self.Flash.Error(err.Error())
			return err
		}
	}

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	TplNames := "root/create_topic"
	return self.Render(TplNames)

}

func PostRCreateTopicHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))

	title := self.FormValue("title")
	//images := self.FormValue("images")

	cid := self.Args("cid").MustInt64()
	nid := self.Args("nodeid").MustInt64()

	pid := self.Args("pid").MustInt64()
	if pid == 0 {
		pid = self.Param("pid").MustInt64()
	}

	if (len(content) > 0) && (_usr_.Id > 0) {
		if tid, err := models.AddTopic(title, content, _usr_.Avatar, _usr_.AvatarLarge, _usr_.AvatarMedium, _usr_.AvatarSmall, pid, cid, nid, _usr_.Id); err != nil || tid <= 0 {
			self.Flash.Error(fmt.Sprint("增加话题出现错误:", err))
			return self.Redirect("/root/create/topic/")

		} else {
			if pid == 0 {
				pid = tid
			}

			self.Flash.Success("新增话题成功!")
			return self.Redirect("/root/create/" + strconv.FormatInt(pid, 10) + "/topic/")
			//return self.Redirect( "/root/read/topic/"+strconv.FormatInt(tid, 10)+"/")

		}

	} else {
		return self.Redirect("/root/read/topic/")
	}

}
