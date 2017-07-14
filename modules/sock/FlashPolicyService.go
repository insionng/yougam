package sock

import (
	"bufio"
	"log"
	"net"
	"time"
)

/*
	Flash Policy Service必须在843端口提供服务
	所以若果防火墙在此端口设置过拦截的话，那么防火墙须要设置开放843端口
	sudo /sbin/iptables -I INPUT -p tcp --dport 843 -j ACCEPT
*/
func FlashPolicyService() error {
	hostPort := "0.0.0.0:843"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", hostPort)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}

}

func handleClient(conn net.Conn) {

	defer func() {
		time.Sleep(time.Second) //等待客户端响应之时间
		conn.Close()
	}()

	sendFlashPolicy(conn)

}

func sendFlashPolicy(conn net.Conn) {
	//Flash Player会先到请求的843端口请求策略文件的内容，
	//所以这里直接返回策略内容，而public目录不需要放置crossdomain.xml
	/*
			flashPolicy := `<?xml version="1.0"?>
		<!DOCTYPE cross-domain-policy SYSTEM "http://www.macromedia.com/xml/dtds/cross-domain-policy.dtd">
		<cross-domain-policy>
			<site-control permitted-cross-domain-policies="master-only"/>
			<allow-access-from domain="*" to-ports="*" />
		</cross-domain-policy>`
	*/
	FlashPolicy := `<?xml version="1.0"?>
<cross-domain-policy>
  <allow-access-from domain="*" to-ports="*" />
</cross-domain-policy>`
	writer := bufio.NewWriter(conn)
	writer.WriteString(FlashPolicy)
	writer.Flush()
}
