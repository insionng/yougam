package routers

import (
	"github.com/insionng/makross"

	"time"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetFavoriteHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}
	if helper.IsSpider(self.UserAgent()) != true {
		name := self.Param("name").String()
		id := self.Param("id").MustInt64()
		uid := _usr_.Id
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
						qs.Hotness = helper.Qhot(qs.Views, qs.ReplyCount, qs.Hotscore, models.GetAScoresByPid(id), qs.Created, qs.ReplyTime)

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
			} else*/
		if name == "topic" {
			itm := models.IsTopicMark(uid, id)

			switch {
			case itm == true: //已经被收藏 则以下取消收藏
				{

					if tp, err := models.GetTopic(id); err == nil {

						tp.Hotdown = tp.Hotdown + 1
						tp.Hotscore = helper.Score(tp.Hotup, tp.Hotdown)
						tp.Hotvote = helper.QhotVote(tp.Hotup, tp.Hotdown)
						tp.Hotness = helper.Hotness(tp.Hotup, tp.Hotdown, tp.Created)

						models.DelTopicMark(uid, id)

						//统计话题被收藏数
						tp.FavoriteCount, _ = models.TopicMarkCount(id)
						models.PutTopic(id, tp)

						if usr, e := models.GetUser(uid); e == nil && usr != nil {
							tmc, _ := models.TopicMarkCountByUid(uid)

							//用户自己收藏的话题总数
							usr.FavoriteCount = tmc
							models.PutUser(uid, usr)
						}
						//return false, strconv.FormatInt(tp.FavoriteCount, 10)
						data := map[string]interface{}{}
						data["isFavorite"] = false
						data["FavoriteCount"] = tp.FavoriteCount
						return self.JSON(data)
					}

				}
			case itm == false: //尚未被收藏 则以下进行收藏
				{

					if tp, err := models.GetTopic(id); err == nil {

						tp.Hotup = tp.Hotup + 1
						tp.Hotscore = helper.Score(tp.Hotup, tp.Hotdown)
						tp.Hotvote = helper.QhotVote(tp.Hotup, tp.Hotdown)
						tp.Hotness = helper.Hotness(tp.Hotup, tp.Hotdown, tp.Created)

						models.SetTopicMark(uid, tp.Cid, id)

						//统计话题被收藏数
						tp.FavoriteCount, _ = models.TopicMarkCount(id)
						models.PutTopic(id, tp)

						if usr, e := models.GetUser(uid); e == nil && usr != nil {
							tmc, _ := models.TopicMarkCountByUid(uid)

							//用户自己收藏的话题总数
							usr.FavoriteCount = tmc
							models.PutUser(uid, usr)
						}
						//return true, strconv.FormatInt(tp.FavoriteCount, 10)
						data := map[string]interface{}{}
						data["isFavorite"] = true
						data["FavoriteCount"] = tp.FavoriteCount
						return self.JSON(data)
					}

				}

			}
		} else if name == "node" {

			if nd, err := models.GetNode(id); err == nil {

				nd.Hotup = nd.Hotup + 1
				nd.Hotscore = helper.Score(nd.Hotup, nd.Hotdown)
				nd.Hotness = helper.Hotness(nd.Hotup, nd.Hotdown, time.Now().Unix())
				models.PutNode(id, nd)

				data := map[string]interface{}{}
				data["isFavorite"] = true
				data["FavoriteCount"] = nd.FavoriteCount
				return self.JSON(data)
			} else {

				data := map[string]interface{}{}
				data["isFavorite"] = false
				data["FavoriteCount"] = nd.FavoriteCount
				return self.JSON(data)
			}
		} else {
			return self.NoContent(304)
		}

	}

	return self.NoContent(401)

}
