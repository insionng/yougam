package setting

import (
	"fmt"
	"runtime"

	"github.com/insionng/makross"
	"github.com/insionng/yougam/helper"
	simplejson "github.com/insionng/yougam/libraries/bitly/go-simplejson"
	"github.com/insionng/yougam/models"
)

// AuthWebMiddler 会员或管理员前台权限认证
func AuthWebMiddler() makross.Handler {
	return func(self *makross.Context) error {
		var user = new(models.User)
		if sUser, okay := self.Session.Get("SignedUser").(*models.User); okay {
			user = sUser
		}
		if !(user != nil) {
			if nx := self.Args("next").String(); len(nx) != 0 {
				self.Abort()
				return self.Redirect(fmt.Sprintf("/sigin/?next=%v", nx))
			} else {
				self.Abort()
				return self.Redirect(fmt.Sprintf("/sigin/?next=%v", self.RequestURI()))
			}
		}
		return self.Next()
	}
}

// RootMiddler 管理员后台后台认证
func RootMiddler() makross.Handler {
	return func(self *makross.Context) error {
		var IsRoot, IsSignin bool
		if sUser, okay := self.Session.Get("SignedUser").(*models.User); okay {
			IsSignin = true
			IsRoot = (sUser.Role == -1000)
		}
		if IsSignin {
			if !IsRoot {
				self.Abort()
				return self.Redirect("/root/signin/")
			} else {
				self.Set("remoteproto", self.Scheme())
				self.Set("remotehost", self.RealIP())
				self.Set("remoteos", runtime.GOOS)
				self.Set("remotearch", runtime.GOARCH)
				self.Set("remotecpus", runtime.NumCPU())
				return nil
			}
		}
		self.Abort()
		return self.Redirect("/")
	}
}

// APISessionMiddler Session级权限认证
func APISessionMiddler() makross.Handler {
	return func(self *makross.Context) error {
		if _, okay := self.Session.Get("SignedUser").(*models.User); !okay {
			//返回401未认证状态终止服务
			return self.NoContent(401)
		}
		return self.Next()
	}
}

// APICryptMiddler AES128COM加密验证+SESSION(客户端须开启COOKIES)权限认证
func APICryptMiddler() makross.Handler {
	return func(self *makross.Context) error {

		//验证加密请求是否以form data的形式提交
		var datas string
		if dt := self.FormValue("data"); len(dt) > 0 {
			datas = dt
			//fmt.Println("datas:", datas)
		} else {
			//如果不是form data  则设为http self.Req.Body
			b := []byte(nil)
			self.Write(b)
			datas = string(b)
		}

		if len(datas) == 0 {
			crypted, _ := helper.SetJsonCOMEncrypt(0, "提交的数据为空!", nil)
			return self.String(crypted)
		}

		if s, err := helper.Aes128COMDecrypt(datas, helper.AesConstKey); err != nil {
			crypted, _ := helper.SetJsonCOMEncrypt(0, "无法通过安全校验!", nil)
			return self.String(crypted)
		} else {
			self.Set("decrypts", s)
			if j, err := simplejson.NewJson([]byte(s)); err != nil {

				crypted, _ := helper.SetJsonCOMEncrypt(0, err.Error(), nil)
				return self.String(crypted)

			} else {
				if action, err := j.Get("action").String(); err == nil {
					isPass := bool(false)
					//以下请求动作均跳过self.Session.限检查
					if action == "userSignup" || action == "userSignin" || action == "userSignout" || action == "getHomePostList" || action == "getContent" || action == "getComment" || action == "getUserPostList" {
						isPass = true
					}

					//客户端方面除了上面跳过的请求,其他的每个请求都需要附带含self.Session.d的cookies头
					if /*(self.Session.role == 0) && */ !isPass { //此处即为未在cookies中附带self.Session.d
						crypted, _ := helper.SetJsonCOMEncrypt(0, "尚未登录认证!", nil)
						return self.String(crypted)

					}

				} else {
					crypted, _ := helper.SetJsonCOMEncrypt(0, err.Error(), nil)
					return self.String(crypted)
				}

			}
		}
		return self.Next()
	}
}
