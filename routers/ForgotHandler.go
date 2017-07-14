package routers

import (
	"encoding/base64"
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"strings"
	"time"
	"github.com/insionng/yougam/helper"
	simplejson "github.com/insionng/yougam/libraries/bitly/go-simplejson"
	"github.com/insionng/yougam/models"
)

func GetForgotHandler(self *makross.Context) error {
	

	self.Set("catpage", "ForgotHandler")
	tplname := "forgot"

	var keyin string
	if keyin = self.Param("key").String(); len(keyin) <= 0 {
		keyin = self.FormValue("key")
	}

	userkey := helper.MD5to16(self.RealIP())

	switch {
	case ((len(keyin) > 0) && (len(userkey) > 0)):
		{

			keybyte, e := base64.URLEncoding.DecodeString(keyin)
			if e != nil {
				self.Flash.Error("抱歉,取回密码认证链接非法!")
				goto render
			}
			keyin = string(keybyte)
			if s, err := helper.Aes128COMDecrypt(keyin, helper.AesConstKey); err != nil {

				self.Flash.Error("抱歉,取回密码认证链接不合法!!")
				goto render

			} else {
				if s == "" {
					self.Flash.Error("抱歉,取回密码认证失败!")
					goto render
				}
				if sjson, err := simplejson.NewJson([]byte(s)); err != nil {

					self.Flash.Error("抱歉,取回密码认证失败!!")
					goto render
				} else {

					data := sjson.Get("data").MustMap()
					if data == nil {
						self.Flash.Error("抱歉,取回密码认证失败!!!")

						goto render
					}

					var key, email, randkey, timekey string
					if data["key"] != nil {
						key = data["key"].(string)
					}
					if data["email"] != nil {
						email = data["email"].(string)
					}
					if data["randkey"] != nil {
						randkey = data["randkey"].(string)
					}
					if data["timekey"] != nil {
						timekey = data["timekey"].(string)
					}

					timekeyInt64, _ := strconv.ParseInt(timekey, 10, 0)
					key2 := helper.SHA1(randkey + userkey + email + helper.AesConstKey + timekey)
					if now := time.Now().Unix(); timekeyInt64 < now { //超时

						self.Flash.Error("抱歉,该认证链接已经超时失效!")
						goto render
					}
					if key == key2 {

						tplname = "password"
						self.Set("fmail", email)

						if usr, e := models.GetUserByEmail(email); usr != nil && e == nil {
							tplname = "password"
							self.Flash.Success("恭喜,该次链接校验成功,请设置密码!")
							goto render
						} else {

							self.Flash.Error("抱歉,该Email并没有用作注册账号!")
							goto render
						}
					} else {

						self.Flash.Error("抱歉,该认证链接验证失败!")
						goto render
					}

				}
			}
		}
	default:
		{
			tplname = "forgot"
			goto render
		}
	}
render:
	return self.Render(tplname)
}

