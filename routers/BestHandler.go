package routers

import (
	"github.com/insionng/makross"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

func GetBestHandler(self *makross.Context) error {

	self.Set("catpage", "BestGetHandler")
	TplNames := "best-comments"
	url := "/best/comments/"

	//name := self.Args("name").String()
	ctype := self.Args("type").MustInt64()
	page := self.Args("page").MustInt64()
	limit := self.Args("limit").MustInt64()
	if limit <= 0 {
		limit = 25
	}

	self.Set("isdefault", false)

	totalRecords, err := models.GetReplysByPid4Count(0, 0, 0, ctype)
	if err != nil {
		return err
	}
	if totalRecords <= 0 {
		self.Set("ConfidenceReplys", nil)
	} else {

		pages, page, beginnum, endnum, offset := helper.Pages(totalRecords, page, limit)
		if rps := models.GetReplysByTid(0, ctype, int(offset), int(limit), "confidence"); rps != nil {
			self.Set("ConfidenceReplys", *rps)
			if totalRecords := int64(len(*rps)); totalRecords > 0 {
				self.Set("pagesbar", helper.Pagesbar(url, "", totalRecords, pages, page, beginnum, endnum, 5))
			}
		}

	}

	if cats, err := models.GetCategoriesByNodeCount(0, 0, 0, "id"); cats != nil && err == nil {
		self.Set("categories", cats)
	}

	if nds, err := models.NodesOfNavor(0, 0, "hotness"); nds != nil && err == nil {
		self.Set("nodes", *nds)
	}

	/*
		if nds, err := models.GetNodesByCid(0, 0, 10, "views"); nds != nil && err == nil {
			self.Set["nodes_sidebar_10"] = *nds
		}
	*/

	if tps := models.GetTopicsByPid(0, 0, 10, 0, "confidence"); tps != nil {
		self.Set("topics_sidebar_10", *tps)
	}

	if rpys := models.GetReplysByTid(0, 0, 0, 10, "id"); rpys != nil {
		self.Set("replys", rpys)
	}

	self.Set("GetNodesByCid", func(cid int64, offset int, limit int, field string) *[]*models.Node {
		x, _ := models.GetNodesByCid(cid, offset, limit, field)
		return x
	})

	return self.Render(TplNames)

}
