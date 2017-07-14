package routers

import (
	"github.com/insionng/yougam/models"

	"github.com/insionng/makross"
)

func GetPageHandler(self *makross.Context) error {

	TplNames := "page"
	self.Set("navpage", "page")
	pageid := self.Param("pageid").MustInt64()
	page := self.Param("page").String()
	if len(page) > 0 {
		_p, e := models.GetPageByTitle(page)
		if !(e != nil) {
			self.Set("page", _p)
		}
		self.Set("curpage", page)

	} else if pageid != 0 {
		p, e := models.GetPage(pageid)
		if e == nil && p != nil {
			self.Set("curpage", p.Title)
		}

		self.Set("page", p)

	} else {

		return self.Redirect("/")

	}

	return self.Render(TplNames)

}
