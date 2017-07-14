package version2

import (
	"fmt"
	"log"

	"github.com/insionng/makross"
	"github.com/insionng/makross/jwt"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

// GetTopics 获取话题列表
func GetTopicsHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	offset := self.Args("offset").MustInt()
	page := self.Args("page").MustInt64()
	pid := self.Args("pid").MustInt64()
	limit := self.Args("limit").MustInt()
	field := self.Args("field").String()
	ctype := self.Args("ctype").MustInt64()

	if field == "lastest" {
		field = "id"
	} else if (field == "hotness") || (len(field) == 0) {
		field = "hotness"
	}

	if pid != 0 {
		if offset <= 0 {
			var resultsCount int64
			if objs := models.GetTopicsViaPid(pid, 0, limit, 0, field); objs != nil {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs := models.GetTopicsViaPid(pid, int(offset_), limit, ctype, field); objs != nil {
					return self.JSON(objs)
				} else {
					herr.Message = "获取话题数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}
			} else {
				herr.Message = "没有获取到话题数据"
				herr.Status = makross.StatusOK
				return self.JSON(herr)
			}
		} else {
			if objs := models.GetTopicsViaPid(pid, offset, int(limit), ctype, field); objs != nil {
				return self.JSON(objs)
			} else {
				herr.Message = "获取话题数据出错"
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	} else {
		if offset <= 0 {
			if resultsCount, err := models.GetTopicsCount(offset, limit); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetTopics(int(offset_), limit, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetTopics(offset, int(limit), field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	}

}

func GetTopicsByUserHandler(self *makross.Context) error {

	userid := self.Args("id").MustInt64()
	offset := self.Args("offset").MustInt64()
	page := self.Args("page").MustInt64()
	limit := self.Args("limit").MustInt64()
	ctype := self.Args("ctype").MustInt64()
	field := self.Args("field").String()

	if field == "lastest" {
		field = "id"
	} else if (field == "hotness") || (len(field) == 0) {
		field = "hotness"
	}

	if limit < 0 {
		limit = 0
	}

	if page <= 0 {
		page = 1
	}

	if userid > 0 {
		if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {
			if offset <= 0 {
				_, _, _, _, offset := helper.Pages(usrinfo.TopicCount, page, limit)
				if objs := models.GetTopicsByUid(userid, int(offset), int(limit), ctype, field); objs != nil {
					return self.JSON(objs)
				}
			} else {
				if objs := models.GetTopicsByUid(userid, int(offset), int(limit), ctype, field); objs != nil {
					return self.JSON(objs)
				}
			}

		} else {
			herr.Message = "获取用户数据出错!"
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}

	}

	herr.Message = "没有用户ID!"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

// GetTopic 获取特定话题
func GetTopicHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable
	id := self.Args("id").MustInt64()
	if id > 0 {
		obj, err := models.GetTopic(id)
		if err != nil {
			herr.Message = err.Error()
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
		return self.JSON(obj)
	}
	herr.Message = "没有获取到话题ID"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

// GetContent 获取话题
func GetContentHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable
	tid := self.Args("id").MustInt64()
	field := self.Args("field").String()
	if field == "lastest" {
		field = "id"
	} else if (field == "hotness") || (len(field) == 0) {
		field = "hotness"
	}

	if tid > 0 {
		if tp, err := models.GetTopic(tid); tp != nil && err == nil {
			tp.Views = tp.Views + 1
			//tp.Hotup = tp.Hotup + 1
			if row, e := models.PutTopic(tid, tp); e != nil {
				log.Printf("GetContent更新话题ID%v访问次数数据错误, row:%v, error:%v", tid, row, e)
			}

			if objs := models.GetTopicsByPid(tid, 0, 0, 0, field); objs != nil && (len(*objs) > 0) {
				return self.JSON(objs)
			} else {
				herr.Message = fmt.Sprintf("读取主题ID为%v的数据发生错误!", tid)
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}

		}

	}

	herr.Message = "不存在此话题ID"
	return self.JSON(herr, makross.StatusServiceUnavailable)

}

// PostContent 发布话题 或 更新话题
func PostContentHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
	}
	var isRoot bool
	if jwtIsRoot, okay := claims["IsRoot"].(bool); okay {
		isRoot = jwtIsRoot
	} else {
		herr.Message = "尚无权限"
		return self.JSON(herr, makross.StatusServiceUnavailable)
	}

	var author string
	if jwtUsername, okay := claims["Username"].(string); okay {
		author = jwtUsername
	}

	var tp models.Topic
	self.Bind(&tp)

	id := self.Args("id").MustInt64()
	if id > 0 {
		tp.Id = id
	}

	if usrinfo, err := models.GetUser(uid); (err == nil) && (usrinfo != nil) && isRoot {
		tp.Uid = uid
		tp.Author = author

		if tp.Id <= 0 {
			//全新发布
			if tid, err := models.PostTopic(&tp); err != nil || tid <= 0 {
				herr.Message = fmt.Sprintf("发布内容写入数据库时发生错误：%v", err)
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				if tp, err := models.GetTopic(tid); err == nil {
					return self.JSON(tp)
				} else {
					herr.Message = fmt.Sprintf("发布内容并获取话题内容数据出错：%v", err)
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			//对指定topicid的话题进行更新
			if tid, err := models.PutTopic(tp.Id, &tp); (err != nil) || (tid <= 0) {
				herr.Message = fmt.Sprintf("更新内容写入数据库时发生错误：%v", err)
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				if tp, err := models.GetTopic(tid); err != nil {
					herr.Message = fmt.Sprintf("更新内容并获取话题内容数据出错：%v", err)
					return self.JSON(herr, makross.StatusServiceUnavailable)
				} else {
					return self.JSON(tp)
				}

			}
		}

	} else {
		herr.Message = "获取用户数据出错!"
		return self.JSON(herr, makross.StatusServiceUnavailable)

	}
}

// 更新话题
func PutContentHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
	}
	var isRoot bool
	if jwtIsRoot, okay := claims["IsRoot"].(bool); okay {
		isRoot = jwtIsRoot
	} else {
		herr.Message = "尚无权限"
		return self.JSON(herr, makross.StatusServiceUnavailable)
	}

	var author string
	if jwtUsername, okay := claims["Username"].(string); okay {
		author = jwtUsername
	}

	var tp models.Topic
	self.Bind(&tp)

	id := self.Args("id").MustInt64()
	if id > 0 {
		tp.Id = id
	}

	if usrinfo, err := models.GetUser(uid); (err == nil) && (usrinfo != nil) && isRoot {
		tp.Uid = uid
		tp.Author = author

		if tp.Id > 0 {
			//对指定topicid的话题进行更新
			if tid, err := models.PutTopic(tp.Id, &tp); (err != nil) || (tid <= 0) {
				herr.Message = fmt.Sprintf("更新内容写入数据库时发生错误：%v", err)
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				if tp, err := models.GetTopic(tid); err != nil {
					herr.Message = fmt.Sprintf("更新内容并获取话题内容数据出错：%v", err)
					return self.JSON(herr, makross.StatusServiceUnavailable)
				} else {
					return self.JSON(tp)
				}

			}
		}
	}
	herr.Message = "获取用户数据出错!"
	return self.JSON(herr, makross.StatusServiceUnavailable)

}

func DelContentHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
	}
	var isRoot bool
	if jwtIsRoot, okay := claims["IsRoot"].(bool); okay {
		isRoot = jwtIsRoot
	} else {
		herr.Message = "尚无权限"
	}

	id := self.Args("id").MustInt64()

	if isRoot && (id > 0) {
		err := models.DelTopicsByPid(id, uid, -1000)
		if err != nil {
			return self.JSON(err, makross.StatusServiceUnavailable)
		}
		herr.Message = "删除话题成功"
		herr.Status = makross.StatusOK
		return self.JSON(herr)
	}
	herr.Message = "没有获取到话题ID"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}
