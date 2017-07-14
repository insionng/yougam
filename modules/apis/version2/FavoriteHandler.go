package version2

import (
	"fmt"
	"github.com/insionng/makross"
	"github.com/insionng/makross/jwt"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetFavoriteTopicHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
		}
	}

	tid := self.Param("id").MustInt64()
	if tid <= 0 {
		tid = self.Args("id").MustInt64()
	}

	var data = map[string]interface{}{}
	if tid > 0 {
		if models.IsTopicMark(uid, tid) {
			data["IsTopicMark"] = true
			return self.JSON(data)

		} else {
			data["IsTopicMark"] = false
			return self.JSON(data)
		}
	}
	herr.Message = "没有获取到查询参数"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func PostFavoriteTopicHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
		}
	}

	tid := self.Param("id").MustInt64()
	if tid <= 0 {
		tid = self.Args("id").MustInt64()
	}
	cid := self.Args("cid").MustInt64()

	if (tid > 0) && (cid > 0) {
		if models.IsTopicMark(uid, tid) {
			herr.Message = "话题已收藏过"
			herr.Status = makross.StatusInternalServerError
			return self.JSON(herr, makross.StatusInternalServerError)

		} else {
			if objs, err := models.GetTopic(tid); err != nil {
				herr.Message = "此话题不存在，无法收藏！"
				herr.Status = makross.StatusNotFound
				return self.JSON(herr, makross.StatusNotFound)
			} else {
				objs.Hotup = objs.Hotup + 1
				objs.Hotscore = helper.Score(objs.Hotup, objs.Hotdown)
				objs.Hotvote = helper.QhotVote(objs.Hotup, objs.Hotdown)
				objs.Hotness = helper.Hotness(objs.Hotup, objs.Hotdown, objs.Created)
				objs.Confidence = helper.Confidence(objs.Hotup, objs.Hotdown)

				models.PutTopic(tid, objs)
				row, err := models.SetTopicMark(uid, cid, tid)
				if (row <= 0) || (err != nil) {
					herr.Message = "收藏话题时数据库发生错误！"
					herr.Status = makross.StatusInternalServerError
					return self.JSON(herr, makross.StatusInternalServerError)
				}
				herr.Message = "成功收藏话题！"
				herr.Status = makross.StatusOK
				return self.JSON(herr)
			}
		}
	}
	herr.Message = "没有获取到话题ID"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func DelFavoriteTopicHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
		}
	}

	tid := self.Param("id").MustInt64()
	if tid <= 0 {
		tid = self.Args("id").MustInt64()
	}

	if tid > 0 {
		herr.Message = "取消收藏话题成功"
		herr.Status = makross.StatusOK
		row, err := models.DelTopicMark(uid, tid)
		if (row <= 0) || (err != nil) {
			herr.Message = "取消收藏话题失败"
			herr.Status = makross.StatusInternalServerError
			return self.JSON(herr, makross.StatusInternalServerError)
		}
		return self.JSON(herr)
	}
	herr.Message = "没有获取到话题ID"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func GetFavoriteTopicsHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	offset := self.Args("offset").MustInt()
	limit := self.Args("limit").MustInt()
	field := self.Args("field").String()
	cid := self.Args("cid").MustInt64()

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
		}
	}

	var questions []*models.Topic
	var questionMark = new([]*models.TopicMark)
	var err error
	if cid > 0 {
		if questionMark, err = models.GetTopicMarksViaUidWithCid(offset, limit, uid, cid, field); err != nil {
			herr.Message = fmt.Sprintf("获取收藏数据发生错误:%v", err)
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
	} else {
		if questionMark, err = models.GetTopicMarksViaUid(offset, limit, uid, field); err != nil {
			herr.Message = fmt.Sprintf("获取收藏数据发生错误:%v", err)
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
	}

	for _, v := range *questionMark {
		q, _ := models.GetTopic(v.Tid)
		questions = append(questions, q)
	}

	if questions != nil {
		return self.JSON(questions)
	}

	herr.Message = "没有从数据库中获取到收藏数据"
	return self.JSON(herr, makross.StatusServiceUnavailable)

}
