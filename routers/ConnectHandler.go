package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetConnectHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	if !(_usr_ != nil) {
		self.Flash.Warning("请先行登录再进行聊天操作！")
		return self.Redirect("/")
	}

	allow := false
	self.Set("contactHome", allow)
	self.Set("contactSearch", allow)
	self.Set("catpage", "ConnectHandler")

	receiver := self.Param("name").String()
	if len(receiver) == 0 {
		receiver = self.FormValue("name")
	}

	if len(receiver) > 0 {
		if receiver == _usr_.Username {
			allow = false
			self.Set("allow", allow)
			self.Flash.Warning("不允许自己与自己聊天！")
			return self.Redirect("/contact/")
		}

		recipient, err := models.GetUserByUsername(receiver)
		if (err != nil) || (recipient == nil) {
			allow = false
			self.Set("allow", allow)
			self.Flash.Error("聊天对象不存在此世间！")
			return self.Redirect("/contact/")
		}

		if !models.IsFriend(_usr_.Id, recipient.Id) {
			allow = false
			self.Set("allow", allow)
			self.Flash.Error("你们还不是好友，请先行建立好友关系～")
			return self.Redirect(fmt.Sprintf("/friend/add/%v/", recipient.Id))
		}

		key := _usr_.Username + ":" + helper.AesKey + ":" + receiver
		token := helper.EncryptHash(key, nil)
		cache.Store(self).Set(token, _usr_, 60)

		allow = true
		self.Set("token", token)
		self.Set("receiver", receiver)
		self.Set("recipient", *recipient)
	}

	self.Set("messager", true)
	self.Set("allow", allow)
	return self.Render("connect")

}
