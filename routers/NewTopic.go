package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"log"
	"strconv"
	"time"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetNewTopicHandler(self *makross.Context) error {

	var _usr_ = new(models.User)
	var IsSignin bool
	if sUser, okay := self.Session.Get("SignedUser").(*models.User); okay {
		_usr_ = sUser
		IsSignin = okay
	}

	if IsSignin {
		TplNames := "new-topic"
		if len(_usr_.Avatar) > 0 {
			//获取所有针对前端可用之节点
			//if nds, _ := models.AvailableNodes(0, 0, "hotness"); len(*nds) <= 0 {
			if nds, _ := models.GetNodes(0, 0, "hotness"); len(*nds) <= 0 {
				//不存在节点跳转到首页,须要创建至少一个默认节点
				self.Flash.Error("不存在节点,须要创建至少一个默认节点~", false)
				return self.Redirect("/new/node/")

			} else {
				if cats, err := models.GetCategoriesByNodeCount(0, 0, 0, "id"); cats != nil && err == nil {
					self.Set("categories", *cats)
				}

				/*
					if nds_, err := models.NodesOfNavor(0, 0, "hotness"); nds_ != nil && err == nil {
						self.Set("nodes", *nds_)
					}
				*/
				self.Set("nodes", *nds)

				self.Set("hasnid", false)
				node := self.Param("node").String()
				if len(node) > 0 {
					if nd, e := models.GetNodeByTitle(node); nd != nil && e == nil {
						self.Set("hasnid", true)
						self.Set("curnid", nd.Id)
					}
				} else {
					nid := self.Param("nid").MustInt64()
					if nid > 0 {
						if nd, e := models.GetNode(nid); nd != nil && e == nil {
							self.Set("hasnid", true)
							self.Set("curnid", nid)
						}
					}
				}

				self.Set("haspid", false)
				self.Set("catpage", "NewTopicHandler")

				tid := self.Param("tid").MustInt64()
				if tid > 0 {
					if tp, e := models.GetTopic(tid); e == nil && tp != nil {

						if _usr_.Id != tp.Uid { //如果当前用户不是主题的作者

							sjid := int64(0)
							if tp.Pid != 0 {
								sjid = tp.Pid
							} else {
								sjid = tp.Id
							}

							self.Flash.Error("你并非作者,无权为他人主题增添附言!", false)
							return self.Redirect("/topic/" + strconv.FormatInt(sjid, 10) + "/")
						}

						self.Set("topic", tp)
						self.Set("haspid", true)
						self.Set("curpid", tid)

						if tps := models.GetTopicsByPid(tid, 0, 0, 0, "id"); tps != nil && (len(*tps) > 0) {
							self.Set("topics", *tps)
						} else {
							return self.Redirect("/subject/" + strconv.FormatInt(tid, 10) + "/topic/")
						}
					}
				}

			}

			return self.Render(TplNames)
		} else {
			//没有头像 不许发布话题 跳转到头像设置页面
			self.Flash.Error("请设置你的个人头像,没有头像发布的主题将难以引起其他成员的注意~", false)
			return self.Redirect("/settings/avatar/")
		}
	} else {
		return self.Redirect("/")
	}

}

