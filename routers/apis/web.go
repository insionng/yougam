package apis

import (
	"github.com/insionng/makross"

	"time"
	"github.com/insionng/yougam/models"
)

func GetMessagesHandler(self *makross.Context) error {

	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			return self.JSON(nil, 204)
		default:
			if messages, e := models.GetMessagesViaReceiver(0, 0, _usr_.Username, "created"); (!(e != nil)) && (len(*messages) > 0) {
				return self.JSON(messages)

			} else {
				time.Sleep(time.Second * 3)
			}
		}
	}

}
