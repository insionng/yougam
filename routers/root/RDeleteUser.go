package root

import (
	"fmt"

	"github.com/insionng/makross"

	"strconv"
	"strings"
	"time"

	"github.com/insionng/yougam/models"
)

func GetRDeleteUserHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	opuid := _usr_.Id
	if userid := self.Param("uid").MustInt64(); userid > 0 {
		if e := models.DelUser(userid, opuid, _usr_.Role); e != nil {
			self.Flash.Error("删除 User id:"+strconv.FormatInt(userid, 10)+"出现错误 "+fmt.Sprintf("%s", e)+"!", false)
			return self.Redirect("/root/read/user/")
		} else {

			self.Flash.Success("删除 User id:"+strconv.FormatInt(userid, 10)+"成功!", false)
			if strings.HasPrefix(self.RequestURI(), "/root/delete/administrator/") {
				return self.Redirect("/root/read/administrator/")
			} else {
				return self.Redirect("/root/read/user/")
			}

		}
	}

	return self.Redirect("/root/read/user/?version=" + strconv.Itoa(int(time.Now().Unix())))

}
