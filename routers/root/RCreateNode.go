package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func RCreateNodeGetHandler(self *makross.Context) error {

	

	TplNames := "root/create_node"
	self.Set("catpage", "RCreateNodeHandler")

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	if categories, err := models.GetCategories(0, 0, "id"); categories != nil && err == nil {
		self.Set("categories", categories)
	}

	if nid := self.Param("nid").MustInt64(); nid > 0 {

		if nd, err := models.GetNode(nid); nd != nil && err == nil {
			self.Set("node", *nd)
		} else {
			self.Flash.Error(err.Error())
			return err
		}
	}

	return self.Render(TplNames)

}

func PostRCreateNodeHandler(self *makross.Context) error {

	

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))
	pid := self.Args("nodeid").MustInt64()
	title := self.FormValue("title")
	images := self.FormValue("images")

	cid := self.Args("cid").MustInt64()
	var uid int64
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if okay {
		uid = _usr_.Id
	}

	if (len(title) > 0) /* && (content != "") */ && (uid > 0) {

		if nid, err := models.AddNode(title, content, images, pid, cid, uid); err != nil && uid > 0 {

			self.Flash.Error(fmt.Sprint("增加节点出现错误:", err))
			return self.Redirect("/root/create/node/")

		} else {

			//新增点成功后 就去统计有多少个同样分类id的节点,把统计出来的数目写入该分类的nodecount项
			if cid > 0 {
				if nc, e := models.GetNodesByCid(cid, 0, 0, "id"); e == nil {
					if catz, e := models.GetCategory(cid); e == nil {
						catz.NodeCount = int64(len(*nc))
						models.PutCategory(cid, catz)
					}
				}
			}

			self.Flash.Success("新增节点成功!")
			return self.Redirect("/root/read/node/" + strconv.FormatInt(nid, 10) + "/")

		}
	} else {
		self.Flash.Error("节点最低要求标题不能为空！")
		return self.Redirect("/root/create/node/")

	}

}
