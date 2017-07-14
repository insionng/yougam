package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"time"
	"github.com/insionng/yougam/models"
)

func GetRDeleteCategoryHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	if cid := self.Param("cid").MustInt64(); cid > 0 {

		if e := models.DelCategory(cid, _usr_.Id, _usr_.Role); e != nil {
			self.Flash.Error(fmt.Sprintf("删除 Category id:", cid, "出现错误 ", "%s", e, "!"))
			return e
		} else {
			self.Flash.Success(fmt.Sprintf("删除 Category id:", cid, "成功!"))
			return self.Redirect("/root/read/category/")
		}
	}

	return self.Redirect("/root/dashboard/?version=" + strconv.FormatInt(time.Now().Unix(), 10))

}
