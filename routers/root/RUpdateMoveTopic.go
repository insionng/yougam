package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"github.com/insionng/yougam/models"
)

func GetRMoveTopicHandler(self *makross.Context) error {

	

	TplNames := "root/move_topic"
	self.Set("catpage", "RMoveTopicHandler")
	tid := self.Param("tid").MustInt64()
	if tid <= 0 {
		self.Flash.Error("非法参数!")
		return self.Redirect("/root/read/topic/")

	}

	if tp, err := models.GetTopic(tid); tp != nil && err == nil {
		self.Set("tp", *tp)

		if nodes, err := models.GetNodes(0, 0, "id"); nodes != nil && err == nil {
			self.Set("nodes", *nodes)
		}

	} else {
		self.Flash.Error(fmt.Sprint(err))
		return self.Render(TplNames)

	}
	return self.Render(TplNames)
}

func PostRMoveTopicHandler(self *makross.Context) error {
	

	nodeid := self.Args("nodeid").MustInt64()

	tid := self.Param("tid").MustInt64()
	if tid <= 0 || nodeid <= 0 {
		self.Flash.Error("非法参数!")
		return self.Redirect("/root/read/topic/")

	}

	tp, e := models.GetTopic(tid)
	if e != nil || tp == nil {
		self.Flash.Error("话题不存在!")
		return self.Redirect("/root/read/topic/")

	}

	node, e := models.GetNode(nodeid)
	if e != nil || node == nil {
		self.Flash.Error("节点不存在!")
		return self.Redirect("/root/read/topic/")
	}

	tp.Nid = node.Id
	tp.Node = node.Title
	if _, e := models.PutTopic(tid, tp); e == nil {
		self.Flash.Success(fmt.Sprintf("移动话题[%v]成功!", tp.Title))
	} else {
		self.Flash.Error(fmt.Sprintf("移动话题[%v]失败!", tp.Title))
	}
	return self.Redirect(fmt.Sprintf("/root/update/topic/move/%v/", tid))

}
