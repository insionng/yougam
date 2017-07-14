package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"strconv"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func PostNewReplyHandler(self *makross.Context) error {

	//TODO:暂时登录才能发评论，以后在后台增加允许设置是否对游客开放评论

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}
	cc := cache.Store(self)

	tid := self.Param("tid").MustInt64()
	pid := self.Param("pid").MustInt64()

	author := self.FormValue("author")
	nickname := self.FormValue("nickname")
	email := self.FormValue("email")
	website := self.FormValue("website")

	rc := self.FormValue("comment")
	if len(rc) > 0 {
		policy := helper.ObjPolicy()
		rc = policy.Sanitize(rc)
	}

	images := self.FormValue("images")

	if len(_usr_.Avatar) > 0 {
		self.Set("hasavatar", true)
	} else {
		self.Set("hasavatar", false)
		self.Flash.Warning("请设置你的头像,有头像能更好的引起其他成员的注意~")

	}

	if _usr_.Id > 0 {
		if (tid > 0) && (len(rc) > 0) {

			allow := bool(false)
			if _usr_.Balance >= 1 {
				allow = true
			} else {
				self.Flash.Error(fmt.Sprintf("你的余额为[%v]，不足以评论!", _usr_.Balance))
			}

			if tp, err := models.GetTopic(tid); err == nil && allow {

				if tp.Pid != 0 {
					self.Flash.Error("子话题不允许评论!")
					return self.Redirect("/topic/" + self.Param("tid").String() + "/")
				}

				//为安全计,先行保存回应,顺手获得rid
				//ctype不等于0,即是注册用户或管理层的回复 此时把ctype设置为1 主要是为了区分游客
				if rid, err := models.AddReply(tid, pid, _usr_.Id, 1, rc, images, _usr_.Username, _usr_.Avatar, _usr_.AvatarLarge, _usr_.AvatarMedium, _usr_.AvatarSmall, _usr_.Content, _usr_.Nickname, _usr_.Email, _usr_.Website); err != nil {
					self.Flash.Error(fmt.Sprint("#", rid, ":", err))
				} else {
					//更新统计缓存 这里不需要全局计算 意思下就行 在PrepareHandler中间件会定期统计全局数据
					if cc.IsExist("ReplysCount") {
						var ReplysCount int
						cc.Get("ReplysCount", &ReplysCount)
						ReplysCount = ReplysCount + 1
						cc.Set("ReplysCount", ReplysCount, 60*60*2)
						self.Set("ReplysCount", ReplysCount)
					} else {
						ReplysCount := 1
						cc.Set("ReplysCount", ReplysCount, 60*60*2)
						self.Set("ReplysCount", ReplysCount)
					}

					//是否允许写入标记
					var allow = false
					switch kind := tp.Ctype; {
					/*
					   tp.Ctype ==  0 [普通话题]
					   tp.Ctype == -1 [回复可见]
					   tp.Ctype == -2 [付费可见]
					   tp.Ctype == -3 [会员可见]
					*/
					case kind == 0:
						allow = false //如果是普通类型则无须写入标记
					case kind == -1:
						allow = true      //须要回复 并写入标记
						if _usr_ != nil { //若果是登录状态
							if _usr_.Id == tp.Uid { //如果当前用户即是作者则无须写入标记
								allow = false
							}

							if _usr_.Role == -1000 { //如果是管理层则无须写入标记
								allow = false
							}
						}
					case kind == -2:
						allow = false //须要另外付费 并写入标记
					case kind == -3:
						allow = false //另外判断即可 无须写入标记
					}

					if allow {
						//添加UID标记
						models.PutTailinfo2Topic(tid, fmt.Sprint(_usr_.Id))
					}

					hasSend := false //判断作者是否已经通知
					if users := helper.AtUsers(rc); len(users) > 0 {
						for _, v := range users {
							//不是评论者本人自己@自己则继续执行
							if _usr_.Username != v {
								//判断被通知之用户名是否真实存在
								if u, e := models.GetUserByUsername(v); e == nil && u != nil {

									n_, e_ := models.AddNotification(tp.Id, rid, u.Id, 0, helper.Substr(helper.HTML2str(tp.Title), 0, 100, "..."), helper.Substr(helper.HTML2str(rc), 0, 200, "..."), _usr_.Username, _usr_.Avatar, _usr_.AvatarLarge, _usr_.AvatarMedium, _usr_.AvatarSmall)
									if (n_ != -1) && (e_ == nil) {
										_u, _e := models.GetUser(u.Id)
										if (_u != nil) && (_e == nil) {
											cc.Set(fmt.Sprintf("SignedUser:%v", u.Id), _u, 60*60*24)
										}
									}
								}
								if v == tp.Author {
									hasSend = true
								}
							}
						}
					}

					if e := models.SetAmountByUid(_usr_.Id, -1, -1, "创建回复,付出1金币"); e == nil {
						if usr, e := models.GetUser(_usr_.Id); e == nil && usr != nil {
							cc.Set(fmt.Sprintf("SignedUser:%v", usr.Id), usr, 60*60*24)
						}
					}

					if tp.Uid != _usr_.Id { //若果当前评论用户不是主题作者则处理
						models.SetAmountByUid(tp.Uid, 1, +1, "主题被评论,收益1金币")
					} else { //若果当前评论用户是主题作者则处理
						hasSend = true //不对自己发通知
					}

					if !hasSend {
						//通知话题作者
						n_, e_ := models.AddNotification(tp.Id, rid, tp.Uid, 0, helper.Substr(helper.HTML2str(tp.Title), 0, 100, "..."), helper.Substr(helper.HTML2str(rc), 0, 200, "..."), _usr_.Username, _usr_.Avatar, _usr_.AvatarLarge, _usr_.AvatarMedium, _usr_.AvatarSmall)
						if (n_ != -1) && (e_ == nil) {
							_u, _e := models.GetUser(tp.Uid)
							if (_u != nil) && (_e == nil) {
								cc.Set(fmt.Sprintf("SignedUser:%v", tp.Uid), _u, 60*60*24)
							}
						}

					}

					return self.Redirect("/topic/" + self.Param("tid").String() + "/#reply" + strconv.FormatInt(rid, 10))
				}
			}

			return self.Redirect("/topic/" + self.Param("tid").String() + "/")

		} else if tid > 0 {
			return self.Redirect("/topic/" + self.Param("tid").String() + "/")
		} else {
			return self.Redirect("/")
		}
	} else { //游客回应 此时把ctype设置为-1   游客不开放@通知功能
		if len(author) > 0 && len(email) > 0 && tid > 0 && len(rc) > 0 {
			if rid, err := models.AddReply(tid, pid, _usr_.Id, -1, rc, images, author, "", "", "", "", "", nickname, email, website); err != nil {
				self.Flash.Error(fmt.Sprint("#", rid, ":", err))
				return self.Redirect("/topic/" + self.Param("tid").String() + "/")

			} else {
				return self.Redirect("/topic/" + self.Param("tid").String() + "/#reply" + strconv.Itoa(int(rid)))

			}
		} else if tid > 0 {
			return self.Redirect("/topic/" + self.Param("tid").String() + "/")

		} else {
			return self.Redirect("/")

		}

	}
}
