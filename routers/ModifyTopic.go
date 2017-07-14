package routers

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"time"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetModifyTopicHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}
	/*
		修改前判断发帖日期,如果超过1天,则不允许再修改
		主题时 title不能为空 content可为空
		附言时,title为空,content不能为空
	*/

	TplNames := "modify-topic"
	self.Set("catpage", "ModifyTopicHandler")
	tid := self.Param("tid").MustInt64()

	if tp, err := models.GetTopic(tid); err == nil && tp != nil {

		if tp.Pid == 0 { //没有上级pid，即本身id就是pid
			self.Set("curpid", tp.Id)
			self.Set("haspid", false)
		} else {
			self.Set("curpid", tp.Pid)
			self.Set("haspid", true)
		}

		if ((tp.Uid == _usr_.Id) && (tp.ReplyCount <= 0)) || (_usr_.Role < 0) { //没有回复的时候才允许修改
			//如果是作者或管理员则可修改
			if len(_usr_.Avatar) > 0 {
				if nds, e := models.GetNodes(0, 0, "id"); nds != nil && e == nil {

					self.Set("nodes", nds)

					//self.Set["hasnid"] = false
					node := self.Param("node").String()
					if node != "" {
						if nd, e := models.GetNodeByTitle(node); nd != nil && e == nil {
							//self.Set["hasnid"] = true
							self.Set("curnid", nd.Id)
						}
					} else {
						nid := self.Param("nid").MustInt64()
						if nid > 0 {
							if nd, e := models.GetNode(nid); nd != nil && e == nil {
								//self.Set("hasnid", true)
								self.Set("curnid", nid)
							}
						}
					}

					self.Set("topic", tp)
					self.Set("curnid", tp.Nid)
					self.Set("curpid", tid)

					if tps := models.GetTopicsByPid(tid, 0, 0, 0, "id"); tps != nil && (len(*tps) > 0) {
						self.Set("topics", *tps)
					} else {
						return self.Redirect("/subject/" + strconv.FormatInt(tid, 10) + "/topic/")
					}

					return self.Render(TplNames)

				} else {
					self.Flash.Error("不存在节点跳转到首页,管理员需要创建至少一个默认节点!")
					return self.Redirect("/")
				}
			} else {
				//没有头像 不许修改话题 跳转到头像设置页面
				self.Flash.Error("请设置你的个人头像,没有头像发布的主题将难以引起其他成员的注意~")
				return self.Redirect("/settings/avatar/")
			}

		} else {

			if (tp.Uid == _usr_.Id) && (tp.ReplyCount > 0) {
				self.Flash.Error("话题被评论过,已被锁定不能修改!")
				return self.Redirect("/topic/" + strconv.FormatInt(self.Get("curpid").(int64), 10) + "/")
			} else {

				self.Flash.Error("你没有权限修改本话题!")
				return self.Redirect("/topic/" + strconv.FormatInt(self.Get("curpid").(int64), 10) + "/")
			}

		}
	} else {
		self.Flash.Error("非法话题!")
		return self.Redirect("/")
	}

}

func PostModifyTopicHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}
	TplNames := "modify-topic"

	ctype := self.Args("ctype").MustInt64()
	tid := self.Param("tid").MustInt64()
	nid := self.Args("nodeid").MustInt64()
	cid := self.Args("cid").MustInt64()

	images := self.FormValue("images")

	policy := helper.ObjPolicy()

	tptitle := self.FormValue("title")
	if len(tptitle) > 0 {
		tptitle = policy.Sanitize(tptitle)
	}

	tpexcerpt := self.FormValue("excerpt")
	if len(tpexcerpt) > 0 {
		tpexcerpt = policy.Sanitize(tpexcerpt)
	}

	tpcontent := self.FormValue("content")
	if len(tpcontent) > 0 {
		tpcontent = policy.Sanitize(tpcontent)
	}

	nd, err := models.GetNode(nid)
	if nd != nil && cid <= 0 {
		cid = nd.Cid
	}

	if err != nil || nid == 0 {

		self.Flash.Error("节点不存在,请创建或指定正确的节点!")

		goto render
	}

	if tid > 0 { //修改话题
		//核实主题是否真实存在
		if tp, err := models.GetTopic(tid); err == nil && tp != nil {

			if tp.Pid == 0 {
				self.Set("curpid", tp.Id)
			} else {
				self.Set("curpid", tp.Pid)
			}

			if ((tp.Uid == _usr_.Id) && (tp.ReplyCount <= 0)) || (_usr_.Role < 0) { //没有回复的时候才允许修改
				//如果是作者或管理员则可修改

				tp.Title = tptitle
				tp.Content = tpcontent
				tp.Excerpt = tpexcerpt
				tp.Cid = cid
				tp.Nid = nid
				tp.Ctype = ctype
				tp.Node = nd.Title
				//tp.Author = _usr_.Username
				tp.Created = time.Now().Unix()
				if images != "" {
					tp.Attachment = images
				}

				if tp.Pid == 0 { //如果是主题
					if tptitle != "" {
						if tid, err := models.PutTopic(tid, tp); err == nil {
							self.Flash.Success(fmt.Sprint("主题 #", tid, " [", tptitle, "] 修改成功!"))
						} else {
							self.Flash.Error(fmt.Sprint("主题 #", tid, " [", tptitle, "] 修改失败!"))
						}
						return self.Redirect("/topic/" + strconv.FormatInt(tid, 10) + "/")

					} else {

						self.Flash.Error("主题标题不能为空!")
						goto render

					}
				} else { //如果是附言
					if len(tpcontent) > 0 {
						if newtopicid, err := models.PutTopic(tid, tp); err == nil && newtopicid > 0 {

							self.Flash.Success(fmt.Sprint("附言 #", newtopicid, " [", tptitle, "] 保存成功!"))
							//return self.Redirect("/subject/" + strconv.FormatInt(tp.Id, 10) + "/topic/")
							return self.Redirect("/topic/" + strconv.FormatInt(tid, 10) + "/")

						} else {

							self.Flash.Error(fmt.Sprint("附言内容写入数据库发生错误,", err.Error()))
							goto render

						}
					} else {

						self.Flash.Error("附言内容不能为空!")
						goto render
					}
				}
			} else {

				if (tp.Uid == _usr_.Id) && (tp.ReplyCount > 0) {
					self.Flash.Error("话题含有评论被锁定不能修改!")
					return self.Redirect("/topic/" + strconv.FormatInt(self.Get("curpid").(int64), 10) + "/")

				} else {

					self.Flash.Error("你没有权限修改本话题!")
					return self.Redirect("/topic/" + strconv.FormatInt(self.Get("curpid").(int64), 10) + "/")
				}

			}

		} else { //由于可能存在管理员删话题,而用户又在继续修改,但提交新的修改时却已经被删,所以允许用户提交的话题转成新增话题继续递交

			return self.Redirect("/new/topic/") //使用307状态继续传POST
		}
	} else {
		self.Flash.Error("非法话题!")
		goto render
	}

render:
	return self.Render(TplNames)
}
