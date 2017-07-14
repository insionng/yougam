package version1

import (
	"fmt"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"

	"time"

	"github.com/insionng/yougam/helper"
	sjson "github.com/insionng/yougam/libraries/bitly/go-simplejson"
	"github.com/insionng/yougam/models"
)

func GetApisHandler(self *makross.Context) error {
	m := map[string]string{}
	m["version"] = "1.0.0"
	crypted, err := helper.SetJsonCOMEncrypt(1, "", m)
	if err != nil {
		return self.Break(makross.StatusInternalServerError, err)
	}
	return self.String(crypted)
}

func PostApisHandler(self *makross.Context) error {

	_usr_, _ := self.Session.Get("SignedUser").(*models.User)
	var decrypts string
	if self.Get("decrypts") != nil {
		if d, okay := self.Get("decrypts").(string); okay {
			decrypts = d
		}
		//fmt.Println("decrypts:", decrypts)
	} else {

		crypted, _ := helper.SetJsonCOMEncrypt(0, "提交的数据不能为空!", nil)
		return self.String(crypted)

	}

	//fmt.Println("ApisPostHandler,36:", decrypts)
	if j, err := sjson.NewJson([]byte(decrypts)); err != nil {

		crypted, _ := helper.SetJsonCOMEncrypt(0, fmt.Sprint("解析json出错:", err.Error()), nil)
		return self.String(crypted)

	} else {

		if action, err := j.Get("action").String(); err == nil {

			switch {
			//注册用户
			case action == "userSignup":
				{
					username, _ := j.Get("username").String()
					nickname, _ := j.Get("nickname").String()
					password, _ := j.Get("password").String()
					mobile, _ := j.Get("mobile").String()
					gender, _ := j.Get("gender").Int64()
					email, _ := j.Get("email").String()
					content, _ := j.Get("content").String() //个人简介 个人签名 个性说明之类
					group, _ := j.Get("group").String()

					if len(password) > 0 {
						if helper.CheckPassword(password) == false {
							crypted, _ := helper.SetJsonCOMEncrypt(0, "密码含有非法字符或密码过短(至少4~30位密码)!", nil)
							return self.String(crypted)

						}
					} else {
						crypted, _ := helper.SetJsonCOMEncrypt(0, "密码为空!", nil)
						return self.String(crypted)

					}

					if len(username) == 0 {
						crypted, _ := helper.SetJsonCOMEncrypt(0, "用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空!", nil)
						return self.String(crypted)

					}

					if len(email) > 0 {
						if helper.CheckEmail(email) == false {
							crypted, _ := helper.SetJsonCOMEncrypt(0, "Email格式错误!", nil)
							return self.String(crypted)

						}
					} else {
						crypted, _ := helper.SetJsonCOMEncrypt(0, "Email地址为空!", nil)
						return self.String(crypted)

					}

					if len(email) > 0 {
						if usrinfo, err := models.GetUserByEmail(email); usrinfo != nil {

							if usrinfo, err := models.GetUserByUsername(username); usrinfo != nil {
								crypted, _ := helper.SetJsonCOMEncrypt(0, "此用户名不能使用!", nil)
								return self.String(crypted)

							} else if err != nil {

								crypted, _ := helper.SetJsonCOMEncrypt(0, "检索用户名账号期间出错!", nil)
								return self.String(crypted)

							}

							crypted, _ := helper.SetJsonCOMEncrypt(0, "此Email不能使用!", nil)
							return self.String(crypted)

						} else if err != nil {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "检索EMAIL账号期间出错!", nil)
							return self.String(crypted)

						}
					} else {
						if usrinfo, err := models.GetUserByUsername(username); usrinfo != nil {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "此用户名已经被注册,请重新命名!", nil)
							return self.String(crypted)

						} else if err != nil {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "检索账号数据期间出错!", nil)
							return self.String(crypted)

						}
					}

					if usrid, err := models.AddUser(email, username, nickname, "", helper.EncryptHash(password, nil), group, content, mobile, gender, 1); err != nil && usrid <= 0 {

						crypted, _ := helper.SetJsonCOMEncrypt(0, "用户注册信息写入数据库时发生错误!", nil)
						return self.String(crypted)

					} else {

						if usrinfo, err := models.GetUser(usrid); err == nil && usrinfo != nil {

							///注册成功 设置self.Session.on
							self.Session.Set("SignedUser", usrinfo)

							IsSigned := true
							IsRoot := (usrinfo.Role == -1000)
							self.Set("IsSigned", IsSigned)
							self.Set("IsRoot", IsRoot)
							self.Set("User", usrinfo)
							self.Set("UserID", usrinfo.Id)
							self.Set("UserName", usrinfo.Username)
							cc := cache.Store(self)
							cc.Set(fmt.Sprintf("SignedUser:%v", usrinfo.Id), usrinfo, 60*60*60)
							models.PutSignin2User(usrinfo.Id, time.Now().Unix(), usrinfo.SigninCount+1, self.RealIP())

							//返回数据

							crypted, _ := helper.SetJsonCOMEncrypt(1, "", usrinfo)
							return self.String(crypted)

						} else {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "获取用户数据出错!", nil)
							return self.String(crypted)

						}

					}
				}

			//用户登录
			case action == "userSignin":
				{
					username, _ := j.Get("username").String()
					password, _ := j.Get("password").String()
					mobile, _ := j.Get("mobile").String()
					email, _ := j.Get("email").String()

					if password != "" {
						if helper.CheckPassword(password) == false {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "密码含有非法字符或密码过短(至少4~30位密码)!", nil)
							return self.String(crypted)

						}
					} else {

						crypted, _ := helper.SetJsonCOMEncrypt(0, "密码为空!", nil)
						return self.String(crypted)

					}

					if len(email) > 0 {
						if helper.CheckEmail(email) == false {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "邮箱格式错误!", nil)
							return self.String(crypted)

						}

					}

					if len(mobile) > 0 {
						if helper.CheckUsername(mobile) == false {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "手机号码不能包含非法字符,不能少于4个字或多于30个字!", nil)
							return self.String(crypted)

						}

					}

					if len(username) > 0 {

						//TODO  增加手机与用户名交换  以及 校验
						//若果 username实际是email,则交换到email 并清空username
						if helper.CheckEmail(username) == true {
							email = username
							username = ""
						} else if helper.CheckUsername(username) == false {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "用户名不能含有非法字符,不能少于4个字或多于30个字,不能为空!", nil)
							return self.String(crypted)

						}

					}

					if (len(username) == 0) && (len(email) == 0) && (len(mobile) == 0) {
						crypted, _ := helper.SetJsonCOMEncrypt(0, "用户名不能少于4个字或多于30个字,登录账号至少有email或手机以及用户名之一进行登录,不能都为空!", nil)
						return self.String(crypted)

					}

					switch {
					//email账号校验分支
					case len(email) > 0:
						{
							if usrinfo, err := models.GetUserByEmail(email); usrinfo != nil && err == nil {

								if helper.ValidateHash(usrinfo.Password, password) {
									///登录成功 设置self.Session.on
									self.Session.Set("SignedUser", usrinfo)

									curUser := usrinfo
									curIsSigned := true
									curIsRoot := (curUser.Role == -1000)
									self.Set("IsSigned", curIsSigned)
									self.Set("IsRoot", curIsRoot)
									self.Set("SignedUser", curUser)

									cc := cache.Store(self)
									cc.Set(fmt.Sprintf("SignedUser:%v", curUser.Id), curUser, 60*60*60)
									models.PutSignin2User(curUser.Id, time.Now().Unix(), curUser.SigninCount+1, self.RealIP())

									//返回数据

									crypted, _ := helper.SetJsonCOMEncrypt(1, "", usrinfo)
									return self.String(crypted)

								} else {
									crypted, _ := helper.SetJsonCOMEncrypt(0, "密码无法通过校验!", nil)
									return self.String(crypted)

								}
							} else {

								crypted, _ := helper.SetJsonCOMEncrypt(0, "该账号不存在!", nil)
								return self.String(crypted)

							}
						}

					//mobile账号校验分支
					case len(mobile) > 0:
						{
							if usrinfo, err := models.GetUserByMobile(mobile); usrinfo != nil && err == nil {

								if helper.ValidateHash(usrinfo.Password, password) {
									///登录成功 设置self.Session.on
									self.Session.Set("SignedUser", usrinfo)

									curUser := usrinfo
									curIsSigned := true
									curIsRoot := (curUser.Role == -1000)
									self.Set("IsSigned", curIsSigned)
									self.Set("IsRoot", curIsRoot)
									self.Set("SignedUser", curUser)

									cc := cache.Store(self)
									cc.Set(fmt.Sprintf("SignedUser:%v", curUser.Id), curUser, 60*60*60)
									models.PutSignin2User(curUser.Id, time.Now().Unix(), curUser.SigninCount+1, self.RealIP())

									//返回数据

									crypted, _ := helper.SetJsonCOMEncrypt(1, "", usrinfo)
									return self.String(crypted)

								} else {

									crypted, _ := helper.SetJsonCOMEncrypt(0, "密码无法通过校验!", nil)
									return self.String(crypted)

								}
							} else {

								crypted, _ := helper.SetJsonCOMEncrypt(0, "该手机号码不存在!", nil)
								return self.String(crypted)

							}
						}

					//默认为Username账号校验分支
					default:
						{
							if usrinfo, err := models.GetUserByUsername(username); usrinfo != nil && err == nil {

								if helper.ValidateHash(usrinfo.Password, password) {
									///登录成功 设置self.Session.on
									self.Session.Set("SignedUser", usrinfo)

									curUser := usrinfo
									curIsSigned := true
									curIsRoot := (curUser.Role == -1000)
									self.Set("IsSigned", curIsSigned)
									self.Set("IsRoot", curIsRoot)
									self.Set("SignedUser", curUser)

									cc := cache.Store(self)
									cc.Set(fmt.Sprintf("SignedUser:%v", curUser.Id), curUser, 60*60*60)
									models.PutSignin2User(curUser.Id, time.Now().Unix(), curUser.SigninCount+1, self.RealIP())

									//返回数据

									crypted, _ := helper.SetJsonCOMEncrypt(1, "", usrinfo)
									return self.String(crypted)

								} else {

									crypted, _ := helper.SetJsonCOMEncrypt(0, "密码无法通过校验!", nil)
									return self.String(crypted)

								}
							} else {

								crypted, _ := helper.SetJsonCOMEncrypt(0, "该账号不存在!", nil)
								return self.String(crypted)

							}
						}

					}
				}

			//退出登录
			case action == "userSignout":
				{
					if _usr_ != nil {
						cc := cache.Store(self)
						cc.Delete(fmt.Sprintf("User:%v", _usr_.Id))
					}
					//退出，销毁self.Session.on
					self.Session.Delete("User")
					crypted, _ := helper.SetJsonCOMEncrypt(1, "", nil)
					return self.String(crypted)

				}
			//修改用户资料
			case action == "saveUser":
				{
					userid, _ := j.Get("userid").Int64()
					nickname, _ := j.Get("nickname").String()
					gender, _ := j.Get("gender").Int64()
					avatar, _ := j.Get("avatar").String()

					if usrinfo, err := models.GetUser(userid); err == nil {

						if (_usr_ != nil) && (usrinfo.Id == _usr_.Id) {

							usrinfo.Gender = gender
							usrinfo.Nickname = nickname
							usrinfo.Avatar = avatar
							if row, err := models.PutUser(userid, usrinfo); err != nil || row == 0 {

								crypted, _ := helper.SetJsonCOMEncrypt(0, "用户资料写入数据库时发生错误!", nil)
								return self.String(crypted)

							} else {

								if usrinfo, err := models.GetUser(userid); err == nil {

									crypted, _ := helper.SetJsonCOMEncrypt(1, "", usrinfo)
									return self.String(crypted)

								} else {

									crypted, _ := helper.SetJsonCOMEncrypt(0, "获取用户数据出错!", nil)
									return self.String(crypted)

								}

							}
						} else {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "不是当前用户无权修改数据!", nil)
							return self.String(crypted)

						}

					} else {

						crypted, _ := helper.SetJsonCOMEncrypt(0, "获取用户数据出错!", nil)
						return self.String(crypted)

					}
				}

			//发布内容 or 更新内容
			case action == "postContent":
				{
					id, _ := j.Get("id").Int64()
					userid, _ := j.Get("userid").Int64()
					nickname, _ := j.Get("author").String()
					title, _ := j.Get("title").String()
					excerpt, _ := j.Get("excerpt").String()
					content, _ := j.Get("content").String()
					attachment, _ := j.Get("attachment").String() // '图片使用,分割'
					tags, _ := j.Get("tags").String()
					pid, _ := j.Get("pid").Int64()
					cid, _ := j.Get("cid").Int64()
					tailinfo, _ := j.Get("tailinfo").String()
					nodeid, _ := j.Get("nodeid").Int64()
					node, _ := j.Get("node").String()
					category, _ := j.Get("category").String()
					posttime, _ := j.Get("posttime").String()
					latitude, _ := j.Get("latitude").Float64()   //纬度
					longitude, _ := j.Get("longitude").Float64() //经度

					if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {

						if (_usr_ != nil) && (usrinfo.Id == _usr_.Id) {

							tp := &models.Topic{}
							tp.Uid = userid
							tp.Nid = nodeid
							tp.Cid = cid
							tp.Category = category
							tp.Node = node
							tp.Author = nickname
							tp.Title = title
							tp.Excerpt = excerpt
							tp.Content = content
							tp.Tailinfo = tailinfo
							tp.Attachment = attachment // '图片使用,分割'
							tp.Tags = tags
							tp.Pid = pid
							tp.Created, _ = helper.String2UnixNano(posttime)
							tp.Latitude = latitude
							tp.Longitude = longitude

							if id <= 0 {
								//全新发布
								if tid, err := models.PostTopic(tp); err != nil || tid <= 0 {

									crypted, _ := helper.SetJsonCOMEncrypt(0, fmt.Sprint("发布内容写入数据库时发生错误!", err), nil)
									return self.String(crypted)

								} else {

									if tp, err := models.GetTopic(tid); err == nil {

										crypted, _ := helper.SetJsonCOMEncrypt(1, "", tp)
										return self.String(crypted)

									} else {

										crypted, _ := helper.SetJsonCOMEncrypt(0, "获取话题内容数据出错!", nil)
										return self.String(crypted)

									}

								}
							} else {
								//对指定topicid的话题进行更新
								if row, err := models.PutTopic(id, tp); err != nil || row <= 0 {

									crypted, _ := helper.SetJsonCOMEncrypt(0, "更新内容写入数据库时发生错误!", nil)
									return self.String(crypted)

								} else {

									if tp, err := models.GetTopic(id); err == nil {

										crypted, _ := helper.SetJsonCOMEncrypt(1, "", tp)
										return self.String(crypted)

									} else {

										crypted, _ := helper.SetJsonCOMEncrypt(0, "获取话题内容数据出错!", nil)
										return self.String(crypted)

									}

								}
							}
						} else {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "不是当前用户无权修改数据!", nil)
							return self.String(crypted)

						}

					} else {
						crypted, _ := helper.SetJsonCOMEncrypt(0, "获取用户数据出错!", nil)
						return self.String(crypted)

					}
				}

			//获取内容
			case action == "getContent":
				{
					if tid, err := j.Get("contentid").Int64(); err != nil {
						crypted, _ := helper.SetJsonCOMEncrypt(0, err.Error(), nil)
						return self.String(crypted)

					} else {

						if tid > 0 {
							if tp, err := models.GetTopic(tid); tp != nil && err == nil {
								tp.Views = tp.Views + 1
								//tp.Hotup = tp.Hotup + 1
								if row, e := models.PutTopic(tid, tp); e != nil {
									fmt.Println("getcontent更新话题ID", tid, "访问次数数据错误,row:", row, e)
									//self.Ctx.Output.SetStatus(500)
								}

								crypted := ""
								if tps := models.GetTopicsByPid(tid, 0, 0, 0, "id"); tps != nil && (len(*tps) > 0) {
									crypted, _ = helper.SetJsonCOMEncrypt(1, "", tps)

								} else {
									crypted, _ = helper.SetJsonCOMEncrypt(0, fmt.Sprintf("读取主题ID为%v的数据发生错误!", tid), nil)

								}
								//j["data"] = tp
								//j["subcontent"], _ = models.GetSubTopics(tid, 0, 100, "id")
								//crypted, _ := helper.SetJsonCOMEncrypt(1, "", j)
								return self.String(crypted)

							}

						} else {
							return self.NoContent(501)
						}

					}
				}

			//发布评论
			case action == "postComment":
				{
					rid, _ := j.Get("id").Int64() //reply id
					userid, _ := j.Get("userid").Int64()
					nickname, _ := j.Get("author").String()
					content, _ := j.Get("content").String()
					attachment, _ := j.Get("attachment").String() // '图片使用,分割'
					pid, _ := j.Get("pid").Int64()                //another reply id
					tid, _ := j.Get("topicid").Int64()
					tailinfo, _ := j.Get("tailinfo").Int64()
					latitude, _ := j.Get("latitude").Float64()   //纬度
					longitude, _ := j.Get("longitude").Float64() //经度

					if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {

						if (_usr_ != nil) && (usrinfo.Id == _usr_.Id) {

							rpy := &models.Reply{}
							rpy.Tid = tid
							rpy.Uid = userid
							rpy.Pid = pid
							rpy.Author = nickname
							rpy.Content = content
							rpy.Tailinfo = tailinfo
							rpy.Attachment = attachment // '图片使用,分割'
							rpy.Latitude = latitude
							rpy.Longitude = longitude

							if rid <= 0 {
								//全新发布
								if rid, err := models.PostReply(tid, rpy); err != nil || rid <= 0 {
									crypted, _ := helper.SetJsonCOMEncrypt(0, "回复内容写入数据库时发生错误", nil)
									return self.String(crypted)

								} else {

									if rp, err := models.GetReply(rid); err == nil {

										crypted, _ := helper.SetJsonCOMEncrypt(1, "", rp)
										return self.String(crypted)

									} else {
										crypted, _ := helper.SetJsonCOMEncrypt(0, "获取回复内容数据出错", nil)
										return self.String(crypted)

									}

								}
							} else {
								//对指定的回复内容进行更新
								if row, err := models.PutReply(rid, rpy); err != nil || row <= 0 {

									crypted, _ := helper.SetJsonCOMEncrypt(0, "更新回复写入数据库时发生错误", nil)
									return self.String(crypted)

								} else {

									if rp, err := models.GetReply(rid); err == nil {

										crypted, _ := helper.SetJsonCOMEncrypt(1, "", rp)
										return self.String(crypted)

									} else {
										crypted, _ := helper.SetJsonCOMEncrypt(0, "获取话题内容数据出错", nil)
										return self.String(crypted)
									}

								}
							}
						} else {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "不是当前用户无权修改数据", nil)
							return self.String(crypted)

						}

					} else {

						crypted, _ := helper.SetJsonCOMEncrypt(0, "获取用户数据出错", nil)
						return self.String(crypted)

					}
				}

			//获取评论
			case action == "getComment":
				{
					if tid, err := j.Get("contentid").Int64(); err != nil {

						crypted, _ := helper.SetJsonCOMEncrypt(0, err.Error(), nil)
						return self.String(crypted)

					} else {

						if tid > 0 {
							if rps := models.GetReplysByTid(tid, 0, 0, 0, "id"); rps != nil {

								crypted, _ := helper.SetJsonCOMEncrypt(1, "", rps)
								return self.String(crypted)

							}

						} else {
							return self.NoContent(501)
						}

					}
				}

			//举报
			case action == "reportContent":
				{
					id, _ := j.Get("contentid").Int64()
					rid, _ := j.Get("commentid").Int64()
					tid, _ := j.Get("topicid").Int64()
					userid, _ := j.Get("userid").Int64()
					content, _ := j.Get("content").String()
					ctype, _ := j.Get("ctype").Int64()

					if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {

						if (_usr_ != nil) && (usrinfo.Id == _usr_.Id) {
							if id <= 0 {

								if rid <= 0 && tid > 0 {
									id = tid
									ctype = 1
								} else if rid > 0 && tid <= 0 {
									id = rid
									ctype = -1
								} else {

									return self.NoContent(401)
								}
							}

							//如果已经举报过..
							crypted := ""
							d := map[string]int64{}
							if models.IsReportMark(userid, id, ctype) {

								d["id"] = id
								crypted, _ = helper.SetJsonCOMEncrypt(1, "", d)
								return self.String(crypted)

							} else {
								//保存举报内容
								if row, err := models.SetReportMark(userid, id, ctype, content); err != nil || row <= 0 {

									crypted, _ = helper.SetJsonCOMEncrypt(0, fmt.Sprint(err), nil)
								} else {
									d["id"] = id
									crypted, _ = helper.SetJsonCOMEncrypt(1, "", d)

								}

								return self.String(crypted)

							}

						} else {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "不是当前用户无权操作!", nil)
							return self.String(crypted)

						}

					} else {

						crypted, _ := helper.SetJsonCOMEncrypt(0, "获取用户数据出错!", nil)
						return self.String(crypted)

					}
				}

			//获取用户话题列表
			case action == "getUserPostList":
				{
					userid, _ := j.Get("userid").Int64()
					offset, _ := j.Get("offset").Int64()
					page, _ := j.Get("page").Int64()
					limit, _ := j.Get("limit").Int64()

					ctype, _ := j.Get("ctype").Int64()
					field, _ := j.Get("field").String()

					if field == "" {
						field = "id"
					}

					if limit < 0 {
						limit = 0
					}

					if page <= 0 {
						page = 1
					}

					if userid > 0 {
						if usrinfo, err := models.GetUser(userid); err == nil && usrinfo != nil {
							if offset <= 0 {
								_, _, _, _, offset := helper.Pages(usrinfo.TopicCount, page, limit)
								if tps := models.GetTopicsByUid(userid, int(offset), int(limit), ctype, field); tps != nil {

									crypted, _ := helper.SetJsonCOMEncrypt(1, "", tps)
									return self.String(crypted)

								}
							} else {
								if tps := models.GetTopicsByUid(userid, int(offset), int(limit), ctype, field); tps != nil {

									crypted, _ := helper.SetJsonCOMEncrypt(1, "", tps)
									return self.String(crypted)

								}
							}

						} else {

							crypted, _ := helper.SetJsonCOMEncrypt(0, "获取用户数据出错!", nil)
							return self.String(crypted)

						}

					} else {
						return self.NoContent(501)
					}
				}

			//获取首页话题列表
			case action == "getHomePostList":
				{

					offset, _ := j.Get("offset").Int()
					page, _ := j.Get("page").Int64()
					limit, _ := j.Get("limit").Int64()
					field, _ := j.Get("field").String()

					if field == "lastest" {
						field = "id"
					} else if field == "hotness" {
						field = "hotness"
					} else {
						field = "hotness"
					}

					if offset <= 0 {
						if results_count, err := models.GetTopicsCount(offset, int(limit)); err != nil {

							crypted, _ := helper.SetJsonCOMEncrypt(0, fmt.Sprint(err), nil)
							return self.String(crypted)

						} else {

							_, _, _, _, offset := helper.Pages(results_count, page, limit)

							if tps, err := models.GetTopics(int(offset), int(limit), field); err == nil {

								crypted, _ := helper.SetJsonCOMEncrypt(1, "", tps)
								return self.String(crypted)

							} else {
								crypted, _ := helper.SetJsonCOMEncrypt(0, fmt.Sprint("首页 数据查询出错", err), nil)
								return self.String(crypted)

							}

						}
					} else {
						if tps, err := models.GetTopics(offset, int(limit), field); err == nil {

							crypted, _ := helper.SetJsonCOMEncrypt(1, "", tps)
							return self.String(crypted)

						} else {

							crypted, _ := helper.SetJsonCOMEncrypt(0, fmt.Sprint("首页 数据查询出错", err), nil)
							return self.String(crypted)

						}
					}
				}

			//非内置请求动作均视作非法请求
			default:
				return self.NoContent(401)

			}

		} else {
			//非法请求
			return self.NoContent(401)
		}

	}

	return self.NoContent(401)
}
