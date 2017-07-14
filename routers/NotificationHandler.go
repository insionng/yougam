package routers

import (
	"github.com/insionng/makross"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func NotificGetHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		self.Flash.Error("尚未登录认证")
		return self.Redirect("/")
	}

	self.Set("catpage", "NotificationHandler")

	limit := self.Param("limit").MustInt64()
	if limit <= 0 {
		limit = self.Args("limit").MustInt64()
		if limit <= 0 {
			limit = 25
		}
	}

	offset := self.Param("offset").MustInt64()
	if offset <= 0 {
		offset = self.Args("offset").MustInt64()
	}

	ctype := self.Param("ctype").MustInt64()
	if ctype <= 0 {
		ctype = self.Args("ctype").MustInt64()
	}

	tab := self.Param("tab").String()
	if len(tab) <= 0 {
		tab = self.FormValue("tab")
	}

	page := self.Param("page").MustInt64()
	if page <= 0 {
		page = self.Args("page").MustInt64()
	}

	notificationsCount := models.GetNotificationsByUid(_usr_.Id, ctype, 0, 0, "id")
	rcs := int64(len(*notificationsCount))
	pages, pageout, beginnum, endnum, offset := helper.Pages(rcs, page, limit)
	notifications := models.GetNotificationsByUid(_usr_.Id, ctype, int(offset), int(limit), "id")

	curUser, _ := models.GetUser(_usr_.Id)
	if curUser != nil {
		if curUser.NotificationCount > 0 {
			models.SetNotificationCount(_usr_.Id, 0)
			curUser.NotificationCount = 0
			self.Session.Set("SignedUser", curUser)
		}
	}

	self.Set("notifications", notifications)
	self.Set("pagesbar", helper.Pagesbar("/notifications/", "", rcs, pages, pageout, beginnum, endnum, 5))
	return self.Render("notifications")

}
