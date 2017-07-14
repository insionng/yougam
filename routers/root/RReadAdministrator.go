package root

import (
	"github.com/insionng/makross"
	
	"github.com/insionng/yougam/models"
)

func GetRReadRootHandler(self *makross.Context) error {
	

	TplNames := ""
	self.Set("catpage", "RReadAdministrator")
	switch uid := self.Param("uid").MustInt64(); {
	//单独模式
	case uid > 0:
		{
			TplNames = "root/create_administrator"

			if usr, err := models.GetUser(uid); usr != nil && err == nil {
				self.Set("administrator", *usr)

			} else {
				self.Flash.Error(err.Error())
				return self.Render(TplNames)

			}
		}
	//列表模式
	case uid <= 0:
		{
			TplNames = "root/administrator_table"
			offset := self.Args("offset").MustInt64()
			limit := self.Args("limit").MustInt64()
			field := self.FormValue("field")

			if limit == 0 {
				limit = 1000 //默认限制显示最近1000条,需要显示全部请在提交请求的时候设置limit字段为-1
			}

			if field == "" {
				field = "id"
			}

			if usrs, err := models.GetUsersByRole(-1000, int(offset), int(limit), field); err == nil && usrs != nil {
				self.Set("administrators", *usrs)
			}
		}
	}

	return self.Render(TplNames)

}
