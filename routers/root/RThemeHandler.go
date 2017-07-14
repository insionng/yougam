package root

import (
	"encoding/base64"
	"fmt"
	"github.com/insionng/makross"
	
	"io/ioutil"
	"os"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/libraries/goconfig"
)

func GetRThemeHandler(self *makross.Context) error {
	

	self.Set("catpage", "RThemeHandler")
	self.Set("curtheme", helper.Theme())

	themes := Fstdirs("themes")
	logoz := make(map[string]string, 0)
	for _, theme := range themes {
		b, e := ioutil.ReadFile(fmt.Sprintf("themes/%s/static/%s.png", theme, theme))
		if e != nil {
			continue
		}
		if len(b) > 0 {
			logoz[theme] = base64.StdEncoding.EncodeToString(b)
		}
	}
	self.Set("logoz", logoz)

	if themename := self.Param("name").String(); len(themename) > 0 && themename != helper.Theme() {

		conf, e := goconfig.LoadConfigFile(helper.ConfigPath)
		if e != nil {
			self.Flash.Error(fmt.Sprintf("读取配置文件出错:", e))
			return self.Render("root/theme")

		} else {

			newStaticPath := "themes/" + themename + "/static"
			newTemplatesPath := "themes/" + themename + "/templates"

			if helper.IsExist(newStaticPath) && helper.IsExist(newTemplatesPath) {

				conf.SetValue("themes", "style", themename)
				goconfig.SaveConfigFile(conf, helper.ConfigPath)

				if e := os.RemoveAll("./public"); e != nil {
					self.Flash.Error(fmt.Sprintf("删除目录出错:", e))
					return self.Render("root/theme")

				}

				//复制模板静态资源到public目录
				if e := helper.CopyDir(newStaticPath, "./public"); e != nil {
					self.Flash.Error(fmt.Sprintf("复制模板静态资源到public目录时出现错误:", e))
					return self.Render("root/theme")

				}

				return self.Redirect("/?version=" + helper.MD5(themename))

			} else {
				self.Flash.Error(fmt.Sprint(themename, "主题不存在!"))
				return self.Render("root/theme")
			}
		}

	}

	return self.Render("root/theme")
}

func Fstdirs(dir string) (dirlist []string) {
	f, _ := os.Open(dir)
	defer f.Close()
	dirs, _ := f.Readdir(0)
	for _, fileInfo := range dirs {
		if fileInfo.IsDir() {
			dirlist = append(dirlist, fileInfo.Name())
		}
	}
	return dirlist
}
