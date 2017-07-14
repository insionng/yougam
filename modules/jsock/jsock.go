package jsock

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/libraries/igm/sockjs-go.v2/sockjs"
	"github.com/insionng/yougam/models"
	"github.com/insionng/yougam/modules/setting"
	"github.com/insionng/yougam/routers"

	"github.com/insionng/makross"
)

//Client 客户端
type Client struct {
	session  sockjs.Session
	clientIP string
	userid   int64
}

//Box 盒子
type Box struct {
	sync.RWMutex
	key string
	box []*Client
}

//Boxes 盒子集合
type Boxes struct {
	sync.RWMutex
	boxes []*Box
}

var boxes = newBoxes()

func (b *Box) appendClient(client *Client) {
	b.Lock()
	defer b.Unlock()
	b.box = append(b.box, client)
	/*
		for _, c := range b.box {
			if c != client {
				log.Println(client.clientIP, "进入聊天~")
			}
		}
	*/

}

func (b *Box) removeClient(client *Client) {
	b.Lock()
	defer b.Unlock()

	for index, c := range b.box {
		if c == client {
			b.box = append(b.box[:index], b.box[(index+1):]...)
			if err := c.session.Send(`<div class="alert alert-warning alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="alert-heading">警告</h4>` + `你的账号在别处连接，当前通信已经失效，如非本人操作请立即修改账号密码..` + `</div>`); err != nil {
				return
			}
		} else {
			if err := client.session.Send(`<div class="alert alert-warning alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="alert-heading">警告</h4>` + `你的账号已经重新连接，原连接已经失效，系统禁止同一账号在多客户端同时连接..` + `</div>`); err != nil {
				return
			}
		}
	}
}

func (b *Box) removeUser(client *Client) {
	b.Lock()
	defer b.Unlock()

	for index, c := range b.box {
		if c.userid == client.userid {
			b.box = append(b.box[:index], b.box[(index+1):]...)
			if err := c.session.Send(`<div class="alert alert-warning alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="alert-heading">警告</h4>` + `你的账号在别处连接，当前通信已经失效，如非本人操作请立即修改账号密码..` + `</div>`); err != nil {
				return
			}
		} else {
			if err := client.session.Send(`<div class="alert alert-warning alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="alert-heading">警告</h4>` + `你的账号已经重新连接，原连接已经失效，系统禁止同一账号在多客户端同时连接..` + `</div>`); err != nil {
				return
			}
		}
	}
}

func tpl(align, username, avatar, datetime, message string) string {

	var bg = "b-light"
	if align == "right" {
		bg = "bg-light"
	}

	if len(avatar) == 0 {
		avatar = "/identicon/" + username + "/48/default.png"
	}

	var tpl = `
                <article class="chat-item ` + align + `">
                    <a href="/user/` + username + `/" class="pull-` + align + ` thumb-sm avatar">
                    	<img src="` + avatar + `" alt="` + username + `"/>
                    	<div class="text-ellipsis text-center">` + username + `</div>
                    </a>
                    <section class="chat-body">
                        <div class="panel ` + bg + ` text-sm m-b-none">
                            <div class="panel-body">
                                <span class="arrow ` + align + `"></span>
                                <div class="m-b-none">` + message + `</div>
                                <div class="clear"></div>
                            </div>
                        </div>
                        <small class="text-muted"> <i class="fa fa-ok text-success"></i>
                            ` + datetime + `
                        </small>
                    </section>
                </article>
`
	return tpl
}

func (b *Box) broadcastMessage(currentUserid int64, username, avatar, datetime, message string) {
	b.Lock()
	defer b.Unlock()

	var align string
	for _, client := range b.box {

		if client.userid == currentUserid {
			align = "right"
		} else {
			align = "left"
		}

		if err := client.session.Send(tpl(align, username, avatar, datetime, message)); err != nil {
			return
		}
	}
}

func (b *Box) broadcastMessageTo(isMyself bool, currentUserid int64, username, avatar, datetime, message string) {
	b.Lock()
	defer b.Unlock()

	var align string
	for _, client := range b.box {

		if client.userid == currentUserid {
			align = "right"
		} else {
			align = "left"
		}

		if isMyself {
			if client.userid == currentUserid {
				if err := client.session.Send(tpl(align, username, avatar, datetime, message)); err != nil {
					return
				}
			}
		} else {
			if client.userid != currentUserid {
				if err := client.session.Send(tpl(align, username, avatar, datetime, message)); err != nil {
					return
				}
			}
		}

	}
}

func (b *Boxes) getBox(key string) *Box {
	b.Lock()
	defer b.Unlock()

	for _, Box := range b.boxes {
		if Box.key == key {
			return Box
		}
	}

	box := &Box{sync.RWMutex{}, key, make([]*Client, 0)}
	b.boxes = append(b.boxes, box)
	return box
}

