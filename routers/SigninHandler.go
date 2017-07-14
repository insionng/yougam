package routers

import (
	"fmt"
	"strings"
	"time"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"
	"github.com/insionng/makross/captcha"
)

func GetSigninHandler(self *makross.Context) error {

	var IsSignin bool
	if _, okay := self.Session.Get("SignedUser").(*models.User); okay {
		IsSignin = okay
	}

	TplNames := "signin"
	self.Set("catpage", "SigninHandler")
	self.Set("IsCaptcha", helper.IsCaptcha)

	remember, _ := self.GetCookie("remember")
	if IsSignin { //如果已登录
		if next := self.Args("next").String(); len(next) > 0 {
			return self.Redirect(next)
		}
		return self.Redirect("/")
	} else { //如果未登录
		if remember != nil {
			if remember.Value == "true" {
				self.Set("remember", "true")
			} else {
				self.Set("remember", nil)
			}
		}

	}
	return self.Render(TplNames)
}

func PostSigninHandler(self *makross.Context) error {

	TplNames := "signin"
	cpt := new(captcha.Captcha)
	allow := false
	if helper.IsCaptcha {
		cpt = captcha.Store(self)
		allow = cpt.VerifyReq(self)
	}
	if helper.IsCaptcha && (!allow) {
		if len(self.Args(cpt.FieldCaptchaName).String()) > 0 {
			self.Flash.Error("验证码错误~")
		} else {
			self.Flash.Error("验证码为空~")
		}
		return self.Render(TplNames)
	}

	cc := cache.Store(self)

	//Secret := helper.MD5(self.Req.UserAgent() + helper.AesConstKey)
	self.Set("catpage", "SigninHandler")

	password := self.Args("password").String()
	self.Set("tmppassword", password)
	self.Set("tmpemail", self.Args("email").String())
	remember := self.Args("remember").String()

	if len(password) == 0 {
		self.Flash.Error("密码为空~")
		return self.Render(TplNames)
	}

	if helper.CheckPassword(password) == false {
		self.Flash.Error("密码含有非法字符或密码过短(至少4~30位密码)!")
		return self.Render(TplNames)
	}

	var err error
	var usr = new(models.User)
	var email, username string
	if isEmail := strings.Contains(self.Args("email").String(), "@"); isEmail {
		email = self.Args("email").String()
		if len(email) == 0 {
			self.Flash.Error("EMAIL为空~")
			goto render
		}

		if helper.CheckEmail(email) == false {
			self.Flash.Error("Email格式不合符规格~")
			goto render
		}

		usr, err = models.GetUserByEmail(email)
	} else {
		username = self.Args("email").String()
		if len(username) == 0 {
			self.Flash.Error("用户名称为空~")
			goto render
		}

		if helper.CheckUsername(username) == false {
			self.Flash.Error("用户名称格式不合符规格~")
			goto render
		}

		usr, err = models.GetUserByUsername(username)
	}

	if (usr != nil) && (err == nil) {

		if helper.ValidateHash(usr.Password, password) {

			//登录成功设置session
			self.Session.Set("SignedUserID", usr.Id)
			self.Session.Set("SignedUserName", usr.Username)
			self.Session.Set("SignedUser", usr)

			self.Set("IsSigned", true)
			self.Set("IsRoot", (usr.Role == -1000))
			self.Set("SignedUser", usr)
			self.Set("SignedUserID", usr.Id)
			self.Set("SignedUserName", usr.Username)
			cc.Set(fmt.Sprintf("SignedUser:%v", usr.Id), usr, 60*60*24)
			models.PutSignin2User(usr.Id, time.Now().Unix(), usr.SigninCount+1, self.RealIP())

			//设置cookie
			cookie := self.NewCookie()
			cookie.Name = "remember"
			if remember == "true" {
				cookie.Value = "true"
				cookie.Expires = (time.Now().Add(time.Duration(31190400))) //361 days
				//使用flower作本地存储时的Email别名
				//self.SetSuperSecureCookie(Secret, "flower", usr.Email, 31190400)
			} else {
				cookie.Value = ("false") //取消记录
				cookie.Expires = (time.Now().Add(time.Duration(-1)))
				//self.SetSuperSecureCookie(Secret, "flower", "", 3600) //删除数据
			}
			self.SetCookie(cookie)

			if next := self.Args("next").String(); next != "" {
				return self.Redirect(next)
			}
			return self.Redirect("/")

		} else {
			self.Flash.Error("密码无法通过校验~")
			goto render
		}
	} else {
		self.Flash.Error("该账号不存在~")
		goto render
	}
render:
	return self.Render(TplNames)
}
