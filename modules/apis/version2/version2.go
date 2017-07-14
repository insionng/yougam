package version2

import (
	"fmt"
	"strings"
	"time"

	"github.com/insionng/makross"
	"github.com/insionng/makross/jwt"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/models"
)

var herr = new(makross.HTTPError)

//Handler命名规范 请求方法+请求对象

// GetVersion 获取版本
func GetVersionHandler(self *makross.Context) error {
	var m = map[string]interface{}{}
	m["version"] = "2.0.0" //当服务端版本迭代产生不兼容时修改此版本号
	return self.JSON(m)
}

// GetPongHandler 乒乓 心跳Handler
func GetPongHandler(self *makross.Context) error {
	var m = map[string]interface{}{}
	if tokenString, okay := self.Get("TokenString").(string); okay && (len(tokenString) > 0) {
		m["Authorization"] = fmt.Sprintf("%v %v", jwt.Bearer, tokenString)
	}
	return self.JSON(m)
}

// PostSignupHandler 注册用户
func PostSignupHandler(self *makross.Context) error {

	username := self.Args("username").String()
	nickname := self.Args("nickname").String()
	password := self.Args("password").String()
	mobile := self.Args("mobile").String()
	gender := self.Args("gender").MustInt64()
	email := self.Args("email").String()
	content := self.Args("content").String() //个人简介 个人签名 个性说明之类
	group := self.Args("group").String()
	role := self.Args("role").MustInt64()

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	if len(password) > 0 {
		if helper.CheckPassword(password) == false {
			herr.Message = "密码含有非法字符或密码过短(至少4~30位密码)!"
			return self.JSON(herr, makross.StatusServiceUnavailable)

		}
	} else {
		herr.Message = "密码为空!"
		return self.JSON(herr, makross.StatusServiceUnavailable)

	}

	if len(username) == 0 {
		herr.Message = "用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空!"
		return self.JSON(herr, makross.StatusServiceUnavailable)

	}

	if len(email) > 0 {
		if helper.CheckEmail(email) == false {
			herr.Message = "Email格式错误!"
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
	} else {
		herr.Message = "Email地址为空!"
		return self.JSON(herr, makross.StatusServiceUnavailable)
	}

	if len(email) > 0 {
		if usrinfo, err := models.GetUserByEmail(email); usrinfo != nil {

			if usrinfo, err := models.GetUserByUsername(username); usrinfo != nil {
				herr.Message = "此用户名不能使用!"
				return self.JSON(herr, makross.StatusServiceUnavailable)

			} else if err != nil {

				herr.Message = "检索用户名账号期间出错!"
				return self.JSON(herr, makross.StatusServiceUnavailable)

			}

			herr.Message = "此Email不能使用!"
			return self.JSON(herr, makross.StatusServiceUnavailable)

		} else if err != nil {

			herr.Message = "检索EMAIL账号期间出错!"
			return self.JSON(herr, makross.StatusServiceUnavailable)

		}
	} else {
		if usrinfo, err := models.GetUserByUsername(username); usrinfo != nil {

			herr.Message = "此用户名已经被注册,请重新命名!"
			return self.JSON(herr, makross.StatusServiceUnavailable)

		} else if err != nil {

			herr.Message = "检索账号数据期间出错!"
			return self.JSON(herr, makross.StatusServiceUnavailable)

		}
	}

	if role == 0 {
		role = 1
	}

	if usrid, err := models.AddUser(email, username, nickname, "", helper.EncryptHash(password, nil), group, content, mobile, gender, role); err != nil && usrid <= 0 {

		herr.Message = "用户注册信息写入数据库时发生错误!"
		return self.JSON(herr, makross.StatusServiceUnavailable)

	} else {

		if usrinfo, err := models.GetUser(usrid); err == nil && usrinfo != nil {
			///注册成功
			models.PutSignin2User(usrinfo.Id, time.Now().Unix(), usrinfo.SigninCount+1, self.RealIP())

			//返回数据
			return self.JSON(usrinfo)

		} else {

			herr.Message = "获取用户数据出错!"
			return self.JSON(herr, makross.StatusServiceUnavailable)

		}

	}
}

// PostSignin 用户登录
func PostSigninHandler(self *makross.Context) error {
	herr.Message = "ErrUnauthorized"
	herr.Status = makross.StatusUnauthorized

	password := self.Args("password").String()
	if len(password) == 0 {
		herr.Message = "密码为空~"
		return self.JSON(herr, makross.StatusUnauthorized)
	}

	if helper.CheckPassword(password) == false {
		herr.Message = "密码含有非法字符或密码过短(至少4~30位密码)!"
		return self.JSON(herr, makross.StatusUnauthorized)
	}

	var err error
	var usr = new(models.User)
	var email, username string
	uoe := self.Args("username").String()
	mobile := self.Args("mobile").String()

	if (len(uoe) == 0) && (len(mobile) == 0) {
		herr.Message = "用户名不能少于4个字或多于30个字,登录账号至少有email或手机以及用户名之一进行登录,不能都为空!"
		return self.JSON(herr, makross.StatusUnauthorized)
	}

	switch {
	//mobile账号校验分支
	case len(mobile) > 0:
		{
			if helper.CheckUsername(mobile) == false {
				herr.Message = "手机号码不能包含非法字符,不能少于4个字或多于30个字!"
				return self.JSON(herr, makross.StatusUnauthorized)
			}

			if usr, err = models.GetUserByMobile(mobile); usr != nil && err == nil {
				if !helper.ValidateHash(usr.Password, password) {
					herr.Message = "密码无法通过校验!"
					return self.JSON(herr, makross.StatusUnauthorized)

				}
			} else {
				herr.Message = "该手机号码不存在!"
				return self.JSON(herr, makross.StatusUnauthorized)

			}
		}

	//默认账号校验分支
	default:
		if isEmail := strings.Contains(uoe, "@"); isEmail {
			email = uoe
			if len(email) == 0 {
				herr.Message = "EMAIL为空~"
				return self.JSON(herr, makross.StatusUnauthorized)
			}

			if helper.CheckEmail(email) == false {
				herr.Message = "Email格式不合符规格~"
				return self.JSON(herr, makross.StatusUnauthorized)
			}

			usr, err = models.GetUserByEmail(email)
		} else {
			username = uoe
			if len(username) == 0 {
				herr.Message = "用户名称为空~"
				return self.JSON(herr, makross.StatusUnauthorized)
			}

			if helper.CheckUsername(username) == false {
				herr.Message = "用户名称格式不合符规格~"
				return self.JSON(herr, makross.StatusUnauthorized)
			}

			usr, err = models.GetUserByUsername(username)
		}

	}

	if (usr != nil) && (err == nil) {
		if helper.ValidateHash(usr.Password, password) {
			models.PutSignin2User(usr.Id, time.Now().Unix(), usr.SigninCount+1, self.RealIP())
			claims := jwt.NewMapClaims()
			claims["IsRoot"] = (usr.Role == -1000)
			claims["UserId"] = usr.Id
			claims["Username"] = usr.Username
			claims["exp"] = time.Now().Add(jwt.DefaultJWTConfig.Expires).Unix()
			var data = map[string]interface{}{}
			var secret string
			if signingKey, okay := jwt.DefaultJWTConfig.SigningKey.(string); okay {
				secret = signingKey
			}
			data["token"], _ = jwt.NewTokenString(secret, "HS256", claims)
			data["user"] = usr
			return self.JSON(data)
		} else {
			herr.Message = "密码无法通过校验~"
			return self.JSON(herr, makross.StatusUnauthorized)
		}
	} else {
		herr.Message = "该账号不存在~"
		return self.JSON(herr, makross.StatusUnauthorized)
	}
}

// GetSignout 客户端执行清除 cookie 或 local storage时触发GetSignout进行记录动作
func GetSignoutHandler(self *makross.Context) error {
	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			return self.JSON(nil)
		}
	}
	_, e := models.PutSignout2User(uid, time.Now().Unix(), self.RealIP())
	return self.JSON(e)
}

