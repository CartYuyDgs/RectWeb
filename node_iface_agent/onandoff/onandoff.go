package onandoff

import (
	"context"
	chans "mecas/cmd/node_iface_agent/channel"
	"mecas/cmd/node_iface_agent/hostconf"
	"mecas/pkg/errors"
	"time"
)

type OnOrOffDisc struct {
	NicName string
	Cancle  context.CancelFunc
	C       chan int
	Ctx     context.Context
}

func (o *OnOrOffDisc) ToStart(conn *chans.Connection) error {
	o.Ctx, o.Cancle = hostconf.CreateCTX()
	stat, err := hostconf.NicStatusUp(o.NicName)
	if err != nil || stat != "up" {
		return errors.New("up nic status errors")
	}

	go o.SendOnOrOffDiscoverInfo(conn)

	for {
		select {
		case <-o.Ctx.Done():
			o.C <- 2
			//stat, err := NicStatusDown(o.NicName)
			//if err != nil || stat != "up" {
			//	return errors.New("up nic status errors")
			//}
			return nil
		}
	}
}

func (o *OnOrOffDisc) SendOnOrOffDiscoverInfo(conn *chans.Connection) {
	ticker := time.NewTicker(hostconf.NicFrequency)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			stat := hostconf.LinkDetectedget(o.NicName)
			conn.WriteMessage([]byte(stat))
		}
	}
}

func (o *OnOrOffDisc) ToStop() error {
	stat, err := hostconf.NicStatusDown(o.NicName)
	if err != nil || stat != "up" {
		return errors.New("up nic status errors")
	}
	o.Cancle()
	return nil
}
