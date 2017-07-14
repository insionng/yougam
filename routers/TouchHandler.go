package routers

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"strconv"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetTouchHandler(self *makross.Context) error {

	if helper.IsSpider(self.UserAgent()) != true {

		_usr_, _ := self.Session.Get("SignedUser").(*models.User)
		cc := cache.Store(self)

		name := self.Param("name").String()
		id := self.Param("id").MustInt64()
		uid := _usr_.Id
		switch type_ := self.Param("type").String(); {
		case type_ == "like":
			{
				/*
					if name == "question" {
						if models.IsQuestionMark(uid, id) {
							self.Ctx.Output.SetStatus(304)
							return

						} else {
							if qs, err := models.GetQuestion(id); err == nil {

								qs.Hotup = qs.Hotup + 1
								qs.Hotscore = helper.Qhot_QScore(qs.Hotup, qs.Hotdown)
								qs.Hotvote = helper.Qhot_Vote(qs.Hotup, qs.Hotdown)
								qs.Hotness =go:89: cannot use strco/githunv.FormatInt(usr.Hotscore, 10) (type st helper.Qhot(qs.Views, qs.ReplyCount, qs.Hotscore, models.GetAScoresByPid(id), qs.Created, qs.ReplyTime)

								if _, err := models.PutQuestion(id, qs); err != nil {
									fmt.Println("PutQuestion执行错误:", err)
								} else {
									models.SetQuestionMark(uid, id)
								}
								//&hearts; 有用 ({{.article.Hotup}})
								self.Ctx.WriteString(strconv.Itoa(int(qs.Hotscore)))
							} else {
								return
							}
						}
					} else if name == "answer" {
						if models.IsAnswerMark(uid, id) {
							//self.Abort("304")
							self.Ctx.Output.SetStatus(304)
							return

						} else {
							if ans, err := models.GetAnswer(id); err == nil {

								ans.Hotup = ans.Hotup + 1
								ans.Views = ans.Views + 1
								ans.Hotscore = helper.Qhot_AScore(ans.Hotup, ans.Hotdown)
								ans.Hotvote = helper.Qhot_Vote(ans.Hotup, ans.Hotdown)
								ans.Hotness = helper.Qhot(ans.Views, ans.ReplyCount, ans.Hotscore, ans.Views, ans.Created, ans.ReplyTime)

								if _, err := models.PutAnswer(id, ans); err != nil {
									fmt.Println("PutAnswer执行错误:", err)
								} else {
									models.SetAnswerMark(uid, id)
								}
								self.Ctx.WriteString(strconv.Itoa(int(ans.Hotscore)))
							} else {
								return
							}
						}
					} else */
				if name == "user" {
					if models.IsUserMark(uid, id) {
						return self.NoContent(304)

					} else {
						if usr, err := models.GetUser(id); err == nil {
							usr.Hotup = usr.Hotup + 1
							usr.Hotscore = helper.Score(usr.Hotup, usr.Hotdown)
							usr.Hotvote = helper.QhotVote(usr.Hotup, usr.Hotdown)
							usr.Hotness = helper.Hotness(usr.Hotup, usr.Hotdown, usr.Created)
							usr.Confidence = helper.Confidence(usr.Hotup, usr.Hotdown)
							models.PutUser(id, usr)
							models.SetUserMark(uid, id)
							return self.String(strconv.FormatInt(usr.Hotscore, 10))
						} else {
							return self.NoContent(200)
						}
					}
				} else if name == "topic" {
					if models.IsTopicMark(uid, id) {
						return self.NoContent(304)

					} else {
						if tp, err := models.GetTopic(id); err == nil {

							tp.Hotup = tp.Hotup + 1
							tp.Hotscore = helper.Score(tp.Hotup, tp.Hotdown)
							tp.Hotvote = helper.QhotVote(tp.Hotup, tp.Hotdown)
							tp.Hotness = helper.Hotness(tp.Hotup, tp.Hotdown, tp.Created)
							tp.Confidence = helper.Confidence(tp.Hotup, tp.Hotdown)

							if tp.Uid != _usr_.Id { //如果不是自己点自己，那么...
								//ctype==4 为了鼓励用户投票，系统付出金币，操作用户不用付费
								if e := models.SetAmountByUid(tp.Uid, 4, +1, fmt.Sprintf("话题 <a href=\"/topic/%v/\">[%s]</a> 被 <a href=\"/user/%s/\">%s</a> 点赞, 系统赠予你1金币！", tp.Id, tp.Title, _usr_.Username, _usr_.Username)); e == nil {
									if usr, e := models.GetUser(tp.Uid); e == nil && usr != nil {
										cc.Set(fmt.Sprintf("SignedUser:%v", usr.Id), usr, 60*60*24)
									}
								}
							}

							models.PutTopic(id, tp)
							models.SetTopicMark(uid, tp.Cid, id)
							return self.String(strconv.FormatInt(tp.Hotscore, 10))
						} else {
							return self.NoContent(200)
						}
					}
				} else if name == "reply" {
					if models.IsReplyMark(uid, id) {
						return self.NoContent(304)

					} else {
						if rp, err := models.GetReply(id); err == nil {

							rp.Hotup = rp.Hotup + 1
							rp.Hotscore = helper.Score(rp.Hotup, rp.Hotdown)
							rp.Hotvote = helper.QhotVote(rp.Hotup, rp.Hotdown)
							rp.Hotness = helper.Hotness(rp.Hotup, rp.Hotdown, rp.Created)
							rp.Confidence = helper.Confidence(rp.Hotup, rp.Hotdown)
							models.PutReply(id, rp)
							models.SetReplyMark(uid, id)
							return self.String(strconv.FormatInt(rp.Hotscore, 10))
						} else {
							return self.NoContent(304)
						}
					}
				} else if name == "node" {

					if models.IsNodeMark(uid, id) {
						return self.NoContent(304)

					} else {
						if nd, err := models.GetNode(id); err == nil {

							nd.Hotup = nd.Hotup + 1
							nd.Hotscore = helper.Score(nd.Hotup, nd.Hotdown)
							nd.Hotvote = helper.QhotVote(nd.Hotup, nd.Hotdown)
							nd.Hotness = helper.Hotness(nd.Hotup, nd.Hotdown, nd.Created)
							nd.Confidence = helper.Confidence(nd.Hotup, nd.Hotdown)
							models.PutNode(id, nd)
							models.SetNodeMark(uid, id)
							return self.String(fmt.Sprintf("%v", nd.Hotscore))
						} else {
							return self.NoContent(304)
						}
					}

				} else {
					return self.NoContent(304)
				}
			}
		case type_ == "hate":
			{
				/*
					if name == "question" {
						if models.IsQuestionMark(uid, id) {
							self.Ctx.Output.SetStatus(304)
							return

						} else {
							if qs, err := models.GetQuestion(id); err == nil {

								qs.Hotdown = qs.Hotdown + 1
								qs.Hotscore = helper.Qhot_QScore(qs.Hotup, qs.Hotdown)
								qs.Hotvote = helper.Qhot_Vote(qs.Hotup, qs.Hotdown)
								qs.Hotness = helper.Qhot(qs.Views, qs.ReplyCount, qs.Hotscore, models.GetAScoresByPid(id), qs.Created, qs.ReplyTime)

								if _, err := models.PutQuestion(id, qs); err != nil {
									fmt.Println("PutQuestion执行错误:", err)
								} else {
									models.SetQuestionMark(uid, id)
								}
								self.Ctx.WriteString(strconv.Itoa(int(qs.Hotscore)))
							} else {
								return
							}
						}
					} else if name == "answer" {
						if models.IsAnswerMark(uid, id) {
							self.Ctx.Output.SetStatus(304)
							return

						} else {
							if ans, err := models.GetAnswer(id); err == nil {

								ans.Hotdown = ans.Hotdown + 1
								ans.Views = ans.Views + 1
								ans.Hotscore = helper.Qhot_QScore(ans.Hotup, ans.Hotdown)
								ans.Hotvote = helper.Qhot_Vote(ans.Hotup, ans.Hotdown)
								ans.Hotness = helper.Qhot(ans.Views, ans.ReplyCount, ans.Hotscore, ans.Views, ans.Created, ans.ReplyTime)

								if _, err := models.PutAnswer(id, ans); err != nil {
									fmt.Println("PutAnswer执行错误:", err)
								} else {
									models.SetAnswerMark(uid, id)
								}
								self.Ctx.WriteString(strconv.Itoa(int(ans.Hotscore)))
							} else {
								return
							}
						}
					} else */
				if name == "user" {
					if models.IsUserMark(uid, id) {
						return self.NoContent(304)

					} else {
						if usr, err := models.GetUser(id); err == nil {
							usr.Hotdown = usr.Hotdown + 1
							usr.Hotscore = helper.Score(usr.Hotup, usr.Hotdown)
							usr.Hotvote = helper.QhotVote(usr.Hotup, usr.Hotdown)
							usr.Hotness = helper.Hotness(usr.Hotup, usr.Hotdown, usr.Created)
							usr.Confidence = helper.Confidence(usr.Hotup, usr.Hotdown)
							models.PutUser(id, usr)
							models.SetUserMark(uid, id)
							return self.String(fmt.Sprintf("%v", usr.Hotscore))
						} else {
							return self.NoContent(200)
						}
					}
				} else if name == "topic" {
					if models.IsTopicMark(uid, id) {
						return self.NoContent(304)

					} else {
						if tp, err := models.GetTopic(id); err == nil {

							tp.Hotdown = tp.Hotdown + 1
							tp.Hotscore = helper.Score(tp.Hotup, tp.Hotdown)
							tp.Hotvote = helper.QhotVote(tp.Hotup, tp.Hotdown)
							tp.Hotness = helper.Hotness(tp.Hotup, tp.Hotdown, tp.Created)
							tp.Confidence = helper.Confidence(tp.Hotup, tp.Hotdown)
							models.PutTopic(id, tp)
							models.SetTopicMark(uid, tp.Cid, id)
							return self.String(fmt.Sprintf("%v", tp.Hotscore))
						} else {
							return self.NoContent(200)
						}
					}
				} else if name == "reply" {
					if models.IsReplyMark(uid, id) {
						return self.NoContent(304)

					} else {
						if rp, err := models.GetReply(id); err == nil {

							rp.Hotdown = rp.Hotdown + 1
							rp.Hotscore = helper.Score(rp.Hotup, rp.Hotdown)
							rp.Hotvote = helper.QhotVote(rp.Hotup, rp.Hotdown)
							rp.Hotness = helper.Hotness(rp.Hotup, rp.Hotdown, rp.Created)
							rp.Confidence = helper.Confidence(rp.Hotup, rp.Hotdown)
							models.PutReply(id, rp)
							models.SetReplyMark(uid, id)
							return self.String(fmt.Sprintf("%v", rp.Hotscore))
						} else {
							return self.NoContent(200)
						}
					}
				} else if name == "node" {

					if models.IsNodeMark(uid, id) {
						return self.NoContent(304)

					} else {

						if nd, err := models.GetNode(id); err == nil {

							nd.Hotdown = nd.Hotdown + 1
							nd.Hotscore = helper.Score(nd.Hotup, nd.Hotdown)
							nd.Hotvote = helper.QhotVote(nd.Hotup, nd.Hotdown)
							nd.Hotness = helper.Hotness(nd.Hotup, nd.Hotdown, nd.Created)
							nd.Confidence = helper.Confidence(nd.Hotup, nd.Hotdown)
							models.PutNode(id, nd)
							models.SetNodeMark(uid, id)
							return self.String(fmt.Sprintf("%v", nd.Hotscore))
						} else {
							return self.NoContent(200)
						}
					}
				} else {
					return self.NoContent(304)
				}
			}
		case type_ == "top":
			{
				if name == "topic" {
					if tp, err := models.GetTopic(id); err == nil {
						tp.Sort = tp.Sort + 1
						//models.PutSort2Topic(id, tp.Sort)
						models.PutSort2TopicViaVersion(id, tp)
						return self.String(fmt.Sprintf("%v", tp.Sort))
					} else {
						return self.NoContent(200)
					}
				}
			}
		case type_ == "bottom":
			{
				if name == "topic" {
					if tp, err := models.GetTopic(id); err == nil {
						tp.Sort = tp.Sort - 1
						//models.PutSort2Topic(id, tp.Sort)
						models.PutSort2TopicViaVersion(id, tp)
						return self.String(fmt.Sprintf("%v", tp.Sort))
					} else {
						return self.NoContent(200)
					}
				}
			}
		default:
			{
				return self.NoContent(404)
			}
		}

	} else {
		return self.NoContent(401)
	}
	return self.NoContent(200)
}
