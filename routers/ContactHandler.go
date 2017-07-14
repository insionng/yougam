package routers

import (
	"github.com/insionng/makross"
	"github.com/insionng/yougam/models"
)

func GetContactHandler(self *makross.Context) error {

	self.Set("catpage", "ContactHandler")
	self.Set("messager", true)

	offset := self.Param("offset").MustInt()
	if offset <= 0 {
		offset = self.Args("offset").MustInt()
	}

	limit := self.Param("limit").MustInt()
	if offset <= 0 {
		limit = self.Args("limit").MustInt()
	}

	limit = 16 //临时设置，以后再增加翻页处理

	allow := false
	self.Set("contactHome", false)
	self.Set("contactSearch", false)

	if self.RequestURI() == "/contact/" {
		allow = false
		self.Set("contactHome", true)
	}

	if self.RequestURI() == "/contact/search/" {
		allow = false
		self.Set("contactSearch", true)
	}

	if (self.RequestURI() == "/contact/") || (self.RequestURI() == "/contact/search/") {
		u, e := models.GetUsersOnHotness(offset, limit, "created")
		if (e == nil) && (u != nil) {
			self.Set("UsersByCreated", u)
		}

		u_, e_ := models.GetUsersOnHotness(offset, limit, "confidence")
		if (e_ == nil) && (u_ != nil) {
			self.Set("UsersByConfidence", u_)
		}
	}

	self.Set("allow", allow)
	return self.Render("contact")

}
