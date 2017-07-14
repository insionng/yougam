package root

import (
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRUpdateCategoryHandler(self *makross.Context) error {
	self.Set("catpage", "RUpdateCategoryHandler")
	TplNames := "root/update_category"

	cid := self.Param("cid").MustInt64()

	category, _ := models.GetCategory(cid)

	self.Set("category", category)
	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", nds)

	}

	return self.Render(TplNames)

}

func PostRUpdateCategoryHandler(self *makross.Context) error {
	

	self.Set("catpage", "RUpdateCategoryHandler")
	//TplNames := "root/update_category"

	cid := self.Param("cid").MustInt64()
	nid := self.Args("nodeid").MustInt64()

	images := self.FormValue("images")
	title := self.FormValue("title")

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", nds)
	}

	if len(title) > 0 && cid > 0 {
		cat, _ := models.GetCategory(cid)

		cat.Title = title
		cat.Content = content
		cat.Attachment = images
		//nd, _ := models.GetNode(nid)
		cat.Pid = nid
		models.PutCategory(cid, cat)
		self.Set("cat", cat)

		self.Flash.Success("更新分类成功!")
		return self.Redirect("/root/update/category/" + strconv.FormatInt(cid, 10) + "/")
	} else {
		self.Flash.Error("更新分类失败,标题为空或非法请求!")
		return self.Redirect("/root/read/category/")
	}

}
