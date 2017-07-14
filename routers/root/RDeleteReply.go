package root

import (
	"fmt"
	"github.com/insionng/makross"
	
	"strconv"
	"time"
	"github.com/insionng/yougam/models"
)

func GetRDeleteReplyHandler(self *makross.Context) error {
	

	if rid := self.Param("rid").MustInt64(); rid > 0 {

		if e := models.DelReply(rid); e != nil {
			self.Flash.Error("删除 Reply id:" + strconv.FormatInt(rid, 10) + "出现错误 " + fmt.Sprintf("%s", e) + "!")
			return e
		} else {
			self.Flash.Success("删除 Reply id:" + strconv.FormatInt(rid, 10) + "成功!")
			return self.Redirect("/root/read/reply/")

		}
	}

	return self.Redirect("/root/dashboard/?version=" + strconv.FormatInt(time.Now().Unix(), 10))

}
