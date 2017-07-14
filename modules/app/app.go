package app

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/libraries/cli"
	"github.com/insionng/yougam/models"
	"github.com/insionng/yougam/modules/auth"
	"github.com/insionng/yougam/modules/jsock"
	"github.com/insionng/yougam/modules/setting"
	"github.com/insionng/yougam/routers"
	"github.com/insionng/yougam/routers/apis"
	"github.com/insionng/yougam/routers/root"

	"github.com/insionng/makross"
	"github.com/insionng/makross/cache"
	"github.com/insionng/makross/captcha"
	"github.com/insionng/makross/cors"
	"github.com/insionng/makross/csrf"
	"github.com/insionng/makross/i18n"
	"github.com/insionng/makross/jwt"
	//"github.com/insionng/makross/macrof"
	//"github.com/insionng/makross/macrus"
	"github.com/insionng/makross/logger"
	"github.com/insionng/makross/pongor"
	"github.com/insionng/makross/recover"
	"github.com/insionng/makross/session"
	"github.com/insionng/makross/static"

	"github.com/insionng/yougam/modules/apis/version1"
	"github.com/insionng/yougam/modules/apis/version2"
)

func App(c *cli.Context) error {

	if err := models.NewEngine(); err != nil {
		return err
	} else {
		models.HasEngine = true
	}

	/*------------------------------------*/
	m := makross.New()
	/*------------------------------------*/
	m.Use(recover.Recover())
	//m.Use(macrus.New())
	m.Use(logger.Logger())
	m.Use(cors.CORS())
	/*------------------------------------*/

	var production bool
	if len(c.Args()) >= 3 {
		production, _ = strconv.ParseBool(c.Args()[2])
	}

	/*------------------------------------*/
	//中间件之间存在顺序及依赖关系,请勿随意移动位置.
	/*------------------------------------*/
	m.Static("/file/", "../file")
	/*------------------------------------*/
	staticPath := fmt.Sprintf("themes/%v/static", helper.Theme())
	/*------------------------------------*/
	if production {
		if e := os.RemoveAll("./public"); e != nil {
			log.Println("删除目录出错:", e)
		}
		//复制模板静态资源到public目录
		if e := helper.CopyDir(staticPath, "./public"); e != nil {
			log.Println("复制模板静态资源到public目录时出现错误:", e)
		}
		m.Use(static.Static("public"))
	} else {
		m.Use(static.Static(staticPath))
	}
	fmt.Println("Production:", production)
	/*------------------------------------*/
	templatesPath := fmt.Sprintf("themes/%v/templates", helper.Theme())
	m.SetRenderer(pongor.Renderor(pongor.Option{
		Directory: templatesPath,
		Reload:    !production,
	}))
	/*------------------------------------*/

	if helper.IsExist("./conf/locale/locale_zh-CN.ini") {
		m.Use(i18n.I18n(i18n.Options{
			Directory:   "conf/locale",
			DefaultLang: "zh-CN",
			Langs:       []string{"en-US", "zh-CN"},
			Names:       []string{"English", "简体中文"},
			Redirect:    true,
		}))
	}

	/*------------------------------------*/
	//如果在前端启用了Nginx自带的Gzip功能，那么此处可以注释掉
	//m.Use(compress.Gzip())
	/*------------------------------------*/
	m.Use(cache.Cacher())
	/*------------------------------------*/
	m.Use(captcha.Captchaer(captcha.Options{
		URLPrefix:        "/captcha/", // URL prefix of getting captcha pictures.
		FieldIDName:      "captchaid", // Hidden input element ID.
		FieldCaptchaName: "captcha",   // User input value element name in request form.
		ChallengeNums:    6,           // Challenge number.
		Width:            276,         // Captcha image width.
		Height:           80,          // Captcha image height.
		Expiration:       60,          // Captcha expiration time in seconds.
		CachePrefix:      "captcha_",  // Cache key prefix captcha characters.
	}))
	/*------------------------------------*/
	m.Use(session.Sessioner())
	/*------------------------------------*/
	//macrof.Wrapper(m)
	//-------------------------------------------------------------------------------------//

	m.File("/favicon.ico", fmt.Sprintf("themes/%v/static/img/favicon.png", helper.Theme()))

	g := m.Group("", setting.BaseMiddler())

	g.Get("/", routers.GetMainHandler)

	//首页 ?page
	g.Get("/page<page:\\d+>/", routers.GetMainHandler)
	//首页 hotness/confidence类
	g.Get("/topics/<tab:[A-Za-z]+>/", routers.GetMainHandler)
	//http://localhost/lastest/2/
	g.Get("/topics/<tab:[A-Za-z]+>/<page:(\\d+)>/", routers.GetMainHandler)
	//http://localhost/lastest/page2/
	g.Get("/topics/<tab:[A-Za-z]+>/page<page:(\\d+)>/", routers.GetMainHandler)

	//搜索话题
	g.Get("/search/", routers.GetMainHandler)
	//同时支持page和keyword参数
	g.Get("/search/<keyword:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/page<page:(\\d+)>/", routers.GetMainHandler)
	//支持keyword参数
	g.Get("/search/<keyword:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/", routers.GetMainHandler)

	//浏览分类 "/category/<tab:([A-Za-z]+)>/<category:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/"优先级必须高于"/category/<category:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/"
	g.Get("/category/<cid:\\d+>/", routers.GetMainHandler)
	g.Get("/category/<cid:\\d+>/page<page:(\\d+)>/", routers.GetMainHandler)
	g.Get("/category/<cid:\\d+>/<tab:([A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)

	g.Get("/category/<tab:([A-Za-z]+)>/<category:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/", routers.GetMainHandler)
	g.Get("/category/<tab:([A-Za-z]+)>/<category:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)

	g.Get("/category/<category:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/", routers.GetMainHandler)
	g.Get("/category/<category:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)
	g.Get("/category/<category:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/<tab:([A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)

	//浏览节点 "/node/<tab:([A-Za-z]+)>/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/"优先级必须高于"/node/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/"
	g.Get("/node/<nid:\\d+>/", routers.GetMainHandler)
	g.Get("/node/<nid:\\d+>/page<page:(\\d+)>/", routers.GetMainHandler)
	g.Get("/node/<nid:\\d+>/<tab:([A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)

	g.Get("/node/<tab:([A-Za-z]+)>/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/", routers.GetMainHandler)
	g.Get("/node/<tab:([A-Za-z]+)>/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)

	g.Get("/node/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/", routers.GetMainHandler)
	g.Get("/node/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)
	g.Get("/node/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/<tab:([A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)

	g.Get("/createdby/<username:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/", routers.GetMainHandler)
	g.Get("/createdby/<username:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)
	g.Get("/createdby/<username:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/<tab:([A-Za-z]+)>/page<page:(\\d+)>/", routers.GetMainHandler)

	//best
	g.Get("/best/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/", routers.GetBestHandler)
	g.Get("/best/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/", routers.GetBestHandler)
	g.Get("/best/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/page<page:(\\d+)>/", routers.GetBestHandler)
	g.Get("/best/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/<tab:([A-Za-z]+)>/page<page:(\\d+)>/", routers.GetBestHandler)

	//详情页面
	g.Get("/<tid:(\\d+)>/", routers.GetTopicHandler)
	g.Get("/topic/<tid:(\\d+)>/", routers.GetTopicHandler)

	g.Get("/identicon/<name>/<size>/default.png", routers.GetIdenticonHandler)

	g.Get("/page/<pageid:\\d+>/", routers.GetPageHandler)
	g.Get("/page/<page:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/", routers.GetPageHandler)

	g.Get("/user/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/", routers.GetUserHandler)
	g.Get("/user/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/<type>/", routers.GetUserHandler)
	g.Get("/user/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/<type>/page<page:(\\d+)>/", routers.GetUserHandler)
	g.Get("/user/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/<type>/<tab:([A-Za-z]+)>/page<page:(\\d+)>/", routers.GetUserHandler)
	g.Get("/users/", routers.GetUsersHandler)

	//touch view
	g.Get("/touch/view/<name:([A-Za-z]+)>/<id:\\d+>/", routers.GetTouchViewHandler)

	//登录
	g.Get("/signin/", csrf.CSRFWithConfig(csrf.CSRFConfig{
		TokenLookup: "form:csrf",
	}), routers.GetSigninHandler).Post(routers.PostSigninHandler)

	//退出
	g.Get("/signout/", routers.GetSignoutHandler)

	//注册
	g.Get("/signup/", routers.GetSignupHandler).Post(routers.PostSignupHandler)
	g.Get("/signup/<key>/", routers.GetSignupHandler)

	//忘记密码
	forgot := g.Group("/forgot")
	forgot.Get("/", routers.GetForgotHandler).Post(routers.PostForgotHandler)
	forgot.Get("/<key>/", routers.GetForgotHandler).Post(routers.PostForgotHandler)

	//root signin
	g.Get("/root/signin/", csrf.CSRFWithConfig(csrf.CSRFConfig{
		TokenLookup: "form:csrf",
	}), root.GetRSigninHandler).Post(root.PostRSigninHandler)
	//无需权限

	//-------------------------------------------------------------------------------------//

	u := g.Group("", setting.AuthWebMiddler())
	//创建节点
	u.Get("/new/node/", routers.GetNewNodeHandler).Post(routers.PostNewNodeHandler)
	u.Get("/new/node/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/", routers.GetNewNodeHandler).Post(routers.PostNewNodeHandler)              //创建下级节点 :node是上级的title
	u.Get("/new/node/<nid:\\d+>/", routers.GetNewNodeHandler).Post(routers.PostNewNodeHandler)                                         //创建下级节点 :nid是上级的id
	u.Get("/new/node/<category:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/category/", routers.GetNewNodeHandler).Post(routers.PostNewNodeHandler) //创建指定:category的下级节点 :nid是上级的id
	u.Get("/new/node/<cid:\\d+>/category/", routers.GetNewNodeHandler).Post(routers.PostNewNodeHandler)                                //创建指定:cid的下级节点 :nid是上级的id

	//创建话题
	u.Get("/subject/<tid>/topic/", routers.GetNewTopicHandler).Post(routers.PostNewTopicHandler)
	u.Get("/new/topic/", routers.GetNewTopicHandler).Post(routers.PostNewTopicHandler)
	u.Post("/new/topic/<tid>/", routers.PostNewTopicHandler)

	u.Get("/new/node/<nid:(\\d+)>/topic/", routers.GetNewTopicHandler).Post(routers.PostNewTopicHandler)
	u.Post("/new/node/<nid:(\\d+)>/topic/<tid>/", routers.PostNewTopicHandler)
	u.Get("/new/node/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/topic/", routers.GetNewTopicHandler).Post(routers.PostNewTopicHandler)
	u.Post("/new/node/<node:([\\x{4e00}-\\x{9fa5}A-Za-z]+)>/topic/<tid>/", routers.PostNewTopicHandler)

	//修改话题
	u.Get("/modify/<tid>/", routers.GetModifyTopicHandler).Post(routers.PostModifyTopicHandler)
	//支付话题
	u.Post("/pay/<name:([A-Za-z]+)>/<amount:\\d+>/<id:\\d+>/", routers.PostPaymentHandler)
	//收藏话题
	u.Get("/favorite/<name:([A-Za-z]+)>/<id:\\d+>/", routers.GetFavoriteHandler)

	//创建回复
	u.Post("/subject/<tid>/comment/", routers.PostNewReplyHandler)

	//通知提醒
	u.Get("/notifications/", routers.NotificGetHandler)
	u.Get("/notifications/<page:(\\d+)>/", routers.NotificGetHandler)
	u.Get("/notifications/page<page:(\\d+)>/", routers.NotificGetHandler)

	//删除通知
	u.Get("/delete/notification/<notificid>/", routers.DeleteNotificGetHandler)

	//个人设置
	u.Get("/settings/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/", routers.SettingsGetHandler).Post(routers.SettingsPostHandler)

	//好友聊天
	jsock.JSock(m)

	u.Get("/friend/<kind>/<uid:(\\d+)>/", routers.GetFriendHandler).Post(routers.PostFriendHandler)

	//需要权限
	//-------------------------------------------------------------------------------------//
	//root routers
	r := g.Group("/root", setting.RootMiddler())

	//管理员注册请在models文件里设置,或使用默认管理员身份登录后台自行添加

	r.Get("/signout/", root.GetRSignoutHandler)
	r.Get("/dashboard/", root.GetRDashboardHandler)

	//主题
	r.Get("/theme/", root.GetRThemeHandler)
	r.Get("/theme/<name>/", root.GetRThemeHandler)

	rs := r.Group("/search")
	rs.Get("/user/", root.GetRReadUserHandler)
	rs.Get("/user/<keyword>/", root.GetRReadUserHandler)

	rs.Get("/category/", root.GetRSettingPasswordHandler)
	rs.Get("/node/", root.GetRSettingPasswordHandler)
	rs.Get("/topic/", root.GetRSettingPasswordHandler)

	//CRUD Create(创建)操作
	rc := r.Group("/create")
	rc.Get("/user/", root.GetRCreateUserHandler).Post(root.PostRCreateUserHandler)
	rc.Get("/category/", root.RCreateCategoryGetHandler).Post(root.RCreateCategoryPostHandler)
	rc.Get("/node/", root.RCreateNodeGetHandler).Post(root.PostRCreateNodeHandler)
	rc.Get("/topic/", root.GetRCreateTopicHandler).Post(root.PostRCreateTopicHandler)
	rc.Get("/<pid>/topic/", root.GetRCreateTopicHandler).Post(root.PostRCreateTopicHandler)
	rc.Get("/administrator/", root.GetRCreateRootHandler).Post(root.PostRCreateRootHandler)
	rc.Get("/page/", root.GetRCreatePageHandler).Post(root.PostRCreatePageHandler)
	rc.Get("/link/", root.GetRCreateLinkHandler).Post(root.PostRCreateLinkHandler)

	//Read(读取)操作
	rr := r.Group("/read")
	rr.Get("/user/", root.GetRReadUserHandler).Post(root.PostRReadUserHandler)
	rr.Get("/user/<uid>/", root.GetRReadUserHandler).Post(root.PostRReadUserHandler)
	rr.Get("/category/", root.GetRReadCategoryHandler)
	rr.Get("/category/<cid>/", root.GetRReadCategoryHandler)
	rr.Get("/node/", root.GetRReadNodeHandler)
	rr.Get("/node/<nid>/", root.GetRReadNodeHandler)

	rr.Get("/<pid>/topic/", root.GetRReadTopicHandler).Post(root.PostRReadTopicHandler)
	rr.Get("/<pid>/topic/<tid>/", root.GetRReadTopicHandler).Post(root.PostRReadTopicHandler)

	rr.Get("/topic/", root.GetRReadTopicHandler).Post(root.PostRReadTopicHandler)
	rr.Get("/topic/<tid>/", root.GetRReadTopicHandler).Post(root.PostRReadTopicHandler)

	rr.Get("/reply/", root.GetRReadReplyHandler).Post(root.PostRReadReplyHandler)
	rr.Get("/reply/<cmid>/", root.GetRReadReplyHandler).Post(root.PostRReadReplyHandler)

	rr.Get("/administrator/", root.GetRReadRootHandler)
	rr.Get("/administrator/<uid>/", root.GetRReadRootHandler)

	rr.Get("/report/", root.GetRReadReportHandler).Post(root.PostRReadReplyHandler)
	rr.Get("/report/<name>/", root.GetRReadReportHandler).Post(root.PostRReadReplyHandler)
	rr.Get("/report/<rid>/", root.GetRReadReportHandler).Post(root.PostRReadReplyHandler)

	rr.Get("/page/", root.GetRReadPageHandler)
	rr.Get("/page/<pageid>/", root.GetRReadPageHandler)

	rr.Get("/link/", root.GetRReadLinkHandler)
	rr.Get("/link/<linkid>/", root.GetRReadLinkHandler)

	//Update(更新)操作
	ru := r.Group("/update")
	ru.Get("/user/<uid>/<holder>/", root.GetRUpdateUserHandler)

	ru.Get("/user/<uid>/", root.GetRUpdateUserHandler).Post(root.PostRUpdateUserHandler)
	ru.Get("/administrator/<uid>/", root.GetRUpdateRootHandler).Post(root.PostRUpdateRootHandler)

	ru.Get("/category/", root.GetRUpdateCategoryHandler).Post(root.PostRUpdateCategoryHandler)
	ru.Get("/category/<cid>/", root.GetRUpdateCategoryHandler).Post(root.PostRUpdateCategoryHandler)

	ru.Get("/node/", root.GetRUpdateNodeHandler).Post(root.PostRUpdateNodeHandler)
	ru.Get("/node/<nid>/", root.GetRUpdateNodeHandler).Post(root.PostRUpdateNodeHandler)

	ru.Get("/topic/", root.GetRUpdateTopicHandler).Post(root.PostRUpdateTopicHandler)
	ru.Get("/topic/<tid>/", root.GetRUpdateTopicHandler).Post(root.PostRUpdateTopicHandler)

	ru.Get("/reply/", root.GetRUpdateReplyHandler).Post(root.PostRUpdateReplyHandler)
	ru.Get("/reply/<rid>/", root.GetRUpdateReplyHandler).Post(root.PostRUpdateReplyHandler)

	ru.Get("/page/", root.GetRUpdatePageHandler).Post(root.PostRUpdatePageHandler)
	ru.Get("/page/<pageid>/", root.GetRUpdatePageHandler).Post(root.PostRUpdatePageHandler)

	ru.Get("/link/", root.GetRUpdateLinkHandler).Post(root.PostRUpdateLinkHandler)
	ru.Get("/link/<linkid>/", root.GetRUpdateLinkHandler).Post(root.PostRUpdateLinkHandler)

	ru.Get("/topic/move/<tid>/", root.GetRMoveTopicHandler).Post(root.PostRMoveTopicHandler)
	ru.Get("/user/recharge/<uid>/", root.GetRRechargeUserHandler).Post(root.PostRRechargeUserHandler)

	//Delete(删除)操作
	rd := r.Group("/delete")
	rd.Get("/user/<uid>/", root.GetRDeleteUserHandler)
	rd.Get("/administrator/<uid>/", root.GetRDeleteUserHandler)
	rd.Get("/category/<cid>/", root.GetRDeleteCategoryHandler)
	rd.Get("/node/<nid>/", root.GetRDeleteNodeHandler)
	rd.Get("/topic/<tid>/", root.GetRDeleteTopicHandler)
	rd.Get("/reply/<rid>/", root.GetRDeleteReplyHandler)
	rd.Get("/page/<pageid>/", root.GetRDeletePageHandler)
	rd.Get("/link/<linkid>/", root.GetRDeleteLinkHandler)

	rse := r.Group("/setting")
	rse.Get("/password/", root.GetRSettingPasswordHandler).Post(root.PostRSettingPasswordHandler)

	//-------------------------------------------------------------------------------------//
	m.Get("/sock/", makross.WrapHTTPHandler(jsock.JSockHandler))
	//-------------------------------------------------------------------------------------//

	//需要权限

	//-------------------------------------------------------------------------------------//

	//APIS
	ra := m.Group("/api", setting.APISessionMiddler())
	ra.Get("/messages/", apis.GetMessagesHandler)

	//QINIU SIGN API
	//m.Get("/sign4qiniu/", apis.ApisSign4QiniuGetHandler)
	ra.Post("/upload/", apis.PostUploadHandler)

	//touch routers
	rat := ra.Group("/touch")
	rat.Get("/<type:([A-Za-z]+)>/<name:([A-Za-z]+)>/<id:\\d+>/", routers.GetTouchHandler)
	rat.Get("/favorite/<name:([A-Za-z]+)>/<id:\\d+>/", routers.GetFavoriteHandler)
	//rat.Get("/delete/topic/<tid>/",routers.TouchDeleteHandler)

	/*-------------APIs Start-------------*/

	v1 := m.Group("/apis/v1")
	//加密通信API
	// GET/POST /apis/v1/
	v1.Get("/", version1.GetApisHandler)
	v1.Post("/", setting.APICryptMiddler(), version1.PostApisHandler)
	//需要权限 针对移动端接口

	/*------------------------------------*/

	v2 := m.Group("/apis/v2")

	// GET from /apis/v2/
	v2.Get("/", version2.GetVersionHandler)

	// POST to /apis/v2/signup/
	v2.Post("/signup/", version2.PostSignupHandler)

	// POST to /apis/v2/signin/
	v2.Post("/signin/", version2.PostSigninHandler)

	// GET from /apis/v2/signout/
	v2.Get("/signout/", version2.GetSignoutHandler)

	// POST to /apis/v2/report/
	v2.Post("/report/", version2.PostReportHandler)

	var secret = "secret"
	jwt.DefaultJWTConfig.SigningKey = secret
	jwt.DefaultJWTConfig.Expires = time.Minute * 60
	j := v2.Group("", jwt.JWT(secret), auth.AuthJwtMiddler())

	// GET from /apis/v2/ping/
	j.Get("/ping/", version2.GetPongHandler)

	// GET from /apis/v2/categories/
	j.Get("/categories/", version2.GetCategoriesHandler)

	// GET from /apis/v2/category/<id>/
	j.Get("/category/<id:\\d+>/", version2.GetCategoryHandler)

	// GET,POST,PUT,DEL /apis/v2/category/
	j.Get("/category/", version2.GetCategoryHandler).Post(version2.PostCategoryHandler).Put(version2.PutCategoryHandler).Delete(version2.DelCategoryHandler)

	// GET from /apis/v2/nodes/
	j.Get("/nodes/", version2.GetNodesHandler)

	// GET from /apis/v2/node/<id>/
	j.Get("/node/<id:\\d+>/", version2.GetNodeHandler)

	// GET,POST,PUT,DEL /apis/v2/node/
	j.Get("/node/", version2.GetNodeHandler).Post(version2.PostNodeHandler).Put(version2.PutNodeHandler).Delete(version2.DelNodeHandler)

	// GET from /apis/v2/topics/
	j.Get("/topics/", version2.GetTopicsHandler)

	// GET from /apis/v2/topic/<id>/
	j.Get("/topic/<id:\\d+>/", version2.GetTopicHandler)

	// GET,POST /apis/v2/content/
	j.Get("/content/<id:\\d+>/", version2.GetContentHandler).Post(version2.PostContentHandler).Put(version2.PutContentHandler).Delete(version2.DelContentHandler)

	j.Get("/content/", version2.GetContentHandler).Post(version2.PostContentHandler).Put(version2.PutContentHandler).Delete(version2.DelContentHandler)

	// GET,POST /apis/v2/comment/
	j.Get("/comment/", version2.GetCommentHandler).Post(version2.PostCommentHandler)

	// GET from /apis/v2/users/
	j.Get("/users/", version2.GetUsersHandler)

	// GET/PUT /apis/v2/user/<id>/
	j.Get("/user/<id:\\d+>/", version2.GetUserHandler).Put(version2.PutUserHandler)

	// GET,POST,PUT,DEL /apis/v2/user/
	j.Get("/user/", version2.GetUserHandler).Post(version2.PostUserHandler).Put(version2.PutUserHandler).Delete(version2.DelUserHandler)

	// GET 特定用户话题数据 /apis/v2/user/<id>/topics/
	j.Get("/user/<id:\\d+>/topics/", version2.GetTopicsByUserHandler)

	// GET,POST /apis/v2/favorite/topic/
	j.Get("/favorite/topic/", version2.GetFavoriteTopicHandler).Post(version2.PostFavoriteTopicHandler).Delete(version2.DelFavoriteTopicHandler)

	// POST to /apis/v2/favorite/topic/<id>/
	j.Post("/favorite/topic/<id:\\d+>/", version2.PostFavoriteTopicHandler)

	// GET from /apis/v2/favorite/topics/
	j.Get("/favorite/topics/", version2.GetFavoriteTopicsHandler)

	// POST to /apis/v2/upload/
	j.Post("/upload/", version2.PostUploadHandler)

	/*------------------APIs End------------------*/

	if len(c.Args()) > 1 {
		switch l := len(c.Args()); {
		case l >= 2:
			{
				addr := fmt.Sprintf("%v:%v", c.Args()[0], c.Args()[1])
				log.Printf("Listen on %v\n", addr)
				m.Listen(addr)
				return nil
			}
		}

	} else {
		if len(c.Args()) == 1 {
			addr := c.Args()[0]
			log.Printf("Listen on %v\n", addr)
			m.Listen(addr)
			return nil
		} else {
			port := 8000
			log.Printf("Listen on %v\n", port)
			m.Listen(port)
			return nil
		}
	}

	return nil

}
