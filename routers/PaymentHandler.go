package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func PostPaymentHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}
	cc := cache.Store(self)

	tid := self.Param("id").MustInt64()
	toName := self.Param("name").String()
	amount := self.Param("amount").MustInt64()
	password := self.FormValue("password")
	if len(password) <= 0 {
		self.Flash.Error("密码不能为空！")
		return self.Redirect(fmt.Sprintf("/topic/%v/", tid))
	}

	if allow := bool(false); (tid > 0) && (amount > 0) {
		if _usr_.Balance >= amount {
			allow = true
		} else {
			self.Flash.Error(fmt.Sprintf("你的余额为[%v]，不足以支付!", _usr_.Balance))
			return self.Redirect(fmt.Sprintf("/topic/%v/", tid))
		}

		if tp, err := models.GetTopic(tid); (err == nil) && allow {

			if (tp.Author != toName) && (tp.Ctype != -2) {
				self.Flash.Error("非法参数！")
				return self.Redirect(fmt.Sprintf("/topic/%v/", tid))
			}

			usr, e := models.GetUser(_usr_.Id)
			if (e == nil) && (usr != nil) {
				if !helper.ValidateHash(usr.Password, password) {
					self.Flash.Error("密码无法通过校验！")
					return self.Redirect(fmt.Sprintf("/topic/%v/", tid))
				}
			}

			if e := models.SetAmountByUid(_usr_.Id, -3, -amount, fmt.Sprintf("支付 %v 金币", amount)); e == nil {

				cc.Set(fmt.Sprintf("SignedUser:%v", _usr_.Id), usr, 60*60*24)

				if tp.Uid != _usr_.Id { //若果当前支付用户不是主题作者则处理
					models.SetAmountByUid(tp.Uid, 5, +amount, fmt.Sprintf("话题卖出，收益 %v 金币", amount))
					if otherusr, e := models.GetUser(tp.Uid); (e == nil) && (otherusr != nil) {
						cc.Set(fmt.Sprintf("SignedUser:%v", tp.Uid), otherusr, 60*60*24)
					}
				} else {
					self.Flash.Error("不能自我支付！")
					return self.Redirect(fmt.Sprintf("/topic/%v/", tid))
				}

				//添加UID标记
				_, e := models.PutTailinfo2Topic(tid, fmt.Sprint(_usr_.Id))
				if e != nil {
					self.Flash.Error("支付失败！")
					return self.Redirect(fmt.Sprintf("/topic/%v/", tid))
				}
				self.Flash.Success("成功支付！")
				return self.Redirect(fmt.Sprintf("/topic/%v/", tid))

			} else {
				self.Flash.Error("支付失败！")
				return self.Redirect(fmt.Sprintf("/topic/%v/", tid))

			}
		}

	}

	return self.NoContent(404)

}
