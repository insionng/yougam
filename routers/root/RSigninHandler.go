package root

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"
	"github.com/insionng/makross/captcha"

	"time"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRSigninHandler(self *makross.Context) error {

	if _, okay := self.Session.Get("SignedUser").(*models.User); okay {
		//如果已登录root
		return self.Redirect("/root/dashboard/")
	} else {
		return self.Render("root/signin")
	}

}

func PostRSigninHandler(self *makross.Context) error {

	username := self.FormValue("username")
	password := self.FormValue("password")

	if cpt := captcha.Store(self); cpt.VerifyReq(self) {
		if len(username) > 0 && len(password) > 0 {

			if usr, err := models.GetUserByUsername(username); usr != nil && err == nil {
				if helper.ValidateHash(usr.Password, password) && usr.Role == -1000 {

					//登录成功设置self.Session.on
					self.Session.Set("SignedUser", usr)
					cache.Store(self).Set(fmt.Sprintf("User:%v", usr.Id), usr, 60*60*60)
					models.PutSignin2User(usr.Id, time.Now().Unix(), usr.SigninCount+1, self.RealIP())
					self.Flash.Success("登录成功！", false)
					return self.Redirect("/root/dashboard/")

				} else {
					self.Flash.Error("密码错误！", false)
					return self.Redirect("/root/signin/")

				}
			} else {
				self.Flash.Error("不存在此用户！", false)
				return self.Redirect("/root/signin/")

			}
		} else {
			if len(username) > 0 {
				self.Flash.Error("用户名不能为空！", false)
			} else {
				self.Flash.Error("密码不能为空！", false)
			}
			return self.Redirect("/root/signin/")

		}
	} else {
		self.Flash.Error("验证码错误~", false)
		return self.Redirect("/root/signin/")

	}

}
