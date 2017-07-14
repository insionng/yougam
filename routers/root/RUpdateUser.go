package root

import (
	"errors"
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRUpdateUserHandler(self *makross.Context) error {
	

	self.Set("catpage", "RUpdateUserHandler")
	TplNames := "root/update_user"

	if userid := self.Param("uid").MustInt64(); userid > 0 {
		if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {
			self.Set("usr", usrinfo)

			// 匹配处理： /root/update/user/:uid/password/
			if holder := self.Param("holder").String(); len(holder) > 0 && holder == "password" {
				randomPassword := helper.StringNewRand(8)
				usrinfo.Password = helper.EncryptHash(randomPassword, nil)
				if row, err := models.PutUser(userid, usrinfo); err != nil || row <= 0 {
					self.Flash.Error("用户信息更新到数据库时发生错误!")
					return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

				} else {
					body := "<html><body><div>"
					body = body + "<p>亲爱的 <strong>" + usrinfo.Username + "</strong> ，您好：</p>"
					body = body + `<p style="padding-left:2em;">` + helper.SiteName + `的管理员已为你重置了` + helper.SiteName + `的登录密码为：</p>`
					body = body + `<p>[<b>` + randomPassword + `</b>]</p>`
					body = body + "<p>我们将一如既往竭诚为您服务！</p>"
					body = body + `<p style="border-bottom:dotted 1px gray;height:1px;"></p>`
					body = body + "<p>" + helper.MailAdline + "</p>"
					body = body + "</div></body></html>"

					subject := "【" + helper.SiteName + "】重置密码"

					//发送邮件
					if err := helper.SendEmail(usrinfo.Email, subject, body, "html"); err != nil {
						es := fmt.Sprintf("为用户【%s】发送随机密码到【%s】失败,请你稍后重新尝试.", usrinfo.Username, usrinfo.Email)
						self.Flash.Error(es)
						return errors.New(es)
					} else {
						self.Flash.Success(fmt.Sprintf("为用户【%s】生成随机密码并已成功发送到【%s】.", usrinfo.Username, usrinfo.Email))
						return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")
					}

				}

			}
		}
		return self.Render(TplNames)

	} else {
		return self.Redirect("/root/read/user/")
	}

}

func PostRUpdateUserHandler(self *makross.Context) error {
	

	self.Set("catpage", "RUpdateUserHandler")

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
				return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

			} else if password != repassword {

				self.Flash.Error("密码前后不一致,请确认你输入的密码正确无误!")
				return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

			}
		} else {

			self.Flash.Error("密码为空!")
			return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

		}

		if len(username) == 0 {
			self.Flash.Error("用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空!")
			return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

		}

		if len(email) > 0 {
			if helper.CheckEmail(email) == false {
				self.Flash.Error("邮箱格式错误!")
				return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

			}

		} else {
			if helper.CheckEmail(username) == true {
				email = username
			} else if helper.CheckUsername(username) == false {
				self.Flash.Error("用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空!")
				return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

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
		//usr.Role = 1

		if row, err := models.PutUser(userid, usr); err != nil || row <= 0 {
			self.Flash.Error("用户注册信息写入数据库时发生错误!")
			return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

		} else {

			if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {
				self.Set("usr", usrinfo)
				self.Flash.Success("用户账号保存成功!")
				return self.Redirect("/root/update/user/" + strconv.FormatInt(userid, 10) + "/")

			} else {

				self.Flash.Error("获取用户数据出错!")
				return self.Redirect("/root/read/user/")

			}

		}
	} else {

		self.Flash.Error("用户不存在!")
		return self.Redirect("/root/read/user/")
	}
}