func PostNewTopicHandler(self *makross.Context) error {

	TplNames := "new-topic"

	var _usr_ = new(models.User)
	if sUser, okay := self.Session.Get("SignedUser").(*models.User); okay {
		_usr_ = sUser
	}
	cc := cache.Store(self)

	subjectid := self.Param("tid").MustInt64()
	nid := self.Args("nodeid").MustInt64()
	ctype := self.Args("ctype").MustInt64()
	images := self.Args("images").String()

	policy := helper.ObjPolicy()

	tptitle := self.Args("title").String()
	if len(tptitle) > 0 {
		tptitle = policy.Sanitize(tptitle)
	}

	tpexcerpt := self.Args("excerpt").String()
	if len(tpexcerpt) > 0 {
		tpexcerpt = policy.Sanitize(tpexcerpt)
	}

	tpcontent := self.Args("content").String()
	if len(tpcontent) > 0 {
		tpcontent = policy.Sanitize(tpcontent)
	}

	if subjectid <= 0 { //新建主题

		cid := int64(0)
		nd, err := models.GetNode(nid)
		if err != nil || nd == nil || nid == 0 {

			self.Flash.Error("节点不存在,请创建或指定正确的节点!")
			goto render
		} else {
			cid = nd.Cid
		}

		if len(tptitle) > 0 {

			tp := new(models.Topic)
			tp.Title = tptitle
			tp.Content = tpcontent
			tp.Excerpt = tpexcerpt
			tp.Cid = cid
			tp.Nid = nid
			tp.Uid = _usr_.Id
			tp.Ctype = ctype
			tp.Node = nd.Title
			tp.Author = _usr_.Username
			tp.Avatar = _usr_.Avatar
			tp.AvatarLarge = _usr_.AvatarLarge
			tp.AvatarMedium = _usr_.AvatarMedium
			tp.AvatarSmall = _usr_.AvatarSmall
			tp.Created = time.Now().Unix()
			/*else {
				if s, e := helper.GetBannerThumbnail(tpcontent); e == nil {
					tp.Attachment = s
				}
			}

				if thumbnails, thumbnailslarge, thumbnailsmedium, thumbnailssmall, e := helper.GetThumbnails(tpcontent); e == nil {
					tp.Thumbnails = thumbnails
					tp.ThumbnailsLarge = thumbnailslarge
					tp.ThumbnailsMedium = thumbnailsmedium
					tp.ThumbnailsSmall = thumbnailssmall
				}
			*/
			if cat, err := models.GetCategory(cid); err == nil && cat != nil {
				tp.Category = cat.Title
			}

			nodezmap := map[string]interface{}{
				"topic_time":         time.Now().Unix(),
				"topic_count":        models.GetTopicCountByNid(nid),
				"topic_last_user_id": _usr_.Id}

			if e := models.UpdateNode(nid, &nodezmap); e != nil {
				log.Println("NewTopic models.UpdateNode errors:", e)
			}

			if tid, err := models.PostTopic(tp); err == nil {
				if e := models.SetAmountByUid(tp.Uid, 1, 1, "创建主题收益1金币"); e == nil {
					if usr, e := models.GetUser(tp.Uid); (e == nil) && (usr != nil) {
						cc.Set(fmt.Sprintf("User:%v", usr.Id), usr, 60*60*24)
					}
				}

				if len(images) != 0 {
					models.AddAttachment(images, tid, cid, nid, tp.Uid)
				}

				self.Set("curpid", tid)
				self.Flash.Success(fmt.Sprint("主题 #", tid, " [", tp.Title, "] ", "保存成功!"), false)
				//当前用户如果不是主题作者的时候不作处理(如管理员)
				if tp.Uid == _usr_.Id {

					//如果标题或内容中有@通知 则处理以下事件
					if users := helper.AtUsers(fmt.Sprintf("%v%v", tptitle, tpcontent)); len(users) > 0 {

						for _, v := range users {
							//不是作者本人自己@自己则继续执行
							if tp.Author != v {
								//判断被通知之用户名是否真实存在
								if u, e := models.GetUserByUsername(v); e == nil && u != nil {

									n_, e_ := models.AddNotification(tid, 0, u.Id, 0, helper.Substr(helper.HTML2str(tp.Title), 0, 100, "..."), helper.Substr(helper.HTML2str(tpcontent), 0, 200, "..."), _usr_.Username, _usr_.Avatar, _usr_.AvatarLarge, _usr_.AvatarMedium, _usr_.AvatarSmall)
									if (n_ != -1) && (e_ == nil) {
										_u, _e := models.GetUser(u.Id)
										if (_u != nil) && (_e == nil) {
											cc.Set(fmt.Sprintf("User:%v", u.Id), _u, 60*60*24)
										}
									}
								}
							}
						}

					}
				}
				//self.Redirect("/subject/"+strconv.FormatInt(tid, 10)+"/topic/")
				return self.Redirect("/topic/" + strconv.FormatInt(tid, 10) + "/")

			} else {
				self.Set("curpid", tid)
				self.Flash.Error(err.Error())
				goto render

			}
		} else {

			self.Flash.Error("话题标题为空!")
			goto render

		}
	} else { //新加附言 核实主题是否真实存在
		if sj, err := models.GetTopic(subjectid); err == nil && sj != nil {
			if _usr_.Id != sj.Uid { //如果当前用户不是主题的作者
				self.Flash.Error("你并非作者,无权为他人主题增添附言!", false)
				return self.Redirect("/topic/" + strconv.FormatInt(sj.Id, 10) + "/")
			}

			self.Set("curpid", sj.Id)
			if tpcontent != "" {
				allow := bool(false)
				if _usr_.Balance > 0 {
					allow = true
				} else {
					self.Flash.Error("你的余额不足,不允许附言!")
					goto render
				}

				if newtopicid, err := models.AddTopic("", tpcontent, _usr_.Avatar, _usr_.AvatarLarge, _usr_.AvatarMedium, _usr_.AvatarSmall, sj.Id, sj.Cid, sj.Nid, _usr_.Id); (err == nil) && (newtopicid > 0) && allow {
					if e := models.SetAmountByUid(sj.Uid, -2, -1, "创建主题附言,付出1金币"); e == nil {
						if usr, e := models.GetUser(sj.Uid); e == nil && usr != nil {

							cc.Set(fmt.Sprintf("User:%v", usr.Id), usr, 60*60*24)

						}
					}
					if len(images) != 0 {
						models.AddAttachment(images, newtopicid, sj.Cid, sj.Nid, sj.Uid)
					}
					self.Flash.Success("附言保存成功!", false)
					//当前用户如果不是主题作者的时候不作处理(如管理员)
					if sj.Uid == _usr_.Id {

						//如果内容中有@通知 则处理以下事件
						if users := helper.AtUsers(tpcontent); len(users) > 0 {

							for _, v := range users {
								//不是作者本人自己@自己则继续执行
								if sj.Author != v {
									//判断被通知之用户名是否真实存在
									if u, e := models.GetUserByUsername(v); e == nil && u != nil {

										n_, e_ := models.AddNotification(sj.Id, 0, u.Id, 0, helper.Substr(helper.HTML2str(sj.Title), 0, 100, "..."), helper.Substr(helper.HTML2str(tpcontent), 0, 200, "..."), _usr_.Username, _usr_.Avatar, _usr_.AvatarLarge, _usr_.AvatarMedium, _usr_.AvatarSmall)
										if (n_ != -1) && (e_ == nil) {
											_u, _e := models.GetUser(u.Id)
											if (_u != nil) && (_e == nil) {
												cc.Set(fmt.Sprintf("User:%v", u.Id), _u, 60*60*24)
											}
										}
									}
								}
							}

						}
					}
					//self.Redirect("/subject/"+strconv.FormatInt(sj.Id, 10)+"/topic/")
					return self.Redirect("/topic/" + strconv.FormatInt(sj.Id, 10) + "/")
				} else {
					self.Flash.Error(fmt.Sprint("附言内容写入数据库发生错误,", err.Error()))
					goto render
				}
			} else {
				self.Flash.Error("附言内容不能为空!")
				goto render
			}
		} else {
			self.Flash.Error("非法主题!")
			goto render
		}
	}

render:
	return self.Render(TplNames)
}
