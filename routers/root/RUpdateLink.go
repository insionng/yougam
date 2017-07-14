package root

import (
	"fmt"
	"github.com/insionng/makross"

	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRUpdateLinkHandler(self *makross.Context) error {
	self.Set("catpage", "RUpdateLinkHandler")
	TplNames := "root/update_link"

	linkid := self.Param("linkid").MustInt64()

	link, _ := models.GetLink(linkid)

	self.Set("link", link)
	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", nds)

	}

	return self.Render(TplNames)

}

func PostRUpdateLinkHandler(self *makross.Context) error {

	self.Set("catpage", "RUpdateLinkHandler")
	//TplNames := "root/update_link"

	linkid := self.Param("linkid").MustInt64()
	nid := self.Args("nodeid").MustInt64()

	images := self.FormValue("images")
	title := self.FormValue("title")

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))

	if nds, err := models.GetNodes(0, 0, "id"); nds != nil && err == nil {
		self.Set("nodes", nds)
	}

	if len(title) > 0 && len(content) > 0 && linkid > 0 {
		linkz, _ := models.GetLink(linkid)

		linkz.Title = title
		linkz.Content = content
		linkz.Attachment = images
		nd, _ := models.GetNode(nid)
		if nd != nil {
			linkz.Parent = nd.Title
		}
		linkz.Pid = nid
		models.PutLink(linkid, linkz)
		self.Set("linkz", linkz)
		self.Flash.Success(fmt.Sprint("更新 友链ID[", linkid, "] 成功!"))
		return self.Redirect("/root/update/link/" + strconv.FormatInt(linkid, 10) + "/")
	} else {
		self.Flash.Error(fmt.Sprint("更新 友链ID[", linkid, "] 失败!"))
		return self.Redirect("/root/read/link/")

	}

}
