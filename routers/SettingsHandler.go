package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func SettingsGetHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	TplNames := ""
	name := self.Param("name").String()
	self.Set("catpage", "SettingsHandler")

	if usr, err := models.GetUserByUsername(_usr_.Username); usr != nil && err == nil {
		self.Set("user", *usr)
	}

	switch {
	case name == "profile":
		{
			TplNames = "user/profile"

		}
	case name == "avatar":
		{
			TplNames = "user/avatar"

		}
	case name == "avatars":
		{
			TplNames = "user/avatar4qiniu"

		}
	case name == "password":
		{
			TplNames = "user/password"

		}
	default:
		{
			return self.Redirect("/user/" + name + "/")
		}

	}

	return self.Render(TplNames)
}

func SettingsPostHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	cc := cache.Store(self)

	self.Set("catpage", "SettingsHandler")
	name := self.Param("name").String()

	switch {
	case name == "profile":
		{

			username := _usr_.Username
			email := self.FormValue("email")

			nickname := self.FormValue("nickname")
			realname := self.FormValue("realname")

			content := self.FormValue("content")
			age := self.Args("age").MustInt64()
			school := self.FormValue("school")
			weight := self.Args("weight").MustInt64()
			height := self.Args("height").MustInt64()
			zodiacSign := self.FormValue("zodiacsign")
			occupation := self.FormValue("occupation")

			province := self.FormValue("province")
			city := self.FormValue("city")
			company := self.FormValue("company")
			address := self.FormValue("address")

			postcode := self.FormValue("postcode")
			mobile := self.FormValue("mobile")
			website := self.FormValue("website")
			gender := self.Args("gender").MustInt64()
			qq := self.FormValue("qq")
			weixin := self.FormValue("weixin")
			weibo := self.FormValue("weibo")
			/*
				if username == "" {
					self.Flash.Error("用户名不能为空!", false)
					return self.Redirect(fmt.Sprintf("/settings/%s/", name), 302)
					return
				}
			*/
			if len(content) == 0 {
				self.Flash.Error("为了让别人更了解你,请务必填写你的个人签名~")
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))
			}

			if helper.CheckUsername(username) == false {
				self.Flash.Error("用户名包含非法字符,或不合符规格(限4~30个字符)~")
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))
			}

			if len(email) == 0 {
				self.Flash.Error("Email是你的主账号,和主要联系方式,不能留空~")
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			} else {

				if helper.CheckEmail(email) == false {
					self.Flash.Error("Email格式不合符规格~")
					return self.Redirect(fmt.Sprintf("/settings/%s/", name))

				} else {
					cansave := bool(false)
					if usremailinfo, err := models.GetUserByEmail(email); usremailinfo == nil || err != nil {
						//email不存在  可用
						cansave = true
					} else {
						//如果Email存在 则继续判断该Email是否当前用户正在使用的Email
						if _usr_.Email == email { //可以继续使用
							cansave = true
						} else { //已经被其他用户所用 不能使用
							cansave = false
						}

					}

					if cansave {
						if usrinfo, err := models.GetUser(_usr_.Id); usrinfo != nil && err == nil {

							usrinfo.Username = username
							usrinfo.Email = email
							usrinfo.Nickname = nickname
							usrinfo.Realname = realname
							usrinfo.Content = content
							usrinfo.Birth = time.Now().Unix()
							usrinfo.Province = province
							usrinfo.City = city
							usrinfo.Company = company
							usrinfo.Address = address
							usrinfo.Postcode = postcode
							usrinfo.Mobile = mobile
							usrinfo.Website = website
							usrinfo.Gender = gender
							usrinfo.Qq = qq
							usrinfo.Weixin = weixin
							usrinfo.Weibo = weibo
							usrinfo.Occupation = occupation
							usrinfo.Age = age
							usrinfo.Height = height
							usrinfo.School = school
							usrinfo.Weight = weight
							usrinfo.ZodiacSign = zodiacSign

							if _, err := models.PutUser(usrinfo.Id, usrinfo); err == nil {

								//更新self.Session.on

								self.Session.Set("SignedUser", usrinfo)
								_usr_ = usrinfo
								self.Set("SignedUser", _usr_)

								cc.Set(fmt.Sprintf("SignedUser:%v", _usr_.Id), _usr_, 60*60*24)
								self.Flash.Success("设置个人信息成功~")
							} else {
								self.Flash.Error("设置个人信息失败~")
							}

							return self.Redirect(fmt.Sprintf("/settings/%s/", name))

						} else {

							self.Flash.Error("该账号不存在~")
							return self.Redirect(fmt.Sprintf("/settings/%s/", name))

						}
					} else {

						self.Flash.Error("设置Email失败,该Email已被注册使用~")
						return self.Redirect(fmt.Sprintf("/settings/%s/", name))

					}

				}
			}

		}
	case name == "avatar":
		{

			uid := _usr_.Id
			targetFolder := helper.FileStorageDir + "file"
			mf, err := self.FormFile("avatar")
			if err != nil {
				self.Flash.Error(fmt.Sprint("SettingsHandler获取文件错误1,", err.Error()))
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			}

			file, e := mf.Open()
			if e != nil {
				self.Flash.Error(fmt.Sprint("SettingsHandler获取文件错误2,", e.Error()))
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			}

			defer file.Close()

			if file != nil {

				ext := strings.ToLower(path.Ext(mf.Filename))
				filename := helper.MD5(fmt.Sprint(uid)) + ext

				dirpath := fmt.Sprintf("%v/%v", targetFolder, fmt.Sprintf("%v/", uid))

				_err := os.MkdirAll(dirpath, 0755)
				if _err != nil {
					self.Flash.Error(_err.Error())
					return self.Redirect(fmt.Sprintf("/settings/%s/", name))

				}

				filepath := fmt.Sprintf("%s%s", dirpath, filename)
				//f, err := os.OpenFile(fmt.Sprintf(".%s", filepath), os.O_WRONLY|os.O_CREATE, 0755)
				f, err := os.Create(filepath)

				if err != nil {
					self.Flash.Error("SettingsHandler获取文件错误2!")
					return self.Redirect(fmt.Sprintf("/settings/%s/", name))

				}

				defer f.Close()
				io.Copy(f, file)

				inputFile := filepath
				output_file := fmt.Sprintf("%s%v-%s", dirpath, uid, filename)
				output_size := "200x200"
				output_align := "center"
				background := "white"

				e := helper.Thumbnail("resize", inputFile, output_file, output_size, output_align, background)
				if e != nil {
					self.Flash.Error(e.Error())
					return self.Redirect(fmt.Sprintf("/settings/%s/", name))

				} else {

					//所有文件以该加密方式哈希生成文件名  从而实现针对到用户个体的文件权限识别
					filehash, _ := helper.Filehash(helper.URL2local(output_file), nil)

					fname := helper.EncryptHash(filehash+strconv.Itoa(int(uid)), nil)

					newpath := dirpath + fname + ext

					if err := os.Rename(helper.URL2local(output_file), helper.URL2local(newpath)); err != nil {
						log.Println("重命名文件出错:", err)
					} else {
						os.Remove(helper.URL2local(output_file))
					}

					//文件权限校验 通过说明文件上传转换过程中没发生错误
					//首先读取被操作文件的hash值 和 用户请求中的文件hash值  以及 用户当前id的string类型  进行验证

					if fhashed, _ := helper.Filehash(helper.URL2local(newpath), nil); helper.ValidateHash(fname, fhashed+strconv.Itoa(int(uid))) {

						//收到的头像图片存储都设置ctype为 10 与其他图片类型区分开
						/*
							if _, err := models.AddImage(helper.URL2local(newpath), 0, 10, uid); err != nil {
								log.Print("models.AddImage:", err)
							}
						*/
						usr, _ := models.GetUser(uid)

						//清理旧头像文件
						if usr.Avatar != "" {
							//os.Remove(helper.URL2local(usr.Avatar))
							os.Remove(helper.URL2local(targetFolder + usr.Avatar))
						}
						if usr.AvatarLarge != "" {
							//os.Remove(helper.URL2local(usr.AvatarLarge))
							os.Remove(helper.URL2local(targetFolder + usr.AvatarLarge))
						}
						if usr.AvatarMedium != "" {
							//os.Remove(helper.URL2local(usr.AvatarMedium))
							os.Remove(helper.URL2local(targetFolder + usr.AvatarMedium))
						}
						if usr.AvatarSmall != "" {
							//os.Remove(helper.URL2local(usr.AvatarSmall))
							os.Remove(helper.URL2local(targetFolder + usr.AvatarSmall))
						}

						//usr.Avatar = newpath //200x200
						usr.Avatar = dirpath[len(targetFolder):] + fname + ext //200x200
						//log.Println("usr.Avatar...>", usr.Avatar)
						thumbnailpath := strings.Split(usr.Avatar, ".")[0]

						usr.AvatarLarge = thumbnailpath + "_large.jpg" //100x100
						helper.Thumbnail("resize", helper.URL2local(newpath), helper.URL2local(targetFolder+usr.AvatarLarge), "100x100", output_align, background)

						usr.AvatarMedium = thumbnailpath + "_medium.jpg" //48x48
						helper.Thumbnail("resize", helper.URL2local(newpath), helper.URL2local(targetFolder+usr.AvatarMedium), "48x48", output_align, background)

						usr.AvatarSmall = thumbnailpath + "_small.jpg" //32x32
						helper.Thumbnail("resize", helper.URL2local(newpath), helper.URL2local(targetFolder+usr.AvatarSmall), "32x32", output_align, background)

						models.PutUser(uid, usr)

						self.Session.Set("SignedUser", usr)
						_usr_ = usr
						self.Set("SignedUser", _usr_)
						cc.Set(fmt.Sprintf("SignedUser:%v", usr.Id), usr, 60*60*24)
						self.Flash.Success(fmt.Sprint("成功设置头像为:", mf.Filename))

						return self.Redirect(fmt.Sprintf("/settings/%s/", name))

					} else {

						self.Flash.Error("文件权限校验失败!")
						if e := os.Remove(helper.URL2local(newpath)); e != nil {
							log.Println("SettingsHandler清除错误遗留文件", newpath, "出错:", e)
						}

						return self.Redirect(fmt.Sprintf("/settings/%s/", name))

					}

				}

			}

		}
	/*
		case name == "avatars":
			{

				suserid, _ := self.Session.Get("userid").(int64)
				images := template.HTMLEscapeString(strings.TrimSpace(self.Query("images")))
				if images != "" {
					imgs := helper.Split(images, ",")

					if usr, e := models.GetUser(suserid); usr != nil && e == nil && imgs[0] != "" {

						if usr.Avatar != "" {
							////删除七牛云存储中的图片
							//http://yougam.qiniudn.com/2014-8-21-134506AB4F24A56EF60543.jpg?imageView/2/w/200/h/200/q/100
							imgkey := strings.Split(usr.Avatar, "?")
							imgkey2 := strings.Split(imgkey[0], "/")
							imgkey3 := strings.Split(imgkey2[len(imgkey2)-1], ".")
							delkey := imgkey3[0] //key是32位  不含后缀

							rsCli := rs.New(nil)
							rsCli.Delete(nil, BUCKET, delkey)
						}

						usr.Avatar = "http://" + helper.DOMAIN4QINIU + "/" + imgs[0] + "?imageView/2/w/200/h/200/q/100"
						usr.AvatarLarge = "http://" + helper.DOMAIN4QINIU + "/" + imgs[0] + "?imageView/2/w/100/h/100/q/100"
						usr.AvatarMedium = "http://" + helper.DOMAIN4QINIU + "/" + imgs[0] + "?imageView/2/w/48/h/48/q/100"
						usr.AvatarSmall = "http://" + helper.DOMAIN4QINIU + "/" + imgs[0] + "?imageView/2/w/32/h/32/q/100"
						self.Session.Set("useravatar", usr.Avatar)

						if _, e := models.PutUser(suserid, usr); e != nil {
							self.Flash.Error("保存头像失败~头像信息写入数据库出错~")
						} else {
							self.Set["user"] = usr
							self.Session.Set("useravatar", usr.Avatar)
							self.Flash.Success("保存头像成功~")

						}

					}

				} else {
					self.Flash.Error("保存头像失败~你尚未上传头像图片~")
				}

				return self.Redirect(fmt.Sprintf("/settings/%s/", name), 302)
				return
			}
	*/
	case name == "password":
		{

			curpass := self.QueryParam("curpass")
			newpassword := self.QueryParam("password")
			newrepassword := self.QueryParam("repassword")

			if curpass == "" {
				self.Flash.Error("当前密码不能为空!")
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			}

			if newpassword == "" {
				self.Flash.Error("新设密码不能为空!")

				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			}

			if newrepassword == "" {
				self.Flash.Error("确认密码不能为空!")
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			}

			if newpassword != newrepassword {
				self.Flash.Error("两次密码不一致!")
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			}

			if helper.CheckPassword(curpass) == false {
				self.Flash.Error("当前密码含有非法字符或当前密码过短(至少4~30位密码)!")
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			}

			if helper.CheckPassword(newpassword) == false {
				self.Flash.Error("设置密码含有非法字符或设置密码过短(至少4~30位密码)!")
				return self.Redirect(fmt.Sprintf("/settings/%s/", name))
			}

			if usrinfo, err := models.GetUser(_usr_.Id); usrinfo != nil && err == nil {

				if helper.ValidateHash(usrinfo.Password, curpass) {
					usrinfo.Password = helper.EncryptHash(newpassword, nil)

					if _, err := models.PutUser(usrinfo.Id, usrinfo); err == nil {
						//更新self.Session.on
						self.Session.Set("SignedUser", usrinfo)
						_usr_ = usrinfo
						self.Set("SignedUser", _usr_)
						cc.Set(fmt.Sprintf("SignedUser:%v", _usr_.Id), _usr_, 60*60*24)
						self.Flash.Success("设置密码成功~")
					} else {
						self.Flash.Error("设置密码失败~")
					}

					return self.Redirect(fmt.Sprintf("/settings/%s/", name))

				} else {

					self.Flash.Error("当前密码无法通过校验~")
					return self.Redirect(fmt.Sprintf("/settings/%s/", name))

				}
			} else {

				self.Flash.Error("该账号不存在~")

				return self.Redirect(fmt.Sprintf("/settings/%s/", name))

			}

		}
	default:
		{

			return self.Redirect("/")

		}

	}

	return self.Redirect("/")

}