// PostComment 发布评论
func PostCommentHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	claims := jwt.GetMapClaims(self)
	var uid int64
	if jwtUserId, okay := claims["UserId"].(float64); okay {
		uid = int64(jwtUserId)
		if uid <= 0 {
			herr.Message = "尚未登录"
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}
	}

	var author string
	if jwtUsername, okay := claims["Username"].(string); okay {
		author = jwtUsername
	}

	rid := self.Param("id").MustInt64() //reply id
	if rid <= 0 {
		rid = self.Args("id").MustInt64()
	}

	var rpy models.Reply
	self.Bind(&rpy)

	if usrinfo, err := models.GetUser(uid); (err == nil) && (usrinfo != nil) {

		rpy.Uid = uid
		rpy.Author = author

		if rid <= 0 {
			//全新发布
			if rid, err := models.PostReply(rpy.Tid, &rpy); err != nil || rid <= 0 {
				herr.Message = "回复内容写入数据库时发生错误"
				return self.JSON(herr, makross.StatusServiceUnavailable)

			} else {

				if rp, err := models.GetReply(rid); err == nil {
					return self.JSON(rp)

				} else {
					herr.Message = "获取回复内容数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)

				}

			}
		} else {
			//对指定的回复内容进行更新
			if row, err := models.PutReply(rid, &rpy); err != nil || row <= 0 {
				herr.Message = "更新回复写入数据库时发生错误"
				return self.JSON(herr, makross.StatusServiceUnavailable)

			} else {

				if rp, err := models.GetReply(rid); err == nil {
					return self.JSON(rp)

				} else {
					herr.Message = "获取回复内容数据出错"
					return self.JSON(herr, makross.StatusServiceUnavailable)

				}

			}
		}

	} else {
		herr.Message = "获取用户数据出错"
		return self.JSON(herr, makross.StatusServiceUnavailable)

	}
}

