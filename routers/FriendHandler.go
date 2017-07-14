package routers

import (
	"errors"
	"github.com/insionng/makross"

	"log"
	"github.com/insionng/yougam/models"
)

func GetFriendHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	fid := self.Param("uid").MustInt64()
	if fid <= 0 {
		fid = self.Args("uid").MustInt64()
	}

	if fid <= 0 {
		self.Flash.Error("错误参数!")
		return self.Redirect("/contact/")
	}

	if kind := self.Param("kind").String(); len(kind) > 0 && (fid > 0) {
		switch {
		case kind == "add":
			{

				if _usr_.Id == fid {
					self.Flash.Error("不能添加自己为好友!")
					return self.Redirect("/contact/")
				}

				usr, err := models.GetUser(fid)
				if (err != nil) || (usr == nil) {
					log.Println(err)
					if err != nil {
						self.Flash.Error(err.Error())
					}
					return self.Redirect("/contact/")
				}

				if models.IsFriend(_usr_.Id, fid) {
					self.Flash.Warning("已经是好友,无须再度添加!")
					return self.Redirect("/contact/")

				} else {
					self.Set("theUser", usr)
					self.Set("relationship", models.GetRelationship(_usr_.Id, fid))
					goto render

				}

			render:
				self.Set("messager", true)
				return self.Render("add-friend")
			}
		case kind == "delete":
			{
				usr, err := models.GetUser(fid)
				if (err != nil) || (usr == nil) {
					log.Println(err)
					return self.Redirect("/contact/")
				}

				r, e := models.DelFriend(_usr_.Id, fid)
				if (e != nil) || (r <= 0) {
					self.Flash.Error("删除好友失败!")
				} else {
					self.Flash.Success("删除好友成功!")
				}
				return self.Redirect("/contact/")
			}
		case kind == "allow":
			{
				r, e := models.SetFriend(_usr_.Id, fid, 2, "", "default")
				if (e != nil) || (r <= 0) {
					log.Println(e)
					self.Flash.Error("允许好友申请时发生错误!")
				} else {
					friz := models.GetFriendsByUidJoinUser(_usr_.Id, 0, 0, "", "id")
					self.Set("friends", friz)
					self.Flash.Success("经已允许好友申请！今天下英雄﹐惟使君与操耳。本初之徒﹐不足数也。")

				}

				return self.Redirect("/contact/")
			}
		case kind == "deny":
			{
				r, e := models.SetFriend(_usr_.Id, fid, -1, "", "default")
				if (e != nil) || (r <= 0) {
					log.Println(e)
					self.Flash.Error("拒绝好友申请时发生错误!")
				} else {
					friz := models.GetFriendsByUidJoinUser(_usr_.Id, 0, 0, "", "id")
					self.Set("friends", friz)
					self.Flash.Success("经己拒绝好友申请！")

				}

				return self.Redirect("/contact/")
			}
		}

	}

	return errors.New("FriendGetHandler Errors")
}

func PostFriendHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	fid := self.Param("uid").MustInt64()
	if fid <= 0 {
		fid = self.Args("uid").MustInt64()
	}

	if fid <= 0 {
		self.Flash.Error("错误参数!")
		return self.Redirect("/contact/")
	}

	if kind := self.Param("kind").String(); len(kind) > 0 && (fid > 0) {
		switch {
		case kind == "add":
			{

				content := self.FormValue("content")
				r, e := models.SetFriendTo(_usr_.Id, fid, 1, content, "default")
				if (e != nil) || (r <= 0) {
					log.Println(e)
					self.Flash.Error("发送好友申请错误!")
					goto render
				} else {
					friz := models.GetFriendsByUidJoinUser(_usr_.Id, 0, 0, "", "id")
					self.Set("friends", friz)
					self.Flash.Success("发送好友申请成功，请耐心等待对方通过!")
					return self.Redirect("/contact/")
				}

			render:
				self.Set("messager", true)
				return self.Render("add-friend")
			}
		}
	}

	return errors.New("FriendPostHandler Errors")

}
