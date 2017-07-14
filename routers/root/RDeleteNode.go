package root

import (
	"errors"
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"time"
	"github.com/insionng/yougam/models"
)

func GetRDeleteNodeHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	if nid := self.Param("nid").MustInt64(); nid > 0 {
		if nd, e := models.GetNode(nid); e == nil {
			if e := models.DelNode(nid, _usr_.Id, _usr_.Role); e != nil {
				self.Flash.Error("删除 Node id:" + strconv.Itoa(int(nid)) + "出现错误 " + fmt.Sprintf("%s", e) + "!")
				return e
			} else {
				//删除节点成功后 就去统计有多少个同样分类id的节点,把统计出来的数目写入该分类的nodecount项
				if nd.Cid > 0 {
					if nc, e := models.GetNodesByCid(nd.Cid, 0, 0, "id"); e == nil {
						if catz, e := models.GetCategory(nd.Cid); e == nil {
							catz.NodeCount = int64(len(*nc))
							models.PutCategory(nd.Cid, catz)
						}
					}
				}

				self.Flash.Success("删除 Node id:" + strconv.Itoa(int(nid)) + "成功!")
				return self.Redirect("/root/read/node/")

			}
		} else {
			es := "Node id不存在!"
			self.Flash.Error(es)
			return errors.New(es)
		}

	}

	return self.Redirect("/root/dashboard/?version=" + strconv.FormatInt(time.Now().Unix(), 10))

}
