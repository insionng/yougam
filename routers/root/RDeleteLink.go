package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"github.com/insionng/yougam/models"
)

func GetRDeleteLinkHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	if linkid := self.Param("linkid").MustInt64(); linkid > 0 {

		if e := models.DelLink(linkid, _usr_.Id, _usr_.Role); e != nil {
			self.Flash.Error(fmt.Sprintf("删除 Link id: %v 出现错误 %v !", linkid, e))

		} else {
			self.Flash.Success(fmt.Sprintf("删除 Link id:%v成功!", linkid))
		}
	}

	return self.Redirect("/root/read/link/")

}