func PostForgotHandler(self *makross.Context) error {
	

	self.Set("catpage", "ForgotHandler")
	tplname := "forgot"
	var keyin string
	if keyin = self.Param("key").String(); len(keyin) <= 0 {
		keyin = self.FormValue("key")
	}

	userkey := helper.MD5to16(self.RealIP())

	switch {
	case ((len(keyin) > 0) && (len(userkey) > 0)): //进入重置密码分支
		{

			keybyte, e := base64.URLEncoding.DecodeString(keyin)
			if e != nil {
				self.Flash.Error("抱歉,取回密码认证链接非法!")
				goto render
			}
			keyin = string(keybyte)
			if s, err := helper.Aes128COMDecrypt(keyin, helper.AesConstKey); err != nil {

				self.Flash.Error("抱歉,取回密码认证链接不合法!!")
				goto render

			} else {
				if len(s) == 0 {
					self.Flash.Error("抱歉,取回密码认证失败!")
					goto render
				}
				if sjson, err := simplejson.NewJson([]byte(s)); err != nil {

					self.Flash.Error("抱歉,取回密码认证失败!!")
					goto render
				} else {

					data := sjson.Get("data").MustMap()
					if data == nil {
						self.Flash.Error("抱歉,取回密码认证失败!!!")
						goto render
					}

					var key, email, randkey, timekey string
					if data["key"] != nil {
						key = data["key"].(string)
					}
					if data["email"] != nil {
						email = data["email"].(string)
					}
					if data["randkey"] != nil {
						randkey = data["randkey"].(string)
					}
					if data["timekey"] != nil {
						timekey = data["timekey"].(string)
					}

					timekeyInt64, _ := strconv.ParseInt(timekey, 10, 0)
					key2 := helper.SHA1(randkey + userkey + email + helper.AesConstKey + timekey)
					if now := time.Now().Unix(); timekeyInt64 < now { //超时
						self.Flash.Error("抱歉,该认证链接已经超时失效!")
						goto render
					}
					if key == key2 {

						tplname = "password"
						self.Set("fmail", email)

						if usr, e := models.GetUserByEmail(email); usr != nil && e == nil {
							tplname = "password"
							password := self.FormValue("password")
							repassword := self.FormValue("repassword")
							if len(password) == 0 {
								self.Flash.Error("密码不能设置为空!")
								goto render
							}

							if password == repassword {
								usr.Password = helper.EncryptHash(password, nil)

								if rows, e := models.PutUser(usr.Id, usr); rows <= 0 || e != nil {

									self.Flash.Error("抱歉,服务器写入数据发送错误,请你稍后重新尝试.")
									goto render
								} else {

									self.Flash.Success("恭喜,你已经成功设置新密码!")
									goto render
								}
							} else {
								self.Flash.Error("抱歉,新设置密码与确认密码不一致,请你修改后重新尝试.")
								goto render
							}
						} else {

							self.Flash.Error("抱歉,该Email并没有用作注册账号!")
							goto render
						}
					} else {

						self.Flash.Error("抱歉,该认证链接验证失败!")
						goto render
					}

				}
			}
		}
	default: //发送认证链接邮件分支
		{
			var email string
			if email = strings.TrimSpace(self.Param("email").String()); len(email) <= 0 {
				email = strings.TrimSpace(self.FormValue("email"))
			}

			if len(email) <= 0 {
				self.Flash.Error("抱歉,提交的邮件地址不能为空!")
				goto render
			}

			randkey := helper.GUID32BIT()
			d, _ := time.ParseDuration("1h") //设为1小时后失效
			timekey := strconv.FormatInt(time.Now().Add(d).UnixNano(), 10)
			key := helper.SHA1(randkey + userkey + email + helper.AesConstKey + timekey)

			m := map[string]string{}
			m["key"] = key
			m["email"] = email
			m["randkey"] = randkey //string
			m["timekey"] = timekey //string

			crypted, _ := helper.SetJsonCOMEncrypt(1, "", m)
			b64key := base64.URLEncoding.EncodeToString([]byte(crypted))

			if helper.CheckEmail(email) {

				if usr, e := models.GetUserByEmail(email); usr != nil && e == nil {

					url := fmt.Sprintf("http://%s/forgot/?key=%s", helper.Domain, b64key)
					body := "<html><body><div>"
					body = body + "<p>亲爱的 <strong>" + usr.Username + "</strong> ，您好：</p>"
					body = body + `<p style="padding-left:2em;">您申请了重置` + helper.SiteName + `登录密码，请点击以下链接设置新密码.</p>`
					body = body + `<p><a href="` + url + `" target="_blank">` + url + `</a></p>`
					body = body + "<p>(如果无法点击该URL链接地址，请将它复制并粘帖到浏览器的地址输入框访问即可.)</p>"
					body = body + "<p>注意:请您在收到邮件后点击此链接以完成重置，否则该链接将会在1小时后失效.</p>"
					body = body + "<p>我们将一如既往竭诚为您服务！</p>"
					body = body + `<p style="border-bottom:dotted 1px gray;height:1px;"></p>`
					body = body + "<p>" + helper.MailAdline + "</p>"
					body = body + "</div></body></html>"

					subject := "【" + helper.SiteName + "】找回密码"

					//发送邮件
					if err := helper.SendEmail(email, subject, body, "html"); err != nil {
						self.Flash.Error("抱歉," + helper.SiteName + "发送重置邮件失败,请你稍后重新尝试.")
						goto render
					} else {
						self.Flash.Success("恭喜," + helper.SiteName + "已经发送重置邮件到你的注册邮箱,请你注意查收.")
						goto render
					}

				} else {
					self.Flash.Error("抱歉,该Email并没有用作注册账号!")
					return self.Redirect("/forgot/")
				}
			} else {
				self.Flash.Error("抱歉,你提交的邮件地址不合法!")
				goto render
			}
		}
	}
render:
	return self.Render(tplname)
}
