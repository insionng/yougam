package routers

import (
	"github.com/insionng/yougam/models"

	"github.com/insionng/makross"
)

func GetUsersHandler(self *makross.Context) error {

	TplNames := "users"
	self.Set("messager", true)
	self.Set("catpage", "UsersHandler")

	__ActiveUsers, _ := models.GetUsersOnConfidence(0, 66, "last_signin_time")
	self.Set("ActiveUsers", __ActiveUsers)

	//self.Set["UsersConfidence"], _ = models.GetUsersOnConfidence(0, 10, "signin_count")
	__UsersConfidence, _ := models.GetUsers(0, 10, "confidence")
	self.Set("UsersConfidence", __UsersConfidence)

	if tps := models.GetTopicsByPid(0, 0, 10, 0, "hotness"); tps != nil {
		self.Set("topics_sidebar_10", *tps)
	}
	if nds, err := models.NodesOfNavor(0, 0, "hotness"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	return self.Render(TplNames)

}
