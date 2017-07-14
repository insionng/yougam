package version2

import (
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
	"fmt"
	"github.com/insionng/makross"
	"github.com/insionng/makross/jwt"
)

func GetNodeHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable
	nid := self.Args("id").MustInt64()

	if nid > 0 {
		obj, err := models.GetNode(nid)
		if err != nil {
			herr.Message = err.Error()
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
		return self.JSON(obj)
	}
	herr.Message = "没有获取到节点ID"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

// GetNodes
/*
func GetNodes(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	offset := self.Args("offset").MustInt()
	page := self.Args("page").MustInt64()
	limit := self.Args("limit").MustInt()
	field := self.Args("field").String()
	pid := self.Args("pid").MustInt64()
	cid := self.Args("cid").MustInt64()
	ctype := self.Args("ctype").MustInt64()

	if field == "lastest" {
		field = "id"
	} else if (field == "hotness") || (len(field) == 0) {
		field = "hotness"
	}

	switch {
	case (ctype != 0) && (pid > 0): // 获取特定ctype节点里特定pid之列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetNodesByCtypeWithPid(0, limit, ctype, pid, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodesByCtypeWithPid(int(offset_), limit, ctype, pid, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetNodesByCtypeWithPid(offset, limit, ctype, pid, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	case (ctype == 0) && (pid > 0): // 获取特定节点之下级节点列表
		if offset <= 0 {
			var resultsCount int64
			if objs := models.GetNodesByPid(pid, 0, limit, 0, field); objs != nil {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs := models.GetNodesByPid(pid, int(offset_), limit, 0, field); objs != nil {
					return self.JSON(objs)
				} else {
					herr.Message = "获取节点数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			} else {
				herr.Message = "没有获取到节点数据"
				herr.Status = makross.StatusOK
				return self.JSON(herr)
			}
		} else {
			if objs := models.GetNodesByPid(pid, offset, limit, 0, field); objs != nil {
				return self.JSON(objs)
			} else {
				herr.Message = "没有获取到节点数据"
				herr.Status = makross.StatusOK
				return self.JSON(herr)
			}
		}
	case (ctype != 0) && (cid > 0): // 获取特定ctype节点里特定cid之列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetNodesByCtypeWithCid(0, limit, ctype, cid, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodesByCtypeWithCid(int(offset_), limit, ctype, cid, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetNodesByCtypeWithCid(offset, limit, ctype, cid, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	case (ctype == 0) && (cid > 0): // 获取特定分类之下级节点列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetNodesByCid(cid, 0, limit, field); err != nil {
				herr.Message = fmt.Sprintf("获取节点数据发生错误:%v", err)
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodesByCid(cid, int(offset_), limit, field); err != nil {
					herr.Message = "获取节点数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)
				} else {
					return self.JSON(objs)
				}
			}

		} else {
			if objs, err := models.GetNodesByCid(cid, offset, limit, field); err != nil {
				herr.Message = "获取节点数据出错"
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				return self.JSON(objs)
			}
		}
	case (ctype != 0) && (pid <= 0) && (cid <= 0): // 获取特定ctype节点列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetNodesByCtype(0, limit, ctype, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodesByCtype(int(offset_), limit, ctype, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetNodesByCtype(offset, limit, ctype, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	default: // 获取全部节点列表
		if offset <= 0 {
			if resultsCount, err := models.GetNodesCount(offset, limit); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodes(int(offset_), limit, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetNodes(offset, int(limit), field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	}
}
*/
// GetNodes
func GetNodesHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	offset := self.Args("offset").MustInt()
	page := self.Args("page").MustInt64()
	limit := self.Args("limit").MustInt()
	field := self.Args("field").String()
	pid := self.Args("pid").MustInt64()
	cid := self.Args("cid").MustInt64()
	ctype := self.Args("ctype").MustInt64()

	if field == "lastest" {
		field = "id"
	} else if (field == "hotness") || (len(field) == 0) {
		field = "hotness"
	}

	switch {
	case (ctype != 0) && (pid != 0): // 获取特定ctype节点里特定pid之节点列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetNodesByCtypeWithPid(0, limit, ctype, pid, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodesByCtypeWithPid(int(offset_), limit, ctype, pid, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetNodesByCtypeWithPid(offset, limit, ctype, pid, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	case (ctype == 0) && (pid != 0): // 获取特定节点之下级节点列表
		if offset <= 0 {
			var resultsCount int64
			if objs := models.GetNodesViaPid(pid, 0, limit, 0, field); objs != nil {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs := models.GetNodesViaPid(pid, int(offset_), limit, 0, field); objs != nil {
					return self.JSON(objs)
				} else {
					herr.Message = "获取节点数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			} else {
				herr.Message = "没有获取到节点数据"
				herr.Status = makross.StatusOK
				return self.JSON(herr)
			}
		} else {
			if objs := models.GetNodesViaPid(pid, offset, limit, 0, field); objs != nil {
				return self.JSON(objs)
			} else {
				herr.Message = "没有获取到节点数据"
				herr.Status = makross.StatusOK
				return self.JSON(herr)
			}
		}
	case (ctype != 0) && (cid > 0): // 获取特定ctype节点里特定cid之节点列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetNodesByCtypeWithCid(0, limit, ctype, cid, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodesByCtypeWithCid(int(offset_), limit, ctype, cid, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetNodesByCtypeWithCid(offset, limit, ctype, cid, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	case (ctype == 0) && (cid > 0): // 获取特定节点之下级节点列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetNodesByCid(cid, 0, limit, field); err != nil {
				herr.Message = fmt.Sprintf("获取节点数据发生错误:%v", err)
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodesByCid(cid, int(offset_), limit, field); err != nil {
					herr.Message = "获取节点数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)
				} else {
					return self.JSON(objs)
				}
			}

		} else {
			if objs, err := models.GetNodesByCid(cid, offset, limit, field); err != nil {
				herr.Message = "获取节点数据出错"
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				return self.JSON(objs)
			}
		}
	case (ctype != 0) && (pid == 0) && (cid <= 0): // 获取特定ctype节点列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetNodesByCtype(0, limit, ctype, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodesByCtype(int(offset_), limit, ctype, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetNodesByCtype(offset, limit, ctype, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	default: // 获取全部节点列表
		if offset <= 0 {
			if resultsCount, err := models.GetNodesCount(offset, limit); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetNodes(int(offset_), limit, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetNodes(offset, int(limit), field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	}
}

func PostNodeHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
		}
	}

	var author string
	if jwtUsername, okay := claims["Username"].(string); okay {
		author = jwtUsername
	}

	var obj models.Node
	self.Bind(&obj)

	if (obj.Pid >= 0) && (len(obj.Title) > 0) && (uid > 0) {
		obj.Uid = uid
		obj.Author = author
		row, err := models.PostNode(&obj)
		if (err != nil) || row <= 0 {
			return self.JSON(err, makross.StatusServiceUnavailable)
		}
		return self.JSON(obj)
	}
	herr.Message = "新增节点失败"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func PutNodeHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
		}
	}
	var author string
	if jwtUsername, okay := claims["Username"].(string); okay {
		author = jwtUsername
	}

	var obj models.Node
	self.Bind(&obj)

	if obj.Id > 0 {
		obj.Uid = uid
		obj.Author = author
		row, err := models.PutNode(obj.Id, &obj)
		if (err != nil) || (row <= 0) {
			return self.JSON(err, makross.StatusServiceUnavailable)
		}
		return self.JSON(obj)
	}
	herr.Message = "修改节点失败"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func DelNodeHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
		}
	}
	var isRoot bool
	if jwtIsRoot, okay := claims["IsRoot"].(bool); okay {
		isRoot = jwtIsRoot
	} else {
		herr.Message = "尚无权限"
	}

	id := self.Args("id").MustInt64()

	if isRoot && (id > 0) {
		err := models.DelNode(id, uid, -1000)
		if err != nil {
			herr.Message = err.Error()
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
		herr.Status = makross.StatusOK
		herr.Message = "删除节点成功"
		return self.JSON(herr)
	}
	herr.Message = "删除节点失败"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}
