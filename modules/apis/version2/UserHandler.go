package version2

import (
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
	"github.com/insionng/makross"
	"github.com/insionng/makross/jwt"
)

// GetUsers
func GetUsersHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	offset := self.Args("offset").MustInt()
	page := self.Args("page").MustInt64()
	limit := self.Args("limit").MustInt()
	field := self.Args("field").String()
	role := self.Args("role").MustInt64()

	if field == "lastest" {
		field = "id"
	} else if (field == "hotness") || (len(field) == 0) {
		field = "hotness"
	}

	switch {
	case role != 0: // 获取特定角色用户列表
		if offset <= 0 {
			var results_count int64
			if qt, err := models.GetUsersByRole(role, 0, limit, field); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				results_count = int64(len(*qt))
				_, _, _, _, offset_ := helper.Pages(results_count, page, int64(limit))
				if objs, err := models.GetUsersByRole(role, int(offset_), limit, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetUsersByRole(role, offset, limit, field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	default: // 获取全部用户列表
		if offset <= 0 {
			if results_count, err := models.GetUsersCount(offset, limit); err != nil {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				_, _, _, _, offset_ := helper.Pages(results_count, page, int64(limit))
				if objs, err := models.GetUsers(int(offset_), limit, field); err == nil {
					return self.JSON(objs)
				} else {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				}

			}
		} else {
			if objs, err := models.GetUsers(offset, int(limit), field); err == nil {
				return self.JSON(objs)
			} else {
				herr.Message = err.Error()
				return self.JSON(herr, makross.StatusServiceUnavailable)
			}
		}
	}
}

func GetUserHandler(self *makross.Context) error {
	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable
	tid := self.Args("id").MustInt64()

	if tid != 0 {
		tp, err := models.GetUser(tid)
		if err != nil {
			herr.Message = err.Error()
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
		return self.JSON(tp)
	}
	herr.Message = "没有获取到用户ID"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func PostUserHandler(self *makross.Context) error {
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

	var user models.User
	self.Bind(&user)

	if isRoot && (len(user.Username) > 0) && (len(user.Password) > 0) {

		user.Password = helper.EncryptHash(user.Password, nil)
		if !helper.CheckEmail(user.Email) {
			user.Email = ""
		}

		tp, err := models.PostUser(&user)
		if err != nil {
			herr.Message = err.Error()
			return self.JSON(err, makross.StatusServiceUnavailable)
		}
		return self.JSON(tp)
	}
	herr.Message = "没有获取到用户数据"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

// PutUser 更新用户
func PutUserHandler(self *makross.Context) error {
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

	var user models.User
	self.Bind(&user)

	id := self.Args("id").MustInt64()
	if id <= 0 {
		id = user.Id
	}

	var allow bool
	usr, err := models.GetUser(id)
	if isRoot {
		allow = true
	} else {
		if err != nil {
			herr.Message = "获取用户数据出错!"
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
		if (uid > 0) && (usr.Id == uid) {
			allow = true
		} else {
			herr.Message = "不是当前用户无权修改数据!"
			return self.JSON(herr, makross.StatusUnauthorized)
		}
	}

	if allow && (id > 0) {

		user.Id = id
		user.Password = helper.EncryptHash(user.Password, nil)
		if !helper.CheckEmail(user.Email) {
			user.Email = ""
		}

		row, err := models.PutUser(id, &user)
		if (err != nil) || (row == 0) {
			herr.Message = err.Error()
			return self.JSON(herr, makross.StatusServiceUnavailable)
		} else {
			if u, e := models.GetUser(id); e != nil {
				herr.Message = "获取用户数据出错!"
				return self.JSON(herr, makross.StatusServiceUnavailable)
			} else {
				return self.JSON(u)
			}
		}

	}

	herr.Message = "没有获取到用户数据"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}

func DelUserHandler(self *makross.Context) error {
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
		err := models.DelUser(id, uid, -1000)
		if err != nil {
			herr.Message = err.Error()
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
		herr.Message = "删除用户成功"
		return self.JSON(herr)
	}
	herr.Message = "没有获取到用户数据"
	return self.JSON(herr, makross.StatusServiceUnavailable)
}
