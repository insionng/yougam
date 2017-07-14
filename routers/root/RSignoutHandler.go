package root

import (
	"github.com/insionng/makross"
)

func GetRSignoutHandler(self *makross.Context) error {
	return self.Redirect("/signout/?next=/root/signin/")
}
