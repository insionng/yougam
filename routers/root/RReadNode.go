package root

import (
	"github.com/insionng/makross"
	"github.com/insionng/yougam/models"
)

func GetRReadNodeHandler(self *makross.Context) error {

	TplNames := ""
	self.Set("catpage", "RReadNodeHandler")
	switch nid := self.Param("nid").MustInt64(); {
	//单独模式
	case nid > 0:
		{
			TplNames = "root/read_node"

			if nd, err := models.GetNode(nid); nd != nil && err == nil {
				self.Set("node", *nd)

				if nodes, err := models.GetNodes(0, 0, "id"); nodes != nil && err == nil {
					self.Set("nodes", *nodes)
				}

			} else {
				self.Flash.Error(err.Error())
				return self.Render(TplNames)

			}
		}
	//列表模式
	case nid <= 0:
		{
			TplNames = "root/node_table"
			offset := self.Args("offset").MustInt64()
			limit := self.Args("limit").MustInt64()
			field := self.FormValue("field")

			if limit == 0 {
				limit = 1000 //默认限制显示最近1000条,需要显示全部请在提交请求的时候设置limit字段为-1
			}

			if field == "" {
				field = "id"
			}

			if nds, err := models.GetNodes(int(offset), int(limit), field); err == nil && nds != nil {
				self.Set("nodes", *nds)
			}
		}
	}

	return self.Render(TplNames)

}
