package root

import (
	"github.com/insionng/makross"
	
	"github.com/insionng/yougam/models"
)

func GetRReadLinkHandler(self *makross.Context) error {
	

	TplNames := ""
	self.Set("catpage", "RReadLinkHandler")
	switch linkid := self.Param("linkid").MustInt64(); {
	//单独模式
	case linkid > 0:
		{
			TplNames = "root/create_link"

			if cat, err := models.GetLink(linkid); cat != nil && err == nil {
				self.Set("link", *cat)

				if nodes, err := models.GetNodes(0, 0, "id"); nodes != nil && err == nil {
					self.Set("nodes", nodes)
				}

			} else {
				self.Flash.Error(err.Error())
				return self.Render(TplNames)

			}
		}
	//列表模式
	case linkid <= 0:
		{
			TplNames = "root/link_table"
			offset := self.Args("offset").MustInt64()
			limit := self.Args("limit").MustInt64()
			field := self.FormValue("field")

			if limit == 0 {
				limit = 1000 //默认限制显示最近1000条,需要显示全部请在提交请求的时候设置limit字段为-1
			}

			if field == "" {
				field = "id"
			}

			if cats, err := models.GetLinks(int(offset), int(limit), field); err == nil && cats != nil {
				self.Set("links", cats)
			}
		}
	}

	return self.Render(TplNames)
}