// GetComment 获取评论
func GetCommentHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	tid := self.Args("tid").MustInt64()

	if tid > 0 {
		if rps := models.GetReplysByTid(tid, 0, 0, 0, "id"); rps != nil {
			return self.JSON(rps)
		}

	}
	return self.JSON(herr, makross.StatusServiceUnavailable)

}

// PostReport 举报或反馈
func PostReportHandler(self *makross.Context) error {

	herr.Message = "ErrServiceUnavailable"
	herr.Status = makross.StatusServiceUnavailable

	id := self.Args("contentid").MustInt64()
	rid := self.Args("commentid").MustInt64()
	tid := self.Args("topicid").MustInt64()
	userid := self.Args("userid").MustInt64()
	content := self.Args("content").String()
	ctype := self.Args("ctype").MustInt64()

	if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {

		claims := jwt.GetMapClaims(self)
		jwtUserId := claims["UserId"].(float64)
		if suid := int64(jwtUserId); (suid > 0) && (usrinfo.Id == suid) {
			if id <= 0 {

				if rid <= 0 && tid > 0 {
					id = tid
					ctype = 1
				} else if rid > 0 && tid <= 0 {
					id = rid
					ctype = -1
				} else {
					return self.JSON(herr, makross.StatusUnauthorized)
				}
			}

			//如果已经举报过..
			d := map[string]int64{}
			if models.IsReportMark(userid, id, ctype) {

				d["id"] = id
				return self.JSON(d)

			} else {
				//保存举报内容
				if row, err := models.SetReportMark(userid, id, ctype, content); err != nil || row <= 0 {
					herr.Message = err.Error()
					return self.JSON(herr, makross.StatusServiceUnavailable)
				} else {
					d["id"] = id
					return self.JSON(d)
				}

			}

		} else {
			herr.Message = "不是当前用户无权操作!"
			return self.JSON(herr, makross.StatusServiceUnavailable)
		}

	} else {
		herr.Message = "获取用户数据出错!"
		return self.JSON(herr, makross.StatusServiceUnavailable)

	}
}
