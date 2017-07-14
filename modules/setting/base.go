package setting

import (
	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"runtime"
	"time"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func BaseMiddler() makross.Handler {
	return func(self *makross.Context) error {
		tm := time.Now().UTC()
		self.Set("PageStartTime", tm.UnixNano())
		self.Set("Version", helper.Version)
		self.Set("SiteName", helper.SiteName)
		self.Set("SiteTitle", helper.SiteTitle)
		self.Set("requesturi", self.RequestURI())
		self.Set("gorotines", runtime.NumGoroutine())
		self.Set("golangver", runtime.Version())
		self.Set("UsersOnline", self.Session.Count())

		if user, okay := self.Session.Get("SignedUser").(*models.User); okay {
			self.Set("IsSigned", true)
			self.Set("IsRoot", user.Role == -1000)
			self.Set("SignedUser", user)

			self.Set("friends", models.GetFriendsByUidJoinUser(user.Id, 0, 0, "", "id"))
			messages, e := models.GetMessagesViaReceiver(0, 0, user.Username, "created")
			if !(e != nil) {
				self.Set("messages", *messages)
			}
		}

		if cats, err := models.GetCategoriesByNodeCount(0, 0, 0, "id"); cats != nil && err == nil {
			self.Set("categories", *cats)
		}

		if pages, err := models.GetPages(0, 0, "id"); pages != nil && err == nil {
			self.Set("pages", pages)
		}

		if links, err := models.GetLinks(0, 0, "id"); links != nil && err == nil {
			self.Set("links", links)
		}

		var categoriesc, nodesc, topicsc, usersc, ReplysCount, pageviews int
		cc := cache.Store(self)
		if !cc.IsExist("nodescount") {
			categoriesc, nodesc, topicsc, usersc, ReplysCount = models.Counts()
			/*
				cc.Set("categoriescount", categoriesc)
				cc.Set("nodescount", nodesc)
				cc.Set("topicscount", topicsc)
				cc.Set("userscount", usersc)
				cc.Set("ReplysCount", ReplysCount)
			*/
			self.Set("categoriescount", categoriesc)
			self.Set("nodescount", nodesc)
			self.Set("topicscount", topicsc)
			self.Set("userscount", usersc)
			self.Set("ReplysCount", ReplysCount)

		} else {

			cc.Get("categoriescount", &categoriesc)
			self.Set("categoriescount", categoriesc)

			cc.Get("nodescount", &nodesc)
			self.Set("nodescount", nodesc)

			cc.Get("topicscount", &topicsc)
			self.Set("topicscount", topicsc)

			cc.Get("userscount", &usersc)
			self.Set("userscount", usersc)

			cc.Get("ReplysCount", &ReplysCount)
			self.Set("ReplysCount", ReplysCount)

		}

		if cc.Get("pageviews", &pageviews); pageviews == 0 {
			pageviews := int64(1)
			cc.Set("pageviews", pageviews, 60*60)
			self.Set("pageviews", pageviews)

		} else {
			pageviews = pageviews + 1
			cc.Set("pageviews", pageviews, 60*60)
			self.Set("pageviews", pageviews)
		}

		//模板函数
		self.Set("TimeSince", helper.TimeSince)
		self.Set("Split", helper.Split)
		self.Set("Metric", helper.Metric)
		self.Set("Htm2Str", helper.Htm2Str)
		self.Set("Markdown", helper.Markdown)
		self.Set("Markdown2Text", helper.Markdown2Text)
		self.Set("ConvertToBase64", helper.ConvertToBase64)
		self.Set("Unix2Time", helper.Unix2Time)
		self.Set("Compare", helper.Compare)
		self.Set("TimeConsuming", func(start int64) float64 { // 0.001 s
			return float64(time.Now().UnixNano()-start) / 1000000 * 0.001
		})

		self.Set("Text", func(content string, start, length int) string {
			return helper.Substr(helper.HTML2str(content), start, length, "...")
		})

		//self.Set("Split", helper.SplitByPongo2)
		self.Set("Cropword", helper.Substr)
		self.Set("File", helper.File)
		self.Set("GetNodesByCid", func(cid int64, offset int, limit int, field string) *[]*models.Node {
			x, _ := models.GetNodesByCid(cid, offset, limit, field)
			return x
		})
		return self.Next()
	}
}
