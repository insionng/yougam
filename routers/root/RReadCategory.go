package root

import (
	"github.com/insionng/makross"
	
	"github.com/insionng/yougam/models"
)

func GetRReadCategoryHandler(self *makross.Context) error {
	

	TplNames := ""
	self.Set("catpage", "RReadCategoryHandler")
	switch cid := self.Param("cid").MustInt64(); {
	//单独模式
	case cid > 0:
		{
			TplNames = "root/create_category"

			if cat, err := models.GetCategory(cid); cat != nil && err == nil {
				self.Set("category", *cat)

				if nodes, err := models.GetNodes(0, 0, "id"); nodes != nil && err == nil {
					self.Set("nodes", *nodes)
				}

			} else {
				self.Flash.Error(err.Error())
				return self.Render(TplNames)
			}
		}
	//列表模式
	case cid <= 0:
		{
			TplNames = "root/category_table"
			offset := self.Args("offset").MustInt64()
			limit := self.Args("limit").MustInt64()
			field := self.FormValue("field")

			if limit == 0 {
				limit = 1000 //默认限制显示最近1000条,需要显示全部请在提交请求的时候设置limit字段为-1
			}

			if field == "" {
				field = "id"
			}

			if cats, err := models.GetCategories(int(offset), int(limit), field); err == nil && cats != nil {
				self.Set("categories", *cats)
			}
		}
	}

	return self.Render(TplNames)

}
