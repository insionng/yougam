package routers

import (
	"fmt"
	"github.com/insionng/makross"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetTouchViewHandler(self *makross.Context) error {

	name := self.Param("name").String()
	id := self.Param("id").MustInt64()

	if (len(name) > 0) && (id > 0) && (!helper.IsSpider(self.UserAgent())) {
		if name == "topic" {

			if tp, err := models.GetTopic(id); tp != nil && err == nil {
				tp.Views = tp.Views + 1
				if row, e := models.PutViews2TopicViaVersion(id, tp); (e == nil) && (row > 0) {
					return self.String(fmt.Sprintf("%v", tp.Views))
				} else {
					return self.String(fmt.Sprintf("%v", tp.Views-1))
				}
			}
		}

	}
	return self.NoContent(makross.StatusOK)
}
