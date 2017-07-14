package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func RCreateCategoryGetHandler(self *makross.Context) error {

	self.Set("catpage", "RCreateCategoryHandler")

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	TplNames := "root/create_category"
	return self.Render(TplNames)

}

func RCreateCategoryPostHandler(self *makross.Context) error {
	

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))

	title := self.FormValue("title")
	images := self.FormValue("images")
	nid := self.Args("nodeid").MustInt64()
	/*
		if nd, e := models.GetNode(nid); (nd == nil) || (e != nil) {
			self.Flash.Error("节点不存在!", false)
			return self.Redirect("/root/create/category/")
			return
		}
	*/

	if len(title) > 0 {

		if cid, err := models.AddCategory(title, content, images, nid); err != nil {
			self.Flash.Error(fmt.Sprint("增加分类出现错误:", err))
			return self.Redirect("/root/create/category/")

		} else {
			self.Flash.Success("新增分类成功!")
			return self.Redirect("/root/read/category/" + strconv.FormatInt(cid, 10) + "/")

		}
	} else {
		self.Flash.Error("分类标题绝不能为空!")
		return self.Redirect("/root/create/category/")

	}

}
