package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"github.com/insionng/yougam/models"
)

func GetRRechargeUserHandler(self *makross.Context) error {
	

	TplNames := "root/recharge_user"
	self.Set("catpage", "RRechargeUserHandler")
	uid := self.Param("uid").MustInt64()
	if uid <= 0 {
		self.Flash.Error("非法参数!")
		return self.Redirect("/root/read/user/")

	}

	if usr, err := models.GetUser(uid); usr != nil && err == nil {
		self.Set("usr", *usr)

	} else {
		self.Flash.Error(fmt.Sprint(err))
		return self.Render(TplNames)

	}
	return self.Render(TplNames)
}

func PostRRechargeUserHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	amount := self.Args("amount").MustInt64()
	uid := self.Param("uid").MustInt64()
	if uid <= 0 {
		self.Flash.Error("非法参数!")
		return self.Redirect("/root/read/user/")

	}

	usr, e := models.GetUser(uid)
	if e != nil || usr == nil {
		self.Flash.Error("用户不存在!")
		return self.Redirect("/root/read/user/")

	}

	if e := models.SetAmountByUid(uid, 3, amount, "充值"); e != nil {
		self.Flash.Error(fmt.Sprintf("充值用户[%v]失败!", usr.Username))

	} else {
		self.Flash.Success(fmt.Sprintf("充值用户[%v]成功!", usr.Username))
		if _usr_ != nil {
			if _usr_.Id == usr.Id {
				ur, e := models.GetUser(uid)
				if e != nil || ur == nil {
					self.Flash.Error("用户不存在!")
					return self.Redirect("/root/read/user/")

				}

				self.Session.Set("SignedUser", ur)
				self.Set("user", ur)
			}
		}
	}

	return self.Redirect(fmt.Sprintf("/root/update/user/recharge/%v/", uid))

}
