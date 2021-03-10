package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/CartYuyDgs/RectWeb/node_agent/hostconf"
	"github.com/CartYuyDgs/RectWeb/node_agent/message"
	"github.com/CartYuyDgs/RectWeb/node_agent/onandoff"
	"log"

	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	Host = "127.0.0.1:6060"
	Path = "/v2/hosts/ws"
)

var (
	hostname string
	frequency time.Duration
	nicFrequency time.Duration
	identifier   string
	HostSend message.HostSend
	logFile		string
)

func init() {

	conf, err := hostconf.GetConfig()
	if err != nil {
		fmt.Println(err)
	}

	identifier = conf.HostName
	hostname = conf.HostName
	frequency = time.Duration(conf.Frequency)* time.Second
	nicFrequency = time.Duration(conf.NicFrequency)* time.Second
	logFile = conf.

	Hostname = hostname
	HostSend.Conn = ConnectionServer()
}



func main() {
	//conn := ConnectionServer()
	go sendStatistics(HostSend.Conn)
	//collectStatistics()
	receiveCommand(HostSend.Conn)
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
	log.Println("hostname: ", hostname, " ", "", Host)
	dialer := websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	WsConn, _, err := dialer.Dial(u.String(), header)

	if err != nil {
		log.Fatalln("connect err ", err)
	}
	// defer conn.Close()
	conn = CreateConnection(WsConn)
	conn.IsClosed = false
	return
}

type (
	Statistics struct {
	}
	//Command struct {
	//}
)

//func (Command) Run() {
//
//}

func collectStatistics() (stat []byte) {

	hostinfo := hostconf.HostInfo{}
	hostinfo.GetCPUInfo()
	hostinfo.GetMemInfo()
	hostinfo.GetDiskInfo()
	hostinfo.GetHostInfo()
	hostinfo.GetNetInfo()
	hostinfo.GetManageNic()
	hostinfo.GetSystemInfo()
	hostinfo.GetNetNum()

	for _, x := range hostinfo.Nic {
		x.GetLinkDetected()
		x.GetMacAdd()
		x.GetSpeedNic()
	}
	hostinfo.SetUpdateTime()

	//hostinfo.Print()
	var c = &message.Command{}
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
		var c = &message.Command{}
		if msg, err = conn.ReadMessage(); err != nil {
			return
		}
		err = c.UnmarshalJSON(msg)
		if err != nil {
			fmt.Println("read:", err)
			continue
		}

		switch c.Type {
		case message.SetSwitchONOrOFF:
			info, ok := (c.Raw).(*hostconf.NicReq)
			if !ok {
				fmt.Println(err)
				continue
			}
			ONOrOFFDiscover(conn, info.NicName, info.SwitchStatus)

		}

	}
}

func ONOrOFFDiscover(conn *Connection,name string, status int) {

	//整体时常5分钟，允许中断
	o := onandoff.OnOrOffDisc{}
	o.C = make(chan int, 1)

	for {
		select {
		case v, ok := <-o.C:
			if ok == true && v == 1 {
				o.ToStart(conn)
			} else if ok == true && v == 2 {
				o.ToStop()
			} else {
				continue
			}
		}
	}
}



