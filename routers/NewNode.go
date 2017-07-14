package routers

import (
	"fmt"
	"github.com/insionng/makross"
	
	"math"
	"strconv"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
	"github.com/insionng/yougam/modules/setting"
)

func GetNewNodeHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	TplNames := "new-node"
	if len(_usr_.Avatar) > 0 {
		self.Set("hasavatar", true)
		allow := bool(false)
		if _usr_.Balance >= int64(math.Abs(setting.CreateNodesOfGoldCoins)) {
			allow = true
		} else {
			self.Flash.Error(fmt.Sprintf("你的余额为[%v]，不足以创建节点!", _usr_.Balance))
			return self.Redirect(fmt.Sprintf("/user/%v/balance/", _usr_.Username))
		}

		if catz, e := models.GetCategories(0, 0, "id"); (catz != nil) && (e == nil) && allow {

			self.Set("categories", catz)
			self.Set("hascid", false)

			//匹配路由中的:category string
			if category := self.Param("category").String(); category != "" {
				if cat, e := models.GetCategoryByTitle(category); cat != nil && e == nil {
					self.Set("hascid", true)
					self.Set("curcid", cat.Id)
				}
			} else {
				//匹配路由中的:cid int64
				if cid := self.Param("cid").MustInt64(); cid > 0 {
					if cat, e := models.GetCategory(cid); cat != nil && e == nil {
						self.Set("hascid", true)
						self.Set("curcid", cid)
					}
				}

				//当路由都没有带参数的时候 pass

			}

			self.Set("haspid", false)
			self.Set("catpage", "NewNodeHandler")

			nid := self.Param("nid").MustInt64()
			if nid > 0 {
				_nd, _ := models.GetNode(nid)
				self.Set("node", _nd)
				self.Set("haspid", true)
				self.Set("curpid", nid)

				/*
					if nds := models.GetNodesByPid(nid, 0, 0, 0, "id"); nds != nil && (len(*nds) > 0) {
						self.Set["nodes"] = *nds

					} else {
						return self.Redirect("/category/"+strconv.FormatInt(nds.Cid, 10)+"/", 302)
						return
					}
				*/

			}

			nds, _ := models.AvailableNodes(0, 0, "id")
			self.Set("nodes", *nds)
			return self.Render(TplNames)
		} else {
			//不存在分类则跳转到首页,管理员需要创建至少一个默认分类
			self.Flash.Error("请管理员创建至少一个默认分类~")
			return self.Redirect("/")
		}
	} else {
		self.Set("hasavatar", false)
		self.Flash.Warning("请设置你的头像,有头像能更好的引起其他成员的注意~")
		return self.Redirect("/settings/avatar/")
	}

}

