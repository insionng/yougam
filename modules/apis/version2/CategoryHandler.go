package version2

import (
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
	"fmt"
	"github.com/insionng/makross"
	"github.com/insionng/makross/jwt"
)

func GetCategoryHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable
	nid := self.Args("id").MustInt64()

	if nid > 0 {
		obj, err := models.GetCategory(nid)
		if err != nil {
			herr.Message = err.Error()
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
		return self.JSON(obj)
	}
	herr.Message = "没有获取到分类ID"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

// GetCategoriesHandler
func GetCategoriesHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	offset := self.Args("offset").MustInt()
	page := self.Args("page").MustInt64()
	limit := self.Args("limit").MustInt()
	field := self.Args("field").String()
	pid := self.Args("pid").MustInt64()
	nid := self.Args("nid").MustInt64()
	ctype := self.Args("ctype").MustInt64()

	if field == "lastest" {
		field = "id"
	} else if (field == "hotness") || (len(field) == 0) {
		field = "hotness"
	}

	switch {
	case (ctype != 0) && (pid != 0): // 获取特定ctype分类里特定pid之分类列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetCategoriesByCtypeWithPid(0, limit, ctype, pid, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetCategoriesByCtypeWithPid(int(offset_), limit, ctype, pid, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetCategoriesByCtypeWithPid(offset, limit, ctype, pid, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	case (ctype == 0) && (pid != 0): // 获取特定分类之下级分类列表
		if offset <= 0 {
			var resultsCount int64
			if objs := models.GetCategoriesViaPid(pid, 0, limit, 0, field); objs != nil {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs := models.GetCategoriesViaPid(pid, int(offset_), limit, 0, field); objs != nil {
					return self.JSON(objs)
				} else {
					herr.Message = "获取分类数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			} else {
				herr.Message = "没有获取到分类数据"
				herr.Status = makross.StatusOK
				return self.JSON(herr)
			}
		} else {
			if objs := models.GetCategoriesViaPid(pid, offset, limit, 0, field); objs != nil {
				return self.JSON(objs)
			} else {
				herr.Message = "没有获取到分类数据"
				herr.Status = makross.StatusOK
				return self.JSON(herr)
			}
		}
	case (ctype != 0) && (nid > 0): // 获取特定ctype分类里特定nid之分类列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetCategoriesByCtypeWithNid(0, limit, ctype, nid, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetCategoriesByCtypeWithNid(int(offset_), limit, ctype, nid, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetCategoriesByCtypeWithNid(offset, limit, ctype, nid, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	case (ctype == 0) && (nid > 0): // 获取特定分类之下级分类列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetCategoriesByNid(nid, 0, limit, field); err != nil {
				herr.Message = fmt.Sprintf("获取分类数据发生错误:%v", err)
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetCategoriesByNid(nid, int(offset_), limit, field); err != nil {
					herr.Message = "获取分类数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)
				} else {
					return self.JSON(objs)
				}
			}

		} else {
			if objs, err := models.GetCategoriesByNid(nid, offset, limit, field); err != nil {
				herr.Message = "获取分类数据出错"
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				return self.JSON(objs)
			}
		}
	case (ctype != 0) && (pid == 0) && (nid <= 0): // 获取特定ctype分类列表
		if offset <= 0 {
			var resultsCount int64
			if objs, err := models.GetCategoriesByCtype(0, limit, ctype, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				resultsCount = int64(len(*objs))
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetCategoriesByCtype(int(offset_), limit, ctype, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetCategoriesByCtype(offset, limit, ctype, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	default: // 获取全部分类列表
		if offset <= 0 {
			if resultsCount, err := models.GetCategoriesCount(offset, limit); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				_, _, _, _, offset_ := helper.Pages(resultsCount, page, int64(limit))
				if objs, err := models.GetCategories(int(offset_), limit, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetCategories(offset, int(limit), field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	}
}

func PostCategoryHandler(self *makross.Context) error {
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

	var obj models.Category
	self.Bind(&obj)

	if (obj.Pid >= 0) && (len(obj.Title) > 0) && (uid > 0) {
		obj.Uid = uid
		obj.Author = author
		row, err := models.PostCategory(&obj)
		if (err != nil) || row <= 0 {
			return self.JSON(err, makross.StatusServiceUnavailable)
		}
		return self.JSON(obj)
	}
	herr.Message = "新增分类失败"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func PutCategoryHandler(self *makross.Context) error {
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

	var obj models.Category
	self.Bind(&obj)

	if obj.Id > 0 {
		obj.Uid = uid
		obj.Author = author
		row, err := models.PutCategory(obj.Id, &obj)
		if (err != nil) || (row <= 0) {
			return self.JSON(err, makross.StatusServiceUnavailable)
		}
		return self.JSON(obj)
	}
	herr.Message = "修改分类失败"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func DelCategoryHandler(self *makross.Context) error {
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
		err := models.DelCategory(id, uid, -1000)
		if err != nil {
			herr.Message = err.Error()
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
		herr.Status = makross.StatusOK
		herr.Message = "删除分类成功"
		return self.JSON(herr)
	}
	herr.Message = "删除分类失败"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}
