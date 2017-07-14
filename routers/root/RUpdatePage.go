package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRUpdatePageHandler(self *makross.Context) error {

	self.Set("catpage", "RUpdatePageHandler")
	TplNames := "root/update_page"

	pageid := self.Param("pageid").MustInt64()

	page, _ := models.GetPage(pageid)

	self.Set("page", page)
	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", nds)

	}

	return self.Render(TplNames)

}

func PostRUpdatePageHandler(self *makross.Context) error {
	

	self.Set("catpage", "RUpdatePageHandler")
	//TplNames := "root/update_page"

	pageid := self.Param("pageid").MustInt64()
	nid := self.Args("nodeid").MustInt64()

	images := self.FormValue("images")
	title := self.FormValue("title")

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", nds)
	}

	if len(title) > 0 && len(content) > 0 && pageid > 0 {
		pagez, _ := models.GetPage(pageid)

		pagez.Title = title
		pagez.Content = content
		pagez.Attachment = images
		//nd, _ := models.GetNode(nid)
		pagez.Pid = nid
		models.PutPage(pageid, pagez)
		self.Set("pagez", pagez)
		self.Flash.Success(fmt.Sprint("更新 页面ID[", pageid, "] 成功!"))
		return self.Redirect("/root/update/page/" + strconv.FormatInt(pageid, 10) + "/")
	} else {
		self.Flash.Error(fmt.Sprint("更新 页面ID[", pageid, "] 失败!"))
		return self.Redirect("/root/read/page/")

	}

}
