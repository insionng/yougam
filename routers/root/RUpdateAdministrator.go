package root

import (
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRUpdateRootHandler(self *makross.Context) error {
	self.Set("catpage", "RUpdateAdministratorHandler")
	TplNames := "root/update_user"

	if userid := self.Param("uid").MustInt64(); userid > 0 {
		if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {
			self.Set("usr", usrinfo)
		}
		return self.Render(TplNames)

	} else {
		return self.Redirect("/root/read/administrator/")
	}

}

func PostRUpdateRootHandler(self *makross.Context) error {
	

	self.Set("catpage", "RUpdateAdministratorHandler")

	if userid := self.Param("uid").MustInt64(); userid > 0 {
		username := self.FormValue("username")
		nickname := self.FormValue("nickname")
		password := self.FormValue("password")
		repassword := self.FormValue("repassword")
		content := self.FormValue("content")
		mobile := self.FormValue("mobile")
		email := self.FormValue("email")
		group := self.FormValue("group")
		gender := self.Args("gender").MustInt64()

		if len(password) > 0 {
			if helper.CheckPassword(password) == false {
				self.Flash.Error("密码含有非法字符或密码过短(至少4~30位密码)!")
				return self.Redirect("/root/update/administrator/" + strconv.FormatInt(userid, 10) + "/")

			} else if password != repassword {

				self.Flash.Error("密码前后不一致,请确认你输入的密码正确无误!")
				return self.Redirect("/root/update/administrator/" + strconv.FormatInt(userid, 10) + "/")

			}
		} else {

			self.Flash.Error("密码为空!")
			return self.Redirect("/root/update/administrator/" + strconv.FormatInt(userid, 10) + "/")

		}

		if username == "" {
			self.Flash.Error("用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空!")
			return self.Redirect("/root/update/administrator/" + strconv.FormatInt(userid, 10) + "/")

		}

		if email != "" {
			if helper.CheckEmail(email) == false {
				self.Flash.Error("邮箱格式错误!")
				return self.Redirect("/root/update/administrator/" + strconv.FormatInt(userid, 10) + "/")

			}

		} else {
			if helper.CheckEmail(username) == true {
				email = username
			} else if helper.CheckUsername(username) == false {
				self.Flash.Error("用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空!")
				return self.Redirect("/root/update/administrator/" + strconv.FormatInt(userid, 10) + "/")

			}

		}

		usr, _ := models.GetUser(userid)
		//usr.Id = userid
		usr.Email = email
		usr.Username = username
		usr.Nickname = nickname
		usr.Password = helper.EncryptHash(password, nil)
		usr.Group = group
		usr.Content = content
		usr.Mobile = mobile
		usr.Gender = gender
		usr.Role = -1000

		if row, err := models.PutUser(userid, usr); err != nil && row <= 0 {

			self.Flash.Error("用户注册信息写入数据库时发生错误!")
			return self.Redirect("/root/update/administrator/" + strconv.FormatInt(userid, 10) + "/")

		} else {

			if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {
				self.Set("usr", usrinfo)
				self.Flash.Success("用户账号保存成功!")
				return self.Redirect("/root/update/administrator/" + strconv.FormatInt(userid, 10) + "/")

			} else {

				self.Flash.Error("获取用户数据出错!")
				return self.Redirect("/root/read/administrator/")

			}

		}
	} else {

		self.Flash.Error("用户不存在!")
		return self.Redirect("/root/read/administrator/")
	}
}
