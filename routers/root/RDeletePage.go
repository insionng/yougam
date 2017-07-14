package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"time"
	"github.com/insionng/yougam/models"
)

func GetRDeletePageHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	if pageid := self.Param("pageid").MustInt64(); pageid > 0 {

		if e := models.DelPage(pageid, _usr_.Id, _usr_.Role); e != nil {
			self.Flash.Error(fmt.Sprintf("删除 Page id:", pageid, "出现错误 ", "%s", e, "!"))
			return e
		} else {
			self.Flash.Success(fmt.Sprintf("删除 Page id:", pageid, "成功!"))
			return self.Redirect("/root/read/page/")
		}
	}

	return self.Redirect("/root/dashboard/?version=" + strconv.FormatInt(time.Now().Unix(), 10))

}
