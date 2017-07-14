package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"github.com/insionng/yougam/models"
)

func GetRSettingPasswordHandler(self *makross.Context) error {
	self.Set("catpage", "RSettingPasswordHandler")
	TplNames := "root/setting_password"
	return self.Render(TplNames)
}

func PostRSettingPasswordHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	email := self.FormValue("email")
	oldpassword := self.FormValue("oldpassword")
	newpassword := self.FormValue("newpassword")
	username := _usr_.Username

	if len(username) == 0 {
		self.Flash.Error("请填写管理员用户名,不能为空!")
		return self.Redirect("/root/setting/password/")

	}

	if len(email) == 0 {
		self.Flash.Error("请填写管理员注册邮箱,以便系统验证,不能为空!")
		return self.Redirect("/root/setting/password/")

	}

	if len(oldpassword) == 0 {
		self.Flash.Error("请填写管理员当前密码,以便系统验证你的身份,不能为空!")
		return self.Redirect("/root/setting/password/")

	}

	if len(newpassword) == 0 {
		self.Flash.Error("请填写管理员新密码,不能为空!")
		return self.Redirect("/root/setting/password/")

	}

	if len(username) > 0 && len(email) > 0 && len(oldpassword) > 0 && len(newpassword) > 0 {

		if t, e := models.SetUserNewpassword(username, email, oldpassword, newpassword); t {
			self.Flash.Success(fmt.Sprint("设置新密码成功", e))
		} else {
			self.Flash.Error(fmt.Sprint("设置新密码失败", e))
		}

		return self.Redirect("/root/setting/password/")

	} else {
		self.Flash.Error("请完整填写所有项目,不能为空!")
		return self.Redirect("/root/setting/password/")

	}

}
