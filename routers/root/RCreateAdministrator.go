package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRCreateRootHandler(self *makross.Context) error {
	self.Set("catpage", "RCreateAdministrator")
	TplNames := "root/create_administrator"
	return self.Render(TplNames)

}

func PostRCreateRootHandler(self *makross.Context) error {

	
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
			return self.Redirect("/root/create/administrator/")
		} else if password != repassword {
			self.Flash.Error("密码前后不一致,请确认你输入的密码正确无误!")
			return self.Redirect("/root/create/administrator/")
		}
	} else {

		self.Flash.Error("密码为空!")
		return self.Redirect("/root/create/administrator/")
	}

	if len(username) <= 0 {
		self.Flash.Error("用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空!")
		return self.Redirect("/root/create/administrator/")
	}

	if len(email) > 0 {
		if helper.CheckEmail(email) == false {
			self.Flash.Error("邮箱格式错误!")
			return self.Redirect("/root/create/administrator/")
		}

	} else {
		if helper.CheckEmail(username) == true {
			email = username
		} else if helper.CheckUsername(username) == false {
			self.Flash.Error("用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空!")
			return self.Redirect("/root/create/administrator/")
		}

	}

	if len(email) > 0 {
		if usrinfo, err := models.GetUserByEmail(email); (usrinfo != nil) && (err == nil) {

			if usrinfo, err := models.GetUserByUsername(username); usrinfo != nil {
				self.Flash.Error("此用户名不能使用!")
				return self.Redirect("/root/create/administrator/")

			} else if err != nil {

				self.Flash.Error("检索用户名账号期间出错!")
				return self.Redirect("/root/create/administrator/")

			}

			self.Flash.Error("此Email不能使用!")
			return self.Redirect("/root/create/administrator/")

		} else if !((usrinfo == nil) && (err == nil)) {
			self.Flash.Error("检索EMAIL账号期间出错!")
			return self.Redirect("/root/create/administrator/")

		}
	} else {
		if usrinfo, err := models.GetUserByUsername(username); usrinfo != nil {

			self.Flash.Error("此用户名已经被注册,请重新命名!")
			return self.Redirect("/root/create/administrator/")

		} else if err != nil {

			self.Flash.Error("检索账号数据期间出错!")
			return self.Redirect("/root/create/administrator/")

		}
	}

	if usrid, err := models.AddUser(email, username, nickname, "", helper.EncryptHash(password, nil), group, content, mobile, gender, -1000); err != nil || usrid <= 0 {
		self.Flash.Error(fmt.Sprintf("用户注册信息写入数据库时发生错误:%v", err))
		return self.Redirect("/root/create/administrator/")

	} else {

		if usrinfo, err := models.GetUser(usrid); err == nil && usrinfo != nil {
			self.Flash.Success("添加用户账号保存成功!")
			return self.Redirect("/root/read/administrator/" + strconv.FormatInt(usrid, 10) + "/")

		} else {

			self.Flash.Error("获取用户数据出错!")
			return self.Redirect("/root/create/administrator/")

		}

	}

}
