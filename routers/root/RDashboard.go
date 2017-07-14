package root

import (
	"github.com/insionng/makross"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetRDashboardHandler(self *makross.Context) error {

	self.Set("catpage", "RDashboardHandler")
	TplNames := "root/dashboard"

	if tps, err := models.GetTopics(0, 3, "id"); tps != nil && err == nil {
		self.Set("tps", *tps)
	}

	if rpys := models.GetReplysByTid(0, 0, 0, 10, "id"); rpys != nil {
		self.Set("replys", rpys)
	}

	if tps := models.GetTopicsByPid(0, 0, 10, 0, "hotness"); tps != nil {
		self.Set("topics_sidebar_10", *tps)
	}

	//本周趋势
	if tps := models.GetTopicsByPidSinceCreated(0, 0, 10, 0, "hotup", helper.ThisWeek()); len(*tps) > 0 {
		self.Set("topics_thisWeek_10", *tps)
	}
	return self.Render(TplNames)
}
