package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"github.com/insionng/yougam/models"
)

func GetTopicHandler(self *makross.Context) error {

	_usr_, _ := self.Session.Get("SignedUser").(*models.User)
	cc := cache.Store(self)

	TplNames := "topic"
	self.Set("catpage", "topic")

	if tid := self.Param("tid").MustInt64(); tid > 0 {
		if tps := models.GetTopicsByPid(tid, 0, 0, 0, "asc"); (tps != nil) && (len(*tps) > 0) {
			self.Set("topics", *tps)
			if tp, err := models.GetTopic(tid); (tps != nil) && (err == nil) {
				var allow = false
				switch kind := tp.Ctype; {
				/*
				   tp.Ctype ==  0 [普通话题]
				   tp.Ctype == -1 [回复可见]
				   tp.Ctype == -2 [付费可见]
				   tp.Ctype == -3 [会员可见]
				*/
				case kind == 0:
					allow = true //普通话题默认允许显示
				case kind == -1:
					if _usr_ != nil { //若果是登录状态
						if _usr_.Id == tp.Uid { //如果当前用户即是作者则允许
							allow = true
						}

						if _usr_.Role == -1000 {
							allow = true
						}

						if !allow {
							allow = models.IsTailinfoOfUser(tid, _usr_.Id)
						}
					}
				case kind == -2:
					if _usr_ != nil { //若果是登录状态
						if _usr_.Id == tp.Uid { //如果当前用户即是作者则允许
							allow = true
						}

						if _usr_.Role == -1000 {
							allow = true
						}

						if !allow {
							allow = models.IsTailinfoOfUser(tid, _usr_.Id) //须要另外付费 并写入标记
						}
					}
				case kind == -3:
					if _usr_ != nil {
						allow = true
					}
				}

				self.Set("Allow", allow)
				self.Set("article", *tp)
				self.Set("Excerpt", tp.Excerpt)

				__author, _ := models.GetUser(tp.Uid)
				self.Set("author", __author)

				__curnode, _ := models.GetNode(tp.Nid)
				self.Set("curnode", __curnode)

				//侧栏推荐同节点下的优选话题
				/*
					if tps := models.GetSubjectsByNid(tp.Nid, 0, 10, 0, "confidence"); *tps != nil {
						self.Set["topic_sidebar_hotness_10"] = *tps
					}
				*/

				var tps *[]*models.Topic

				if err := cc.Get(fmt.Sprintf("topic_sidebar_hotness_10_%v", tp.Nid), &tps); err != nil {
					if tps = models.GetSubjectsByNid(tp.Nid, 0, 10, 0, "confidence"); *tps != nil {
						cc.Set(fmt.Sprintf("topic_sidebar_hotness_10_%v", tp.Nid), tps, 60*90)
						self.Set("topic_sidebar_hotness_10", *tps)
					}

				} else {
					self.Set("topic_sidebar_hotness_10", *tps)
				}

				//推荐同一作者的优选话题
				/*
					if tps := models.GetSubjectsByUid(tp.Uid, 0, 10, 0, "confidence"); len(*tps) > 0 {
						self.Set["TopicsByUser_10s"] = *tps
					}
				*/

				if err := cc.Get(fmt.Sprintf("TopicsByUser_10s_%v", tp.Uid), &tps); err != nil {
					if tps := models.GetSubjectsByUid(tp.Uid, 0, 10, 0, "confidence"); len(*tps) > 0 {
						cc.Set(fmt.Sprintf("TopicsByUser_10s_%v", tp.Uid), tps, 60*90)
						self.Set("TopicsByUser_10s", *tps)
					}
				} else {
					self.Set("TopicsByUser_10s", *tps)
				}

			}

			if rps := models.GetReplysByTid(tid, 0, 0, 0, "confidence"); (rps != nil) && (len(*rps) > 0) {
				self.Set("replys", *rps)
			}

			return self.Render(TplNames)

		} else {
			self.Flash.Error("话题并不存在！")
			return self.Redirect("/")

		}

	} else {
		self.Flash.Error("非法参数")
		return self.Redirect("/")

	}

}
