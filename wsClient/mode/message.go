package mode

import "encoding/json"

type CmdMessage interface {
	printMessage()
	execMessage()
	returnMessage()
}

type OtherMessage struct {
	types int
	intM  json.RawMessage
}

func (M *ServerMessage) printMessage() {

}

func (M *ServerMessage) execMessage() {

}

func (M *ServerMessage) returnMessage() {

}

//type InforReport struct {
//	ReportTime string
//}

type InforNetworkModify struct {
	NodeName   string
	ReportTime string

	AttribBeMod      string
	AttribBeModValue string
	AttribModValue   string
}

type NodeConfig struct {
	NodeName   string
	ReportTime string
	Cmd        []string
	Args       []string
}

type ServerMessage struct {
	TypeMessage OtherMessage `json:"TypeMessage"` //指令类消息
	NodeConfigs NodeConfig
	NetworkInfo InforNetworkModify
}

//type NodeResult struct {
//	NodeName string
//	ReportTime string
//
//	Resu []string
//}

type NodeMessage struct {
	NodeName         string   `json:"nodename"`      //节点名称
	CPUavailability  float64  `json:"cpu"`           //CPU利用率
	MemAvailability  float64  `json:"mem"`           //内存使用率
	DiskAvailability float64  `json:"disk"`          //磁盘使用率
	NumNetwork       float64  `json:"network"`       //网卡数量
	NodeStatus       bool     `json:"nodestatus"`    //节点状态
	NodeType         int      `json:"nodetype"`      //节点类型
	OwnCluser        string   `json:"cliser"`        //所属集群
	MangleNetwork    string   `json:"managenetwork"` //管理网卡
	MangleIP         string   `json:"manageip"`      //管理IP
	NetworkNames     []string `json:"networknames"`  //网卡名称

	ReportTime string            `json:"reporttime"` //上报时间
	OtherInfo  map[string]string //其他状态消息，通断检测等
}

const (
	NodeAvailable   = 1
	NodeUnAvailable = 2
)
