package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	Host = "127.0.0.1:6060"
	Path = "/v2/hosts/ws"
)

var (
	identifier   string
	frequency    = 5 * time.Second
	nicFrequency = 1 * time.Second
	hostname     = "mec54"

	hostSend HostSend
)

func init() {

	name, err := os.Hostname()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	identifier = name
	hostname = name

	hostSend.hostname = name
	hostSend.conn = ConnectionServer()
}

func main() {
	//conn := ConnectionServer()
	go sendStatistics(hostSend.conn)
	//collectStatistics()
	receiveCommand(hostSend.conn)
}

func ConnectionServer() (conn *Connection) {
	u := url.URL{
		Scheme: "ws",
		Host:   Host,
		Path:   Path,
	}
	// NOTE: 节点标识在这里传递，否则Server只能根据TCPAddr.IP来确定对方的身份
	header := make(http.Header, 4)
	header.Set("WebSocketAgent", hostname)
	var username, password string
	username = "bonc"
	password = "bonc@123"
	header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
	//WsConn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	//WsConn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	//WsConn, _, err :=websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}.Dial(u.String(), header)
	log.Println("hostname: ", hostname, " ", "", Host)
	dialer := websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	WsConn, _, err := dialer.Dial(u.String(), header)

	if err != nil {
		log.Fatalln("connect err ", err)
	}
	// defer conn.Close()
	conn = CreateConnection(WsConn)
	conn.isClosed = false
	return
}

type (
	Statistics struct {
	}
	//Command struct {
	//}
)

func (Command) Run() {

}

func collectStatistics() (stat []byte) {

	hostinfo := HostInfo{}
	hostinfo.getCPUInfo()
	hostinfo.getMemInfo()
	hostinfo.getDiskInfo()
	hostinfo.getHostInfo()
	hostinfo.getNetInfo()
	hostinfo.GetManageNic()
	hostinfo.getSystemInfo()
	hostinfo.getNetNum()

	for _, x := range hostinfo.Nic {
		x.GetLinkDetected()
		x.GetMacAdd()
		x.GetSpeedNic()
	}
	hostinfo.setUpdateTime()

	//hostinfo.Print()
	var c = &Command{}
	c.Raw = &hostinfo
	//stat, err := json.Marshal(hostinfo)
	stat, err := c.MarshalJSON()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return stat
}

func sendStatistics(conn *Connection) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			stat := collectStatistics()
			conn.WriteMessage(stat)
		}
	}
}

// 接收服务端发回的指令，并执行指令
func receiveCommand(conn *Connection) {
	var (
		err error
		msg []byte
	)
	for {
		var c = &Command{}
		if msg, err = conn.ReadMessage(); err != nil {
			return
		}
		err = c.UnmarshalJSON(msg)
		if err != nil {
			fmt.Println("read:", err)
			continue
		}

		switch c.Type {
		case SetSwitchONOrOFF:
			info, ok := (c.Raw).(*NicReq)
			if !ok {
				fmt.Println(err)
				continue
			}
			conn.ONOrOFFDiscover(info.NicName, info.SwitchStatus)

		}

	}
}

func (conn *Connection) ONOrOFFDiscover(name string, status int) {

	//整体时常5分钟，允许中断
	o := OnOrOffDisc{}
	o.c = make(chan int, 1)

	for {
		select {
		case v, ok := <-o.c:
			if ok == true && v == 1 {
				o.toStart(conn)
			} else if ok == true && v == 2 {
				o.toStop()
			} else {
				continue
			}
		}
	}
}
