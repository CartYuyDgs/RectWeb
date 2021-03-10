package main

import (
	"os"
	"text/template"
)

type NicMode string

const (
	PF   NicMode = "pf"
	VF   NicMode = "vf"
	DPDK NicMode = "pfdpdk"
)

type HostNic struct {
	Host string   `json:"host"`
	Nics []string `json:"nics"`
}

type SrIov struct {
	CNIVersion string  `json:"cniVersion"`
	Type       string  `json:"type"`
	HostInfo   HostNic `json:"hostInfo"`
	DeviceID   string  `json:"deviceID"`
	Master     string  `json:"master"`
	Mode       NicMode `json:"mode"`
	LogFile    string  `json:"logFile"`
	LogLevel   string  `json:"logLevel"`
	Namespace  string  `json:"namespace"`
}

const metaCRD = `
apiVersion: "k8s.cni.cncf.io/v1"
kind: NetworkAttachmentDefinition
metadata:
 name: "{{.Type}}"	
 namespace: "{{.Namespace}}"
spec:
 config: '{
   "cniVersion": "{{.CNIVersion}}",
   "type": "{{.Type}}",
   "hostInfo": [{"host":{{.HostInfo.Host}}, "nics":{{.HostInfo.Nics}}}],
//    "deviceID": "0000:02:06.0",
//    "master": "enp11s0f1",
   "mode": "{{.Mode}}",
   "logFile": "{{.LogFile}}",
   "logLevel": "{{.LogLevel}}"
   }'

`

func createTemplate(info SrIov) (string, error) {
	// var result string
	tmpl, err := template.New("createmateCrd").Parse(metaCRD)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, info)
	if err != nil {
		panic(err)
	}

	return "", nil

}

func main() {
	var host = HostNic{
		Host: "host1",
		Nics: []string{"nic1", "nic2"},
	}

	var sriov = SrIov{
		CNIVersion: "0.3.0",
		Type:       "bronv",
		HostInfo:   host,
		Mode:       DPDK,
		LogFile:    "aaaa/aaaa",
		Namespace:  "abc",
		LogLevel:   "111",
	}

	createTemplate(sriov)
}
