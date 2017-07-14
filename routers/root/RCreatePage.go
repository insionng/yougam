package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRCreatePageHandler(self *makross.Context) error {

	self.Set("catpage", "RCreatePageHandler")

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	TplNames := "root/create_page"
	return self.Render(TplNames)

}

func PostRCreatePageHandler(self *makross.Context) error {
	

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))
	title := self.FormValue("title")
	images := self.FormValue("images")
	nid := self.Args("nodeid").MustInt64()
	/*
		if nd, e := models.GetNode(nid); (nd == nil) || (e != nil) {
			self.Flash.Error("节点不存在!", false)
			return self.Redirect("/root/create/page/")
			return
		}
	*/

	if len(title) > 0 && len(content) > 0 {

		if cid, err := models.AddPage(title, content, images, nid); err != nil {
			self.Flash.Error(fmt.Sprint("增加页面出现错误:", err))
			return self.Redirect("/root/create/page/")

		} else {
			self.Flash.Success("新增页面成功!")
			return self.Redirect("/root/read/page/" + strconv.FormatInt(cid, 10) + "/")

		}
	} else {
		self.Flash.Error("页面标题及内容不能为空!")
		return self.Redirect("/root/create/page/")

	}
}
