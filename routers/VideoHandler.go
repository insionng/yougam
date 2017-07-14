package routers

import (
	"fmt"
	"github.com/insionng/makross"
	"github.com/insionng/yougam/helper"
)

func VideoHandler(self *makross.Context) error {

	self.Set("catpage", "video")
	TplNames := "videojs"

	vid := self.Param("vid").MustInt64()
	if vid > 0 {
		fmt.Println("-------------v------------------")
		fmt.Println(vid)
		fmt.Println("-------------v------------------")

	} else {
		return self.Redirect("/")

	}

	self.Set("VideoTags", helper.VideoTags)
	return self.Render(TplNames)

}
