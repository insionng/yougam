package routers

import (
	"encoding/base64"
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"
	"github.com/insionng/makross/captcha"

	"strconv"
	"time"

	"github.com/insionng/yougam/helper"
	simplejson "github.com/insionng/yougam/libraries/bitly/go-simplejson"
	"github.com/insionng/yougam/models"
)

var (
	IsSendmail = false
)

func init() {
	IsSendmail = helper.IsSendMail()
}

func GetSignupHandler(self *makross.Context) error {

	var _usr_ = new(models.User)
	if sUser, okay := self.Session.Get("SignedUser").(*models.User); okay {
		_usr_ = sUser
	}
	cc := cache.Store(self)

	self.Set("catpage", "SignupHandler")
	TplNames := "signup"
	var keyin string
	if keyin = self.Args("key").String(); len(keyin) <= 0 {
		keyin = self.QueryParam("key")
	}

	if len(keyin) > 0 {

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

				var key, email, username, password, randkey, timekey string
				if data["key"] != nil {
					key = data["key"].(string)
				}
				if data["email"] != nil {
					email = data["email"].(string)
				}
				if data["username"] != nil {
					username = data["username"].(string)
				}
				if data["password"] != nil {
					password = data["password"].(string)
				}
				if data["randkey"] != nil {
					randkey = data["randkey"].(string)
				}
				if data["timekey"] != nil {
					timekey = data["timekey"].(string)
				}

				timekeyInt64, _ := strconv.ParseInt(timekey, 10, 0)
				key2 := helper.SHA1(randkey + username + helper.AesConstKey + email + timekey + password)

				if now := time.Now().Unix(); timekeyInt64 < now { //超时
					self.Flash.Error("抱歉,该邮件认证链接已经超时失效!")
					goto render
				}

				if key == key2 {

					if helper.CheckEmail(email) {
						if usr, e := models.GetUserByEmail(email); usr == nil || e != nil {

							if usrinfo, _ := models.GetUserByUsername(username); usrinfo != nil {

								self.Flash.Error("用户名已被注册~")
								goto render

							} else {

								if usrid, err := models.AddUser(email, username, "", "", helper.EncryptHash(password, nil), "", "", "", 0, 1); err != nil {
									self.Flash.Error("用户注册信息写入数据库时发生错误~")
									goto render
								} else {

									if usrinfo, err := models.GetUser(usrid); err == nil || usrinfo != nil {

										//注册账号成功,以下自动登录并设置Session
										_usr_ = usrinfo
										self.Set("IsSigned", true)
										self.Set("SignedUser", _usr_)
										self.Session.Set("SignedUser", _usr_)
										cc.Set(fmt.Sprintf("SignedUser:%v", _usr_.Id), _usr_, 60*60*24)
										models.PutSignin2User(_usr_.Id, time.Now().Unix(), _usr_.SigninCount+1, self.RealIP())

										if e := models.SetAmountByUid(usrinfo.Id, 2, 2000, "注册收益2000金币"); e == nil {
											if usrinfo, err := models.GetUser(usrid); err == nil || usrinfo != nil {
												cc.Set(fmt.Sprintf("SignedUser:%v", usrinfo.Id), usrinfo, 60*60*24)
											}

										}

										self.Flash.Success("账号登录成功~", false)
										//self.Session.写入后直接跳到首页
										return self.Redirect("/")

									} else {

										self.Flash.Success("注册账号成功,但登录失败,请手动登录~", false)
										//注册成功后直接跳转到登录页
										return self.Redirect("/signin/")

									}
								}
							}
						} else {
							self.Flash.Error("抱歉,该Email已经用作注册账号!", false)
							return self.Redirect("/signup/")
						}
					} else {
						self.Flash.Error("抱歉,你提交的邮件地址不合法!")
						goto render
					}

				} else {

					self.Flash.Error("抱歉,该邮件认证链接验证失败!")
					goto render
				}

			}
		}
	}
	//侧栏九宫格推荐榜单
	//先行取出最热门的9个节点 然后根据节点获取该节点下最热门的话题
	/*
		if nd, err := models.GetNodes(0, 9, "hotness"); err == nil {
			if len(*nd) > 0 {
				for _, v := range *nd {

					i := 0
					output_start := `<ul class="widgets-popular widgets-similar clx">`
					output := ""
					if tps := models.GetTopicsByNid(v.Id, 0, 1, 0, "hotness"); err == nil {

						if len(*tps) > 0 {
							for _, v := range *tps {

								i += 1
								if i == 3 {
									output = output + `<li class="similar similar-third">`
									i = 0
								} else {
									output = output + `<li class="similar">`
								}
								output = output + `<a target="_blank" href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `">
													<img src="` + v.ThumbnailsSmall + `" wdith="70" height="70" />
												</a>
											</li>`
							}
						}
					}
					output_end := "</ul>"
					if len(output) > 0 {
						output = output_start + output + output_end
						self.Set["topic_hotness_9"] = template.HTML(output)
					} else {
						self.Set["topic_hotness_9"] = nil
					}

				}
			}
		} else {
			fmt.Println("节点数据查询出错", err)
		}
	*/
render:
	return self.Render(TplNames)

}

