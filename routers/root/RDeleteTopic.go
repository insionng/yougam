package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"time"
	"github.com/insionng/yougam/models"
)

func GetRDeleteTopicHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	uid := _usr_.Id
	role := _usr_.Role
	if tid := self.Param("tid").MustInt64(); tid > 0 {

		if tp, err := models.GetTopic(tid); tp != nil && err == nil {
			if tp.Pid == 0 { //如果是主题 就删光下属所有话题

				if tps := models.GetTopicsByPid(tid, 0, 0, 0, "id"); tps != nil {
					for _, g := range *tps {
						models.DelTopic(g.Id, uid, role) //该操作是在后台进行 所以设为当前管理员uid
						models.DelTopicMark(g.Uid, g.Id) //删除该用户的话题收藏 所以这里是g.Uid
					}

					//删除下属评论
					models.DelReplysByTid(tid)

					self.Flash.Success("删除 Topic id:" + strconv.FormatInt(tid, 10) + "成功!")
					return self.Redirect("/root/read/topic/")

				}

			} else { //pid不等于0是子话题 此时只需要删除该话题 上级和其他子话题都不用删除
				if e := models.DelTopic(tid, uid, role); e != nil {
					self.Flash.Error("删除 Topic id:" + strconv.FormatInt(tid, 10) + "出现错误 " + fmt.Sprintf("%s", e) + "!")
					return e
				} else {

					models.DelTopicMark(tp.Uid, tid) //删除该用户的话题收藏 所以这里是tp.Uid

					//删除下属评论
					models.DelReplysByTid(tp.Pid)

					self.Flash.Success("删除 Topic id:" + strconv.FormatInt(tid, 10) + "成功!")
					return self.Redirect("/root/read/topic/")

				}

			}
		} else {
			//读取数据异常后返回
			self.Flash.Error("删除 Topic id:" + strconv.FormatInt(tid, 10) + "出现错误 " + err.Error() + "!")
			return self.Redirect("/root/read/topic/")

		}
	}

	return self.Redirect("/root/dashboard/?version=" + strconv.FormatInt(time.Now().Unix(), 10))

}
