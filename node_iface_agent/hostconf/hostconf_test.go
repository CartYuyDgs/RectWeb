package hostconf

import (
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	conf, err := GetConfig()
	if err != nil {
		t.Error(err)
	}

	if conf.HostName != "mec54" {
		t.Errorf("测试失败，名称失败")
	}

	fmt.Println(conf.NicFrequency, conf.Frequency, conf.HostName)
}


func TestMain(m *testing.M) {
	m.Run()
}
