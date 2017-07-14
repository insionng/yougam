package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"time"

	"github.com/insionng/yougam/models"
)

func GetSignoutHandler(self *makross.Context) error {

	_usr_, _ := self.Session.Get("SignedUser").(*models.User)
	cc := cache.Store(self)

	if _usr_ != nil {
		//销毁self.Session.on
		self.Session.Delete("User")
		self.Session.Set("SignedUser", nil)

		cc.Set(fmt.Sprintf("SignedUser:%v", _usr_.Id), nil, 60*60*24)
		cc.Delete(fmt.Sprintf("User:%v", _usr_.Id))
		_usr_ = nil
	}

	//退出时关闭记住密码
	//Secret := helper.MD5(self.Req.UserAgent() + helper.AesConstKey)
	//self.SetSuperSecureCookie(Secret, "flower", "", 3600) //删除数据

	cookie := self.NewCookie()
	cookie.Secure = (true)
	cookie.Name = ("remember")
	cookie.Value = ("false")
	cookie.Expires = (time.Now().Add(1 * time.Minute))
	self.SetCookie(cookie)

	if next := self.FormValue("next"); next != "" {
		return self.Redirect(next)
	}

	return self.Redirect(fmt.Sprintf("/?version=%v", time.Now().Unix()))

}
