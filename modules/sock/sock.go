package sock

import (
	//"errors"
	"fmt"
	"github.com/insionng/makross"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"github.com/insionng/yougam/helper"
	"github.com/insionng/yougam/libraries/gorilla/websocket"
	"github.com/insionng/yougam/models"

	"github.com/insionng/yougam/modules/setting"
	"github.com/insionng/yougam/routers"
)

type Client struct {
	websocket *websocket.Conn
	clientIP  net.Addr
	userid    int64
}

type Box struct {
	sync.RWMutex
	key string
	box []*Client
}

type Boxes struct {
	sync.RWMutex
	boxes []*Box
}

func (b *Box) appendClient(client *Client) {
	b.Lock()
	b.box = append(b.box, client)
	for _, c := range b.box {
		if c != client {
			log.Println(client.clientIP, "进入聊天~")
		}
	}
	b.Unlock()
}

func (b *Box) removeClient(client *Client) {
	b.Lock()
	defer b.Unlock()

	for index, c := range b.box {
		if c == client {
			b.box = append(b.box[:index], b.box[(index+1):]...)
			if err := c.websocket.WriteMessage(1, []byte(`<div class="alert alert-warning alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="alert-heading">警告</h4>`+`你的账号在别处连接，当前通信已经失效，如非本人操作请立即修改账号密码..`+`</div>`)); err != nil {
				return
			}
		} else {
			if err := client.websocket.WriteMessage(1, []byte(`<div class="alert alert-warning alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="alert-heading">警告</h4>`+`你的账号已经重新连接，原连接已经失效，系统禁止同一账号在多客户端同时连接..`+`</div>`)); err != nil {
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
			if err := c.websocket.WriteMessage(1, []byte(`<div class="alert alert-warning alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="alert-heading">警告</h4>`+`你的账号在别处连接，当前通信已经失效，如非本人操作请立即修改账号密码..`+`</div>`)); err != nil {
				return
			}
		} else {
			if err := client.websocket.WriteMessage(1, []byte(`<div class="alert alert-warning alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="alert-heading">警告</h4>`+`你的账号已经重新连接，原连接已经失效，系统禁止同一账号在多客户端同时连接..`+`</div>`)); err != nil {
				return
			}
		}
	}
}

func (b *Box) broadcastMessage(messageType int, message []byte) {
	b.Lock()
	defer b.Unlock()

	for _, client := range b.box {
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
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

func Sock(m *makross.Macross) {

	if setting.FlashPolicyService {
		go FlashPolicyService()
	}

	var boxes *Boxes
	boxes = newBoxes()

	m.Get("/sock/connect/:name([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)/", func(self *self.Context) {
		if name := self.ParamsEscape("name"); self.IsSigned && (len(name) > 0) && websocket.IsWebSocketUpgrade(self.Req.Request) {
			recipient, err := models.GetUserByUsername(name)
			if (err != nil) || (recipient == nil) {
				log.Println(err)
				return // err
			}

			ws, err := websocket.Upgrade(self.Resp, self.Req.Request, nil, 1024, 1024)
			if _, okay := err.(websocket.HandshakeError); okay {
				log.Println("Not a websocket handshake")
				http.Error(self.Resp, "Not a websocket handshake", 400)
				return // err
			} else if err != nil {
				log.Println(err)
				return // err
			}

			rAddr := ws.RemoteAddr()
			sockCli := &Client{ws, rAddr, self.User.Id}

			box := boxes.getBox(helper.OrderKey(self.User.Id, recipient.Id))

			if len(box.box) < 2 {
				log.Println("L<2")
				box.appendClient(sockCli)
			} else {
				log.Println("L>=2")
				uid := fmt.Sprintf("%v", self.User.Id)
				//如果是原来的用户则更新连接地址,即禁止用户同时登录2个以上客户端
				if s := strings.Split(box.key, ":"); (s[0] == uid) || (s[1] == uid) {
					box.removeUser(sockCli)
					box.appendClient(sockCli)
				} else {
					return // nil
				}

			}

			for {

				messageType, message, err := ws.ReadMessage()
				if err != nil {
					box.removeClient(sockCli)
					//log.Println("bye")
					log.Println(err)
					return // err
				}

				now := time.Now()
				policy := helper.ObjPolicy()
				s := policy.SanitizeBytes(message)
				if len(self.User.AvatarMedium) <= 0 {
					self.User.AvatarMedium = "/img/d48.png"
				}

				body := fmt.Sprintf(`
                    <article class="panel panel-warning box cell first message">
                        <div class="panel-body">%s</div>
                        <div id="uid%v" class="panel-footer box nobg cell last">
				            <a href="/user/%s/" class="thumb avatar">
				            	<p><img src="%s"></p><p>%s</p>
				            </a>
	                        <div style="clear:both;"></div>
	                        <small class="pull-right">%s</small>
                        </div>
                        <div style="clear:both;"></div>
                    </article>`, s,
					self.User.Id, self.User.Username,
					self.User.AvatarMedium, self.User.Username,
					now.Format("2006.01.02 03:04"))
				box.broadcastMessage(messageType, []byte(body))

			}
		} else {
			return // errors.New("Is not WebSocket Upgrade")
		}
	})

	m.Get("/contact/", routers.GetContactHandler)
	m.Get("/contact/search/", routers.GetContactHandler)
	m.Get("/connect/<name:([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)>/", routers.GetContactHandler)

}
