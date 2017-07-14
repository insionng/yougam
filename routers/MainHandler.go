package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"strconv"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
	"github.com/insionng/yougam/modules/setting"
)

func GetMainHandler(self *makross.Context) error {

	var IsSignin bool
	//var _usr_ = new(models.User)
	if _, okay := self.Session.Get("SignedUser").(*models.User); okay {
		//_usr_ = sUser
		IsSignin = okay
	}

	cc := cache.Store(self)

	self.Set("catpage", "home")
	self.Set("messager", false)

	if v := self.Args("version").String(); len(v) > 0 {
		self.Response.Header().Set("Cache-Control", "no-cache")
	}

	keyword := self.Args("keyword").String()
	category := self.Args("category").String()
	node := self.Args("node").String()
	username := self.Args("username").String()
	cid := self.Args("cid").MustInt64()
	nid := self.Args("nid").MustInt64()
	tab := self.Args("tab").String()
	page := self.Args("page").MustInt64()
	ctype := self.Args("ctype").MustInt64()
	limit := self.Args("limit").MustInt64()
	if limit <= 0 {
		limit = 25
	}

	url := "/topics"
	if tab == "lastest" {
		url = url + "/lastest/"
		tab = "id"
		self.Set("tab", "lastest")
	} else if tab == "hotness" {
		url = url + "/hotness/"
		tab = "hotness"
		self.Set("tab", "hotness")
	} else if tab == "rising" {
		url = url + "/rising/"
		tab = "hotup"
		self.Set("tab", "rising")
	} else if tab == "scores" {
		url = url + "/scores/"
		tab = "hotscore"
		self.Set("tab", "scores")
	} else if tab == "votes" {
		url = url + "/votes/"
		tab = "hotvote"
		self.Set("tab", "votes")
	} else if tab == "controversial" {
		url = url + "/controversial/"
		tab = "reply_count"
		self.Set("tab", "controversial")
	} else if tab == "popular" {
		url = url + "/popular/"
		tab = "views"
		self.Set("tab", "popular")
	} else if tab == "cold" {
		url = url + "/cold/"
		tab = "cold"
		self.Set("tab", "cold")
	} else if tab == "favorites" {
		url = url + "/favorites/"
		tab = "favorite_count"
		self.Set("tab", "favorites")
	} else if tab == "lastest" {
		url = url + "/lastest/"
		tab = "id"
		self.Set("tab", "lastest")
	} else { //最后一个是默认选择
		url = url + "/optimal/"
		tab = "confidence"
		self.Set("tab", "optimal")
	}

	self.Set("isdefault", false)

	TplNames := "main"

	switch {
	case len(keyword) > 0: //搜索模式
		{
			//如果已经登录
			if IsSignin {
				limit = 30
			}

			if rc, err := models.SearchSubject(keyword, 0, 0, "id"); err == nil {

				rcs := int64(len(*rc))
				pages, pageout, beginnum, endnum, offset := helper.Pages(rcs, page, limit)

				if st, err := models.SearchSubject(keyword, int(offset), int(limit), "hotness"); err == nil {

					self.Set("topics", *st)

				}

				if k := self.FormValue("keyword"); len(k) > 0 {
					self.Set("search_keyword", k)
				} else {
					self.Set("search_keyword", keyword)
				}

				self.Set("pagesbar", helper.Pagesbar("/search/", keyword, rcs, pages, pageout, beginnum, endnum, 5))
			} else {

				self.Flash.Error(fmt.Sprintf("SearchSubject errors:%v", err))
				return self.Redirect("/")
			}
		}
	case len(category) > 0: //特定分类名下的话题
		{
			totalRecords, _ := models.GetSubjectsByCategory4Count(category, 0, 0, ctype)
			if totalRecords > 0 {
				pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, page, limit)

				if cat, e := models.GetCategoryByTitle(category); cat != nil && e == nil {
					cat.Views = cat.Views + 1
					cat.Hotup = cat.Hotup + 1
					cat.Hotness = helper.Hotness(cat.Hotup, cat.Hotdown, cat.Created)
					cat.Confidence = helper.Confidence(cat.Hotup, cat.Hotdown)
					cat.Hotvote = helper.QhotVote(cat.Hotup, cat.Hotdown)
					cat.Hotscore = helper.Score(cat.Hotup, cat.Hotdown)
					models.PutCategory(cat.Id, cat)
				}

				if tps := models.GetSubjectsByCategory(category, int(offset), int(limit), ctype, tab); *tps != nil {
					self.Set("topics", *tps)
					url = "/category/" + category + url
				}

				self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))
			}
		}
	case cid > 0: //特定分类ID下的话题
		{
			totalRecords, _ := models.GetSubjectsByCid4Count(cid, 0, 0, ctype)
			if totalRecords > 0 {
				pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, page, limit)

				if cat, e := models.GetCategory(cid); cat != nil && e == nil {
					cat.Views = cat.Views + 1
					cat.Hotup = cat.Hotup + 1
					cat.Hotness = helper.Hotness(cat.Hotup, cat.Hotdown, cat.Created)
					cat.Confidence = helper.Confidence(cat.Hotup, cat.Hotdown)
					cat.Hotvote = helper.QhotVote(cat.Hotup, cat.Hotdown)
					cat.Hotscore = helper.Score(cat.Hotup, cat.Hotdown)
					models.PutCategory(cat.Id, cat)
				}

				if tps := models.GetSubjectsByCid(cid, int(offset), int(limit), ctype, tab); *tps != nil {
					self.Set("topics", *tps)
					url = "/category/" + strconv.FormatInt(cid, 10) + url
				}

				self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))
			}
		}
	case len(node) > 0: //特定节点名下的话题
		{
			self.Set("CurNdTitle", node)

			totalRecords, _ := models.GetSubjectsByNode4Count(node, 0, 0, ctype)
			if totalRecords > 0 {
				pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, page, limit)

				if nd, e := models.GetNodeByTitle(node); nd != nil && e == nil {
					//self.Set["curnodeContent"] = nd.Content
					nd.Views = nd.Views + 1
					/*
						nd.Hotup = nd.Hotup + 1
						nd.Hotness = helper.Hotness(nd.Hotup, nd.Hotdown, nd.Created)
						nd.Confidence = helper.Confidence(nd.Hotup, nd.Hotdown)
						nd.Hotvote = helper.QhotVote(nd.Hotup, nd.Hotdown)
						nd.Hotscore = helper.Score(nd.Hotup, nd.Hotdown)
					*/
					models.PutNode(nd.Id, nd)
					self.Set("curnode", *nd)

				}

				if tps := models.GetSubjectsByNode(node, int(offset), int(limit), ctype, tab); *tps != nil {
					self.Set("topics", *tps)
					url = "/node/" + node + url

				}

				self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))
			}
		}
	case nid > 0: //特定节点ID下的话题
		{
			if nd, e := models.GetNode(nid); nd != nil && e == nil {
				self.Set("CurNdTitle", nd.Title)

				totalRecords, _ := models.GetSubjectsByNid4Count(nid, 0, 0, ctype)
				if totalRecords > 0 {
					pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, page, limit)

					if nd, e := models.GetNode(nid); nd != nil && e == nil {
						//self.Set["curnodeContent"] = nd.Content
						nd.Views = nd.Views + 1
						/*
							nd.Hotup = nd.Hotup + 1
							nd.Hotness = helper.Hotness(nd.Hotup, nd.Hotdown, nd.Created)
							nd.Confidence = helper.Confidence(nd.Hotup, nd.Hotdown)
							nd.Hotvote = helper.QhotVote(nd.Hotup, nd.Hotdown)
							nd.Hotscore = helper.Score(nd.Hotup, nd.Hotdown)
						*/
						models.PutNode(nd.Id, nd)
						self.Set("curnode", *nd)

					}

					if tps := models.GetSubjectsByNid(nid, int(offset), int(limit), ctype, tab); *tps != nil {
						self.Set("topics", *tps)
						url = "/node/" + strconv.FormatInt(nid, 10) + url
					}

					self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))
				}
			}

		}
	case len(username) > 0: //特定用户名下的话题
		{
			self.Set("CurUsrTitle", username)

			totalRecords, _ := models.GetSubjectsByUser4Count(username, 0, 0, ctype)
			if totalRecords > 0 {
				pages, page, beginnum, eusrnum, offset := helper.Pages(totalRecords, page, limit)

				//if usr, e := models.GetUserByUsername(username); usr != nil && e == nil {
				//self.Set["curusernameContent"] = usr.Content
				//usr.Views = usr.Views + 1
				/*
					usr.Hotup = usr.Hotup + 1
					usr.Hotness = helper.Hotness(usr.Hotup, usr.Hotdown, usr.Created)
					usr.Confidence = helper.Confidence(usr.Hotup, usr.Hotdown)
					usr.Hotvote = helper.QhotVote(usr.Hotup, usr.Hotdown)
					usr.Hotscore = helper.Score(usr.Hotup, usr.Hotdown)
				*/
				//models.PutUser(usr.Id, usr)
				//self.Set["curusername"] = *usr

				//}

				if tps := models.GetSubjectsByUser(username, int(offset), int(limit), ctype, tab); *tps != nil {
					self.Set("topics", *tps)
					url = "/createdby/" + username + url

				}

				self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, eusrnum, 5))
			}
		}
	default: //默认显示首页话题列表数据
		{
			self.Set("isdefault", true)
			totalRecords, err := models.GetTopicsByPid4Count(0, 0, 0, ctype)
			if err != nil {
				return err
			}
			if totalRecords <= 0 {
				self.Set("topics", nil)
			} else {

				pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, page, limit)

				tps := new([]*models.Topic)

				if tab == "hotup" { //计算一月内的趋势 若果没有数据则继续缩小时间范围
					tps = models.GetTopicsByPidSinceCreated(0, int(offset), int(limit), ctype, tab, helper.ThisMonth())
					if totalRecords := int64(len(*tps)); totalRecords <= 0 {
						tps = models.GetTopicsByPidSinceCreated(0, int(offset), int(limit), ctype, tab, helper.ThisWeek())
						if totalRecords := int64(len(*tps)); totalRecords <= 0 {
							tps = models.GetTopicsByPidSinceCreated(0, int(offset), int(limit), ctype, tab, helper.ThisDate())
							if totalRecords := int64(len(*tps)); totalRecords <= 0 {
								tps = models.GetTopicsByPidSinceCreated(0, int(offset), int(limit), ctype, tab, helper.ThisHour())
								if totalRecords := int64(len(*tps)); totalRecords > 0 {
									goto rising
								}
							} else {
								goto rising
							}
						} else {
							goto rising
						}
					} else {
						goto rising
					}

				rising:
					{
						pages, page, beginnum, endnum, _ := helper.Pages(totalRecords, page, limit)
						self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))
					}
				} else {

					tps = models.GetTopicsByPid(0, int(offset), int(limit), ctype, tab)
					if totalRecords := int64(len(*tps)); totalRecords > 0 {
						self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))
					}

				}
				self.Set("topics", *tps)

			}
		}
	}

	//置顶话题
	/*
		if tpsBySort := models.GetTopicsByPid(0, 0, 0, 0, "sort"); len(*tpsBySort) > 0 {
			self.Set["TopicsBySort"] = *tpsBySort
		}
	*/
	if item := setting.Cache.Get("Main_TopicsBySort"); item != nil {
		if !item.Expired() {
			if tpsBySort, okay := item.Value().(*[]*models.Topic); okay {
				self.Set("TopicsBySort", *tpsBySort)
			}
		} else {
			if tpsBySort := models.GetTopicsByPid(0, 0, 0, 0, "sort"); len(*tpsBySort) > 0 {
				cc.Set("Main_TopicsBySort", tpsBySort, 60*90)
				self.Set("TopicsBySort", *tpsBySort)
			}
		}
	} else {
		if tpsBySort := models.GetTopicsByPid(0, 0, 0, 0, "sort"); len(*tpsBySort) > 0 {
			cc.Set("Main_TopicsBySort", tpsBySort, 60*90)
			self.Set("TopicsBySort", *tpsBySort)
		}
	}

	if nds, err := models.NodesOfNavor(0, 0, "hotness"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	if nds, err := models.GetNodesByCid(0, 0, 10, "confidence"); nds != nil && err == nil {
		self.Set("nodes_sidebar_confidence_10", *nds)
	}

	/*
		if tab == "hotness" { //当主列表显示最热话题的时候 右侧显示最新话题
			if tps := models.GetTopicsByPid(0, 0, 10, 0, "created"); tps != nil {
				self.Set["topics_sidebar_10"] = *tps
			}
		} else {
			if tps := models.GetTopicsByPid(0, 0, 10, 0, "hotness"); tps != nil {
				self.Set["topics_sidebar_10"] = *tps
			}
		}
	*/
	topics_sidebar_10_key, topics_sidebar_10_tab := "", ""
	if tab == "hotness" { //当主列表显示最热话题的时候 右侧显示最新话题
		topics_sidebar_10_tab = "created"
		topics_sidebar_10_key = "topics_sidebar_10_" + "created"
	} else {
		topics_sidebar_10_tab = "hotness"
		topics_sidebar_10_key = "topics_sidebar_10_" + "hotness"
	}

	var tps []*models.Topic
	if err := cc.Get(topics_sidebar_10_key, &tps); err != nil {
		if tps := models.GetTopicsByPid(0, 0, 10, 0, topics_sidebar_10_tab); tps != nil {
			cc.Set(topics_sidebar_10_key, tps, 60*90)
			self.Set("topics_sidebar_10", *tps)
		}

	} else {

		self.Set("topics_sidebar_10", &tps)

	}

	/*
		if rps := models.GetReplysByTid(0, 0, 0, 5, "confidence"); rps != nil {
			self.Set["ConfidenceReplys"] = *rps
		}
	*/
	var rps *[]*models.Reply
	if err := cc.Get("Main_ConfidenceReplys", &rps); err != nil {

		if rps := models.GetReplysByTid(0, 0, 0, 5, "confidence"); rps != nil {
			cc.Set("Main_ConfidenceReplys", rps, 60*90)
			self.Set("ConfidenceReplys", *rps)
		}

	} else {

		self.Set("ConfidenceReplys", *rps)

	}

	/*
		if rps := models.GetReplysByTid(0, 0, 0, 5, "created"); rps != nil {
			self.Set["replys"] = *rps
		}
	*/
	if err := cc.Get("Main_replys", &rps); err != nil {
		if rps := models.GetReplysByTid(0, 0, 0, 5, "created"); rps != nil {
			cc.Set("Main_replys", rps, 60*90)
			self.Set("replys", *rps)
		}
	} else {

		self.Set("replys", *rps)

	}

	return self.Render(TplNames)

}
