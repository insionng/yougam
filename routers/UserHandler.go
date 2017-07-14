package routers

import (
	"fmt"
	"github.com/insionng/makross"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetUserHandler(self *makross.Context) error {

	_usr_, _ := self.Session.Get("SignedUser").(*models.User)

	TplNames := "user/index"
	name := self.Args("name").String()
	ktype := self.Args("type").String()
	tab := self.Args("tab").String()
	ipage := self.Args("page").MustInt64()
	offset := self.Args("offset").MustInt64()

	limit := int64(10)
	ctype := int64(0)
	field := "id"
	url := "/"

	if tab == "lastest" {
		url = "/lastest/"
		field = "id"
		self.Set("tab", "lastest")
	} else if tab == "hotness" {
		url = "/hotness/"
		field = "hotness"
		self.Set("tab", "hotness")
	} else {
		url = "/lastest/"
		field = "id"
		self.Set("tab", "lastest")
	}

	if usr, err := models.GetUserByUsername(name); usr != nil && err == nil {
		self.Set("userProfile", *usr)
		switch {
		case ktype == "favorites":
			{

				TplNames = "user/favorites"
				url = "/user/" + name + "/" + ktype + url

				if tps := models.JoinTopicmarkJoinUserForGetTopicsByUid(usr.Id, int(offset), 0, field); tps != nil {
					if totalRecords := int64(len(*tps)); totalRecords > 0 {

						pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, ipage, limit)

						if tps := models.JoinTopicmarkJoinUserForGetTopicsByUid(usr.Id, int(offset), int(limit), field); tps != nil {

							self.Set("favorites", *tps)
							self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))

						}

					}

				}

				return self.Render(TplNames)

			}
		case ktype == "topics":
			{

				TplNames = "user/topics"
				url = "/user/" + name + "/" + ktype + url
				if totalRecords, e := models.GetSubjectsCountByUsername(name, int(offset), int(limit), ctype); totalRecords > 0 && e == nil {

					pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, ipage, limit)

					if tps := models.GetSubjectsByUsername(name, int(offset), int(limit), ctype, field); *tps != nil {
						self.Set("topics", *tps)
						self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))
					}
				}

				return self.Render(TplNames)
			}
		case ktype == "replys":
			{

				TplNames = "user/replys"
				url = "/user/" + name + "/" + ktype + url

				if rps := models.GetReplysByTidUsername(0, name, ctype, int(offset), 0, field); rps != nil {

					if totalRecords := int64(len(*rps)); totalRecords > 0 {
						pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, ipage, limit)

						if rpys := models.GetReplysByTidUsernameJoinTopic(0, name, ctype, int(offset), int(limit), field); rpys != nil {
							self.Set("replys", rpys)
							self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))

						}
					}

				}

				return self.Render(TplNames)

			}
		case ktype == "balance":
			{

				if _usr_.Username != name {
					self.Flash.Warning(fmt.Sprintf("警告，请勿非法访问 %v 的钱包！", name))
					url = fmt.Sprintf("/user/%v/%v/", _usr_.Username, ktype)
					return self.Redirect(url)

				}

				TplNames = "user/balance"
				url = "/user/" + name + "/" + ktype + url
				if balances := models.GetBalancesByUsername(name, 0, 0, 0, field); balances != nil {

					if totalRecords := int64(len(*balances)); totalRecords > 0 {
						pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, ipage, limit)

						if balances := models.GetBalancesByUsername(name, ctype, int(offset), int(limit), field); balances != nil {

							self.Set("balances", balances)
							self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))

						}
					}
				}
				return self.Render(TplNames)

			}
		default:
			{

				if tps := models.JoinTopicmarkJoinUserForGetTopicsByUid(usr.Id, 0, 10, "id"); tps != nil {
					self.Set("favorites", *tps)
				}

				if tps := models.GetSubjectsByUsername(name, 0, 10, 0, "id"); *tps != nil {
					self.Set("topics", *tps)
				}

				if rpys := models.GetReplysByTidUsernameJoinTopic(0, name, 0, 0, 10, "id"); rpys != nil {
					self.Set("replys", *rpys)
				}

				return self.Render(TplNames)
			}
		}
	} else {
		return self.Redirect("/")

	}
}
