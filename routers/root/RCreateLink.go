package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRCreateLinkHandler(self *makross.Context) error {

	self.Set("catpage", "RCreateLinkHandler")
	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	TplNames := "root/create_link"
	return self.Render(TplNames)

}

func PostRCreateLinkHandler(self *makross.Context) error {

	policy := helper.StandardURLsPolicy()
	content := policy.Sanitize(self.FormValue("content"))
	
	title := self.FormValue("title")
	images := self.FormValue("images")
	nid := self.Args("nodeid").MustInt64()
	nd, _ := models.GetNode(nid)
	var parentname string
	if nd != nil {
		parentname = nd.Title
	}

	if len(title) > 0 && len(content) > 0 {

		if cid, err := models.AddLink(title, content, images, nid, parentname); err != nil {

			self.Flash.Error(fmt.Sprint("增加友链出现错误:", err))
			return self.Redirect("/root/create/link/")

		} else {
			self.Flash.Success("新增友链成功!")
			return self.Redirect("/root/read/link/" + strconv.FormatInt(cid, 10) + "/")

		}
	} else {
		self.Flash.Error("友链标题及内容不能为空!")
		return self.Redirect("/root/create/link/")

	}

}
