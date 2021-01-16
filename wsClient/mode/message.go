package mode

type NodeMessage struct {
	NodeName         string  `json:"nodename"`
	CPUavailability  float64 `json:"cpu"`
	MemAvailability  float64 `json:"mem"`
	DiskAvailability float64 `json:"disk"`
	NumNetwork       float64 `json:"network"`
	NodeStatus       bool    `json:"nodestatus"`
	NodeType         int     `json:"nodetype"`
	OwnCluser        string  `json:"cliser"`
}

const (
	NodeAvailable   = 1
	NodeUnAvailable = 2
)
