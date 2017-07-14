package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRUpdateTopicHandler(self *makross.Context) error {
	

	self.Set("catpage", "RUpdateTopicHandler")
	if tid := self.Param("tid").MustInt64(); tid > 0 {
		if tp, err := models.GetTopic(tid); tp != nil && err == nil {
			self.Set("topic", *tp)
			self.Set("images", tp.Attachment)
		} else {
			self.Flash.Error(err.Error())
			return err
		}
	}

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	TplNames := "root/update_topic"
	return self.Render(TplNames)

}

func PostRUpdateTopicHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	self.Set("catpage", "RUpdateTopicHandler")
	TplNames := "root/update_topic"
	tid := self.Param("tid").MustInt64()

	title := self.FormValue("title")
	images := self.FormValue("images")

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))

	cid := self.Args("cid").MustInt64()
	nid := self.Args("nodeid").MustInt64()
	pid := self.Args("pid").MustInt64()
	uid := _usr_.Id //当前管理员uid
	if tid > 0 {
		if tp, err := models.GetTopic(tid); tp != nil && err == nil {
			tp.Title = title
			tp.Content = content
			tp.Attachment = images
			tp.Pid = pid
			tp.Cid = cid
			tp.Nid = nid
			if tp.Uid <= 0 {
				tp.Uid = uid //uid默认不作修改,不然会把用户的uid替换掉,当uid为0才设为管理员uid
			}
			if nd, e := models.GetNode(nid); nd != nil && e == nil {
				tp.Node = nd.Title
			}

			if tid, err := models.PutTopic(tid, tp); err != nil || tid <= 0 {
				self.Flash.Error(fmt.Sprint("更新话题出现错误:", err))

			} else {
				self.Flash.Success("更新话题成功!")
			}

		}

		return self.Redirect("/root/read/topic/" + strconv.FormatInt(tid, 10) + "/")

	} else {
		self.Flash.Error("更新话题出现错误:不存在该话题!")
		return self.Redirect("/root/read/node/")

	}
	return self.Render(TplNames)

}
