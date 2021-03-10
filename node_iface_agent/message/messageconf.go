package message

import (
	//"bytes"
	//"encoding/gob"
	"encoding/json"
	"fmt"
	"mecas/cmd/node_iface_agent/channel"
	"mecas/cmd/node_iface_agent/hostconf"
)

type cmdType uint8

const (
	TypeUnkown cmdType = iota
	ReportMessage
	SetSwitchONOrOFF
	TypeDelete
)

var errUnknownType = fmt.Errorf("unknown type")

//type (
//	CmdCreate struct {
//		A string
//	}
//	CmdUpdate struct {
//		B int
//	}
//	CmdDelete struct {
//		C bool
//	}
//)

type HostSend struct {
	Hostname string
	Conn     *channel.Connection
}

type Command struct {
	Type cmdType
	Raw  interface{}
}

func (c Command) MarshalJSON() (encoded []byte, err error) {
	var m struct {
		Type    cmdType         `json:"type"`
		Encoded json.RawMessage `json:"encoded"`
	}

	switch c.Raw.(type) {
	case hostconf.HostInfo, *hostconf.HostInfo:
		m.Type = ReportMessage
	case hostconf.NicReq, *hostconf.NicReq:
		m.Type = SetSwitchONOrOFF
	//case CmdDelete, *CmdDelete:
	//	m.Type = TypeDelete
	default:
		return nil, errUnknownType
	}
	encoded, err = json.Marshal(c.Raw)
	if err != nil {
		return
	}
	m.Encoded = json.RawMessage(encoded)
	return json.Marshal(m)
}

func (c *Command) UnmarshalJSON(data []byte) (err error) {
	var m struct {
		Type    cmdType         `json:"type"`
		Encoded json.RawMessage `json:"encoded"`
	}
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	switch c.Type = m.Type; c.Type {
	case ReportMessage:
		c.Raw = &hostconf.HostInfo{}
	case SetSwitchONOrOFF:
		c.Raw = &hostconf.NicReq{}
	//case TypeDelete:
	//c.Raw = &CmdDelete{}
	default:
		return errUnknownType
	}
	return json.Unmarshal(m.Encoded, c.Raw)
}