func PostSignupHandler(self *makross.Context) error {

	var _usr_ = new(models.User)
	if sUser, okay := self.Session.Get("SignedUser").(*models.User); okay {
		_usr_ = sUser
	}
	cc := cache.Store(self)

	self.Set("catpage", "SignupHandler")
	TplNames := "signup"

	email := self.FormValue("email")
	username := self.FormValue("username")
	password := self.FormValue("password")
	repassword := self.FormValue("repassword")

	self.Set("tmpemail", email)
	self.Set("tmpusername", username)
	self.Set("tmppassword", password)
	self.Set("tmprepassword", repassword)

	if cpt := captcha.Store(self); cpt.VerifyReq(self) {

		if password == "" {
			self.Flash.Error("密码为空~")
			goto render
		}

		if password != repassword {
			self.Flash.Error("两次密码不匹配~")
			goto render
		}

		if helper.CheckPassword(password) == false {
			self.Flash.Error("密码含有非法字符或密码过短(至少4~30位密码)!")
			goto render
		}

		if username == "" {
			self.Flash.Error("用户名是为永久性设定,不能少于2个字或多于30个字,请慎重考虑,不能为空~")
			goto render
		}

		if helper.CheckUsername(username) == false {
			self.Flash.Error("用户名是为永久性设定,不能少于2个字或多于30个字,请慎重考虑,不能为空~")
			goto render
		}

		if helper.CheckEmail(email) == false {
			self.Flash.Error("Email格式不合符规格~")
			goto render
		}

		if usrinfo, _ := models.GetUserByEmail(email); usrinfo != nil {
			self.Flash.Error("抱歉,该Email已经用作注册账号!")
			goto render

		} else {

			if usrinfo, _ := models.GetUserByUsername(username); usrinfo != nil {
				self.Flash.Error("用户名已被注册~")
				goto render

			}
		}

		//开启邮件认证注册
		if !IsSendmail {
			if usrid, err := models.AddUser(email, username, "", "", helper.EncryptHash(password, nil), "", "", "", 0, 1); err != nil {
				self.Flash.Error("用户注册信息写入数据库时发生错误~")
				goto render
			} else {
				if e := models.SetAmountByUid(usrid, 2, 2000, "注册收益2000金币"); e == nil {
					if usrinfo, err := models.GetUser(usrid); err == nil || usrinfo != nil {

						//注册账号成功,以下自动登录并设置self.Session.on
						_usr_ = usrinfo
						self.Set("IsSigned", true)
						self.Set("SignedUser", _usr_)
						self.Session.Set("SignedUser", _usr_)
						cc.Set(fmt.Sprintf("SignedUser:%v", usrinfo.Id), usrinfo, 60*60*24)
						models.PutSignin2User(_usr_.Id, time.Now().Unix(), _usr_.SigninCount+1, self.RealIP())

						self.Flash.Success("账号登录成功~", false)

						//self.Session.写入后直接跳到首页
						return self.Redirect("/")

					} else {

						self.Flash.Success("注册账号成功,但登录失败,请手动登录~", false)
						//注册成功后直接跳转到登录页
						return self.Redirect("/signin/")
					}

				}

			}
		} else {

			randkey := helper.GUID32BIT()
			d, _ := time.ParseDuration("1h") //设为1小时后失效
			timekey := strconv.FormatInt(time.Now().Add(d).UnixNano(), 10)
			key := helper.SHA1(randkey + username + helper.AesConstKey + email + timekey + password)

			m := map[string]string{}
			m["key"] = key
			m["email"] = email
			m["username"] = username
			m["password"] = password
			m["randkey"] = randkey //string
			m["timekey"] = timekey //string

			crypted, _ := helper.SetJsonCOMEncrypt(1, "", m)
			b64key := base64.URLEncoding.EncodeToString([]byte(crypted))
			if len(b64key) > 0 {
				if usr, e := models.GetUserByEmail(email); usr == nil || e != nil {

					url := fmt.Sprintf("http://%s/signup/?key=%s", helper.Domain, b64key)
					body := "<html><body><div>"
					body = body + "<p>亲爱的 <strong>" + username + "</strong> ，您好：</p>"
					body = body + `<p style="padding-left:2em;">您申请了注册` + helper.SiteName + `账号，请点击以下认证链接.</p>`
					body = body + `<p><a href="` + url + `" target="_blank">` + url + `</a></p>`
					body = body + "<p>(如果无法点击该URL链接地址，请将它复制并粘帖到浏览器的地址输入框访问即可.)</p>"
					body = body + "<p>注意:请您在收到邮件后点击此链接以完成注册，否则该链接将会在1小时后失效.</p>"
					body = body + "<p>我们将一如既往竭诚为您服务！</p>"
					body = body + `<p style="border-bottom:dotted 1px gray;height:1px;"></p>`
					body = body + "<p>" + helper.MailAdline + "</p>"
					body = body + "</div></body></html>"

					subject := "【" + helper.SiteName + "】注册账号"

					//发送邮件
					if err := helper.SendEmail(email, subject, body, "html"); err != nil {
						self.Flash.Error("抱歉," + helper.SiteName + "发送注册认证邮件失败,请你稍后重新尝试." + fmt.Sprintf("错误：[%v]", err))
						goto render
					} else {
						self.Flash.Success("恭喜,"+helper.SiteName+"已经发送注册认证邮件到你的注册邮箱,请你注意查收.", false)
						//邮件发送成功后跳转到登录页
						return self.Redirect("/signin/")

					}

				} else {
					self.Flash.Error("抱歉,该Email已经用作注册账号!")
					goto render
				}
			} else {
				self.Flash.Error("抱歉,服务器生成密匙出错，麻烦向管理员报告此问题!")
				goto render
			}

		}
	} else {

		if len(self.Args(cpt.FieldCaptchaName).String()) > 0 {
			self.Flash.Error("验证码错误~")
		} else {
			self.Flash.Error("验证码为空~")
		}
		goto render

	}

render:
	return self.Render(TplNames)
}
