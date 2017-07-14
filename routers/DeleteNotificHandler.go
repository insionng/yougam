package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"github.com/insionng/yougam/models"
)

func DeleteNotificGetHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	notificid := self.Param("notificid").MustInt64()
	if notificid > 0 {

		if models.DelNotificationByRole(notificid, _usr_.Id, -1000) == nil {
			_u, _e := models.GetUser(_usr_.Id)
			if (_u != nil) && (_e == nil) {
				cache.Store(self).Set(fmt.Sprintf("User:%v", _usr_.Id), _u, 60*60*24)
			}
			self.Flash.Success("删除通知讯息成功!")

		} else {
			self.Flash.Error("删除通知讯息失败")
		}

		return self.Redirect("/notifications/")
	} else {
		self.Flash.Error("你无权删除别人的通知讯息")
		return self.Redirect("/notifications/")
	}

}