func newBoxes() *Boxes {
	return &Boxes{sync.RWMutex{}, make([]*Box, 0)}
}

func sockHandler(session sockjs.Session) {
	l, err := url.Parse(session.Request().RequestURI)
	if err != nil {
		session.Close(1008, "Bad Request")
		return
	}

	sender := l.Query()["sender"][0]
	receiver := l.Query()["receiver"][0]
	token := l.Query()["token"][0]

	if (sender == receiver) || (len(sender) == 0) || (len(receiver) == 0) || (len(token) != 48) {
		session.Close(1008, "Bad Request")
		return
	}

	if len(receiver) > 0 {
		recipient, err := models.GetUserByUsername(receiver)
		if (err != nil) || (recipient == nil) {
			log.Println(err)
			session.Close(1008, "Bad Request")
			return
		}

		if key := sender + ":" + helper.AesKey + ":" + receiver; !helper.ValidateHash(token, key) {
			session.Close(1008, "Bad Request")
			return
		}

		var me *models.User
		item := setting.Cache.Get(token)
		if item == nil {
			session.Close(1008, "Bad Request")
			return
		}

		if user, okay := item.Value().(*models.User); okay {

			usr, e := models.GetUserByUsername(sender)
			if (e != nil) || (usr == nil) {
				session.Close(1011, "Unauthorized")
				return
			}

			if (usr.Password != user.Password) || (usr.Username != user.Username) {
				session.Close(1008, "Bad Request")
				return
			}

			if !models.IsFriend(usr.Id, recipient.Id) {
				session.Close(1008, "Bad Request")
				return
			}

			setting.Cache.Delete(token) //用完即弃
			me = usr

		} else {
			session.Close(1008, "Bad Request")
			return
		}

		rAddr := session.Request().RemoteAddr
		sockCli := &Client{session, rAddr, me.Id}
		orderKey := helper.OrderKey(me.Id, recipient.Id)
		box := boxes.getBox(orderKey)

		if len(box.box) < 2 {
			box.appendClient(sockCli)
		} else {
			uid := fmt.Sprintf("%v", me.Id)
			//如果是原来的用户则更新连接地址,即禁止用户同时登录2个以上客户端
			if s := strings.Split(box.key, ":"); (s[0] == uid) || (s[1] == uid) {
				box.removeUser(sockCli)
				box.appendClient(sockCli)
			} else {
				session.Close(1008, "Bad Request")
				return
			}

		}

		//读取并发送离线消息
		messages, e := models.GetMessagesViaReceiverWithSender(0, 0, me.Username, receiver, "asc")
		if (e == nil) && (messages != nil) {
			for _, v := range *messages {
				models.DelMessage(v.Id) //阅后即焚
				box.broadcastMessageTo(false, v.Uid, v.Sender, v.Avatar, helper.TimeSince(v.Created), v.Content)
			}
		}

		for {

			//若果双方好友关系已经解除
			if !models.IsFriend(me.Id, recipient.Id) {
				session.Close(1008, "Bad Request")
				return
			}

			message, err := session.Recv()
			if err != nil {
				box.removeClient(sockCli)
				session.Close(1000, "Closed Request")
				return
			}

			policy := helper.ObjPolicy()
			body := policy.Sanitize(message)
			now := time.Now()
			box.broadcastMessage(me.Id, me.Username, me.AvatarMedium, now.Format("2006-01-02 03:04"), body)
			go func() {
				//若果对方没有连线
				if box := boxes.getBox(orderKey); len(box.box) == 1 {
					m := new(models.Message)
					m.Key = orderKey
					m.Uid = me.Id
					m.Sender = me.Username
					m.Avatar = me.AvatarMedium
					m.Receiver = receiver
					m.Content = body
					m.Created = now.Unix()
					models.PostMessage(m)
				}

				//always save the history message
				m := new(models.HistoryMessage)
				m.Key = orderKey
				m.Uid = me.Id
				m.Sender = me.Username
				m.Avatar = me.AvatarMedium
				m.Receiver = receiver
				m.Content = body
				m.Created = now.Unix()
				models.PostHistoryMessage(m)
			}()

		}
	} else {
		session.Close(1008, "Bad Request")
		return
	}

}

//JSockHandler 对象
var JSockHandler = sockjs.NewHandler("/sock", sockjs.Options{
	Websocket:       true,
	JSessionID:      nil,
	SockJSURL:       fmt.Sprintf("http://%s/%s", helper.Domain, "/libs/sockjs-client-1.1.0/sockjs.min.js"),
	HeartbeatDelay:  25 * time.Second,
	DisconnectDelay: 5 * time.Second,
	ResponseLimit:   128 * 1024,
}, sockHandler)

//JSock 路由
func JSock(m *makross.Makross) {
	m.Get("/contact/", routers.GetContactHandler)
	m.Get("/contact/search/", routers.GetContactHandler)
	m.Get("/connect/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/", routers.GetConnectHandler)
}