/*
func NewNodePostHandler(self *makross.Context) error {

	TplNames := "new-node"
	_usr_, _ := self.Session.Get("user").(*models.User)
	if _usr_ == nil {
		return
	}

	cid := com.StrToself.Param("cid")
	category := self.Params("category")
	catid := com.StrToself.FormValue("catid")
	uid, _ := self.Session.Get("SignedUser").(*models.User)

	images := self.FormValue("images")
	ndtitle := self.FormValue("title")

	policy := helper.ObjPolicy()
	ndcontent := policy.Sanitize(self.Query("content"))

	fmt.Println(category)
	if cid <= 0 {

		cat, err := models.GetCategory(catid)

		cid = int64(0)

		if cat != nil {
			cid = cat.Id
		}

		if err != nil || cid == 0 {

			self.Flash.Error("分类不存在,请创建或指定正确的分类!")

			goto render
			return
		}

		if ndtitle != "" {

			nd := new(models.Node)
			nd.Title = ndtitle
			nd.Content = ndcontent
			nd.Cid = catid
			nd.Pid = 0
			nd.Uid = _usr_.Id
			nd.Author = _usr_.Username
			nd.Created = time.Now().Unix()

			if images != "" {
				nd.Attachment = images
			} ///else {
				///if s, e := helper.GetBannerThumbnail(tpcontent); e == nil {
				///	nd.Attachment = s
				///}
			///}


			if cat, err := models.GetCategory(catid); err == nil && cat != nil {
				nd.Parentname = cat.Title
			}

			catzmap := map[string]interface{}{
				"node_time":         time.Now().Unix(),
				"node_count":        models.GetNodeCountByPid(catid),
				"node_last_user_id": uid}

			if e := models.UpdateCategory(catid, &catzmap); e != nil {
				log.Println("NewNodePostHandler models.UpdateCategory errors:", e)
			}

			if nid_, err := models.PutNode(0, nd); err == nil {
				//models.SetRecordforImageOnPost(nid_, uid)

				self.Set["curpid"] = nid_
				self.Flash.Success(fmt.Sprint("节点#", nid_, "保存成功!"), false)

				//return self.Redirect("/subject/"+strconv.FormatInt(nid_, 10)+"/category/", 302)
				return self.Redirect("/category/"+strconv.FormatInt(nid_, 10)+"/", 302)
				return

			} else {
				self.Set["curpid"] = nid_
				self.Flash.Error(err.Error())

				goto render
				return
			}
		} else {
			if ndtitle == "" {
				self.Flash.Error("节点标题为空!")
			}

			goto render
			return
		}
	} else {
		//核实分类是否真实存在
		if cat, err := models.GetCategory(cid); err == nil && cat != nil {
			self.Set["curpid"] = cat.Id
			if ndcontent != "" {
				if NewNid, err := models.AddNode(ndtitle, ndcontent, images, 0, cat.Id, _usr_.Id); err == nil && NewNid > 0 {

					self.Flash.Success(fmt.Sprint("节点#", NewNid, "保存成功!"), false)
					return self.Redirect(fmt.Sprintf("/user/%v/balance/", _usr_.Username), 302)
					return
				} else {

					self.Flash.Error(fmt.Sprint("节点内容写入数据库发生错误,", err.Error()))

					goto render
					return
				}
			} else {

				self.Flash.Error("节点内容不能为空!")

				goto render
				return
			}
		} else {
			self.Flash.Error("非法分类!")

			goto render
			return
		}
	}

render:
	return self.Render(TplNames)
}
*/

func PostNewNodeHandler(self *makross.Context) error {

	
	_usr_, okay := self.Session.Get("SignedUser").(*models.User)
	if !okay {
		return self.NoContent(401)
	}

	policy := helper.ObjPolicy()
	content := policy.Sanitize(self.FormValue("content"))
	pid := self.Args("nodeid").MustInt64()
	title := self.FormValue("title")
	images := self.FormValue("images")

	cid := self.Args("cid").MustInt64()

	allow := false

	if _usr_.Balance >= int64(math.Abs(setting.CreateNodesOfGoldCoins)) {
		allow = true
	} else {
		self.Flash.Error(fmt.Sprintf("你的余额为[%v]，不足以创建节点!", _usr_.Balance))
		return self.Redirect(fmt.Sprintf("/user/%v/balance/", _usr_.Username))
	}

	if (len(title) > 0) /* && (content != "") */ && allow {

		if nid, err := models.AddNode(title, content, images, pid, cid, _usr_.Id); err != nil {

			self.Flash.Error(fmt.Sprint("增加节点出现错误:", err))
			return self.Redirect("/new/node/")
		} else {

			//新增点成功后 就去统计有多少个同样分类id的节点,把统计出来的数目写入该分类的NodeCount项
			if cid > 0 {
				if nc, e := models.GetNodesByCid(cid, 0, 0, "id"); e == nil {
					if catz, e := models.GetCategory(cid); e == nil {
						catz.NodeCount = int64(len(*nc))
						models.PutCategory(cid, catz)
					}
				}
			}

			if models.SetAmountByUid(_usr_.Id, -3, +setting.CreateNodesOfGoldCoins, fmt.Sprintf("创建节点消耗%v金币", int64(math.Abs(setting.CreateNodesOfGoldCoins)))) == nil {
				_usr_.Balance = _usr_.Balance + setting.CreateNodesOfGoldCoins
				self.Session.Set("user", _usr_)
				self.Set("user", _usr_)
			}
			self.Flash.Success("新增节点成功!")
			self.Set("curnode", title) //传到重定向的页面作为新建话题之参数
			return self.Redirect("/node/" + strconv.FormatInt(nid, 10) + "/")
		}
	} else {
		self.Flash.Error("节点最低要求标题不能为空！")
		return self.Redirect("/new/node/")
	}

}
