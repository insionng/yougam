package root

import (
	"github.com/insionng/makross"
	"github.com/insionng/yougam/models"
)

func GetRReadIssueHandler(self *makross.Context) error {

	TplNames := ""
	self.Set("catpage", "RReadIssueHandler")
	switch rid := self.Param("rid").MustInt64(); {
	//单独模式
	case rid > 0:
		{
			TplNames = "root/read_issue.html"

			if rm, err := models.GetIssue(rid); rm != nil && err == nil {
				self.Set("issue", *rm)

				if Issues, err := models.GetIssues(0, 0, "id", 0); Issues != nil && err == nil {
					self.Set("issues", Issues)
				}

			} else {
				self.Flash.Error(err.Error())
				return self.Render(TplNames)

			}
		}
	//列表模式
	case rid <= 0:
		{

			TplNames = "root/issue_table"
			offset := self.Args("offset").MustInt64()
			limit := self.Args("limit").MustInt64()
			field := self.FormValue("field")
			ctype := self.Args("ctype").MustInt64()

			if limit == 0 {
				limit = 1000 //默认限制显示最近1000条,需要显示全部请在提交请求的时候设置limit字段为-1
			}

			if field == "" {
				field = "id"
			}

			if ctype != -1 || ctype != 1 {

				switch name := self.Param("name").String(); {
				case name == "topic":
					self.Set("sidebar", "rtopics")
					ctype = 1
				case name == "reply":
					self.Set("sidebar", "rreplys")
					ctype = -1
				default:
					self.Set("sidebar", "issues")
					ctype = 0
				}

			}

			if rms, err := models.GetIssues(int(offset), int(limit), field, ctype); err == nil && rms != nil {
				self.Set("issues", rms)
			}
		}
	}
	return self.Render(TplNames)

}
