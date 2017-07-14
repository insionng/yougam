package routers

import (
	"errors"
	"fmt"

	"github.com/insionng/makross"
	//"github.com/insionng/makross/cache"

	"image"
	"image/color"
	"image/png"

	"github.com/insionng/yougam/libraries/identicon"
	"github.com/insionng/yougam/models"
)

func GetIdenticonHandler(self *makross.Context) error {

	size := self.Param("size").MustInt()
	if size <= 0 {
		size = self.Args("size").MustInt()
	}
	if size <= 0 {
		return errors.New("cant used zero in size")
	}

	username := self.Param("name").String()
	if len(username) <= 0 {
		username = self.FormValue("name")
	}
	if len(username) <= 0 {
		return errors.New("cant used nil in name")
	}

	var img image.Image
	key := fmt.Sprintf("identicon:%v:%v", username, size)
	//cc := cache.Store(self)
	var okay bool
	if img, okay = self.Session.Get(key).(image.Image); !okay {

		if usr, e := models.GetUserByUsername(username); (e != nil) || (usr == nil) {
			return e
		}

		ident, e := identicon.New(size, color.RGBA{255, 0, 0, 100}, color.RGBA{0, 255, 255, 100}, color.NRGBA{}, color.NRGBA{})
		img = ident.Make([]byte(username))
		if e != nil {
			return e
		}

		//cc.Set(key, img, 60*60*60)
		self.Session.Set(key, img)

	}

	self.Response.Header().Set("Content-Type", "image/png")
	return png.Encode(self.Response, img)
}
