package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRUpdateNodeHandler(self *makross.Context) error {
	

	TplNames := "root/update_node"
	self.Set("catpage", "RUpdateNodeHandler")

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	if categories, err := models.GetCategories(0, 0, "id"); categories != nil && err == nil {
		self.Set("categories", categories)
	}

	if nid := self.Param("nid").MustInt64(); nid > 0 {

		if nd, err := models.GetNode(nid); nd != nil && err == nil {
			self.Set("node", *nd)
			if category, e := models.GetCategory(nd.Cid); category != nil && e == nil {
				self.Set("category", *category)
			}

		} else {
			self.Flash.Error(err.Error())
			return err
		}
	}

	return self.Render(TplNames)

}

func PostRUpdateNodeHandler(self *makross.Context) error {

	self.Set("catpage", "RUpdateNodeHandler")
	TplNames := "root/update_node"

	nid := self.Param("nid").MustInt64()
	pid := self.Args("nodeid").MustInt64()
	title := self.FormValue("title")
	images := self.FormValue("images")
	cid := self.Args("cid").MustInt64()

	
	var uid int64
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if okay {
		uid = _usr_.Id
	}
	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	if nd, err := models.GetNode(nid); nd != nil && err == nil {
		nd.Id = nid
		nd.Title = title
		nd.Content = content
		nd.Attachment = images
		nd.Pid = pid
		nd.Cid = cid
		nd.Uid = uid

		if row, err := models.PutNode(nid, nd); err != nil || row == 0 {
			self.Flash.Error(fmt.Sprintf("更新节点出现错误:%v", err))

		} else {

			//更新节点成功后 就去统计有多少个同样分类id的节点,把统计出来的数目写入该分类的nodecount项
			if cid > 0 {
				if nc, e := models.GetNodesByCid(cid, 0, 0, "id"); e == nil {
					if catz, e := models.GetCategory(cid); e == nil {
						catz.NodeCount = int64(len(*nc))
						models.PutCategory(cid, catz)
					}
				}
			}

			self.Flash.Success("更新节点成功!")

		}
		return self.Redirect("/root/read/node/" + strconv.FormatInt(nid, 10) + "/")

	} else {
		self.Flash.Error("获取节点出现错误:不存在该节点!")
		return self.Redirect("/root/read/node/")

	}

	return self.Render(TplNames)

}
