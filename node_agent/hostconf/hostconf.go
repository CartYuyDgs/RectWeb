package hostconf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	kernel_path      = "/proc/version"
	net_path         = "/sys/class/net/" ///sys/class也是构成linux统一设备模型的一部分
	virtual_net_path = "/sys/devices/virtual/net/"
	config_path      = "./../config.json" //此地址是一个暂时的配置文件存放地址，后续根据实际需要会改动
)

//var (
//	Hostname string
//	Frequency time.Duration
//	NicFrequency time.Duration
//	Identifier   string
//	LogFile		string
//)

type NicReq struct {
	HostId       string `json:"hostId"`
	NicName      string `json:"nicName"`
	SwitchStatus int    `json:"switchStatus"` //yes or no
}

//主机信息
type HostInfo struct {
	HostName   string     `json:"hostname"`
	HostID     string     `json:"hostid"`
	Cpu        string     `json:"cpu"`
	Mem        string     `json:"mem"`
	Disk       string     `json:"disk"`
	Os         string     `json:"os"`
	Kernel     string     `json:"kernel"`
	NicNum     int        `json:"nicnum"`
	ManageNic  string     `json:"managenic"`
	ManageIP   string     `json:"manageIP"`
	Nic        []*NicInfo `json:"nics"`
	UpdateTime string     `json:"updateTime"`
}

//网卡信息
type NicInfo struct {
	Pci           string `json:"pci"`
	NicDriver     string `json:"ncidriver"`
	DriverVersion string `json:"driverversion"`
	NicName       string `json:"nicname"`
	Bandwidth     string `json:"bandwidth"`
	NicModel      string `json:"nicmodel"`
	LinkStatus    string `json:"linkstatus"`
	MAC           string `json:"mac"`
}



type config struct {
	HostName     string `json:"host_name"`
	Frequency    int64    `json:"frequency"`
	NicFrequency int64    `json:"nic_frequency"`
	//Identifier   string
	LogFile		 string	`json:"log_file"`
}

func GetConfig() (*config, error) {
	file, err := os.Open(config_path)
	if err != nil {
		//log.Errorf("%+v", err)
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := config{}
	err = decoder.Decode(&conf)
	if err != nil {
		//log.Errorf("%+v", err)
		return nil, err
	}
	return &conf, nil
}

func (h HostInfo) Print() {
	fmt.Printf("------------print %s hostinfo-----------\n", h.HostName)
	fmt.Println()
	fmt.Println("HostName: ", h.HostName)
	fmt.Println("HostID: ", h.HostID)
	fmt.Println("CPU: ", h.Cpu)
	fmt.Println("Mem: ", h.Mem)
	fmt.Println("Disk: ", h.Disk)
	fmt.Println("OperationSystem: ", h.Os)
	fmt.Println("Kernel: ", h.Kernel)
	fmt.Println("NicNum: ", h.NicNum)
	fmt.Println("ManageNic: ", h.ManageNic)
	fmt.Println("ManageIP: ", h.ManageIP)
	fmt.Println("UpdateTime: ", h.UpdateTime)

	fmt.Printf("------------print %s Nicinfo-----------\n", h.HostName)
	for _, x := range h.Nic {
		fmt.Printf("-----------print nic %s------------\n", x.NicName)
		fmt.Println("Pci: ", x.Pci)
		fmt.Println("LinkStatus: ", x.LinkStatus)
		fmt.Println("MAC: ", x.MAC)
		fmt.Println("Bandwidth: ", x.Bandwidth)
		fmt.Println("DriverVersion: ", x.DriverVersion)
		fmt.Println("NciDriver: ", x.NicDriver)
		fmt.Println("NciModel: ", x.NicModel)
	}
}

func (h *HostInfo) SetUpdateTime() {
	h.UpdateTime = time.Now().String()
}

//cpu使用率
func (h *HostInfo) GetCPUInfo() {
	per, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalln("get cpuInfo error: ", err)
	}

	cpuPer := fmt.Sprintf("%f", per[0])
	//
	//fmt.Println("Cpu Per: ", cpuPer)
	h.Cpu = cpuPer
	return
}

//内存
func (h *HostInfo) GetMemInfo() {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalln("get memInfo error! ", err)
	}

	//fmt.Println("memInfo: ", memInfo.UsedPercent)
	memPer := fmt.Sprintf("%f", memInfo.UsedPercent)
	h.Mem = memPer
	return
}

//磁盘
func (h *HostInfo) GetDiskInfo() {
	// info, _ := disk.Partitions(true)
	// fmt.Println("info: ", info)

	useinfo, _ := disk.Usage("/")
	//fmt.Println("used: ", useinfo.UsedPercent)
	diskPer := fmt.Sprintf("%f", useinfo.UsedPercent)
	h.Disk = diskPer
	return
}

//主机名
func (h *HostInfo) GetHostInfo() {
	hostinfo, err := host.Info()
	if err != nil {
		log.Fatalln("get hostInfo error! ", err)
	}

	//fmt.Println("hostname: ", hostinfo.Hostname)
	//fmt.Printf("KernelVersion: %s-%s", hostinfo.KernelArch, hostinfo.KernelArch)
	//fmt.Printf("platform: %s-%s", hostinfo.Platform, hostinfo.PlatformVersion)
	hostname := hostinfo.Hostname

	// kernelinfo := fmt.Sprintf("%s-%s", hostinfo.KernelArch, hostinfo.KernelArch)
	// platform := fmt.Sprintf("%s-%s", hostinfo.Platform, hostinfo.PlatformVersion)

	h.HostName = hostname
	h.HostID = hostinfo.HostID
	// h.Kernel = kernelinfo
	// h.OperationSystem = platform
	return
}

func (n *NicInfo) GetLinkDetected() {
	LinkDetected := make([]byte, 20)

	nicCmd := exec.Command("ethtool", n.NicName)
	stdo := SetCommandStd(nicCmd)
	err := nicCmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	LinkCmd := exec.Command("grep", "Link detected")
	LinkCmd.Stdin = stdo

	linkstdo := SetCommandStd(LinkCmd)
	LinkCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	linkstdo.Read(LinkDetected)
	LinkDetected = bytes.Split(LinkDetected, []byte(": "))[1]
	//fmt.Println(strings.TrimSpace(string(LinkDetected)))
	linksta := strings.TrimSpace(string(LinkDetected))
	//linksta = strings.Split(linksta,":")[1]
	//fmt.Println(strings.Split(linksta,"\n"), len(linksta))
	n.LinkStatus = strings.TrimSpace(strings.Split(linksta, "\n")[0])
}

func (n *NicInfo) GetMacAdd() {
	MAC := make([]byte, 40)

	nicCmd := exec.Command("ifconfig", n.NicName)
	stdo := SetCommandStd(nicCmd)
	err := nicCmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	MacCmd := exec.Command("grep", "ether")
	MacCmd.Stdin = stdo
	linkstdo := SetCommandStd(MacCmd)
	MacCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	linkstdo.Read(MAC)
	//fmt.Println(strings.TrimSpace(string(MAC)))
	MACs := strings.TrimSpace(string(MAC))
	Mac := strings.Split(MACs, " ")

	//fmt.Println(MACs,Mac[1])
	//fmt.Println(strings.TrimSpace(string(MACs)))
	n.MAC = string(Mac[1])
}

func (n *NicInfo) GetSpeedNic() {
	CardSpeed := make([]byte, 30)

	//nicName := "ens2f0"
	nicCmd := exec.Command("ethtool", n.NicName)
	stdo := SetCommandStd(nicCmd)
	err := nicCmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	SpeedCmd := exec.Command("grep", "Speed")
	SpeedCmd.Stdin = stdo

	speedstdo := SetCommandStd(SpeedCmd)
	SpeedCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	speedstdo.Read(CardSpeed)
	CardSpeed = bytes.Split(CardSpeed, []byte(": "))[1]
	//fmt.Println(strings.TrimSpace(string(CardSpeed)))
	n.Bandwidth = strings.Split(strings.TrimSpace(string(CardSpeed)), "\n")[0]
}

func (h *HostInfo) GetManageNic() {
	manageNic := "ManageNic"
	manageIp := "ManageIP"
	Nic := os.Getenv(manageNic)
	if Nic == "" {
		fmt.Println("not find env ", manageNic)
		return
	}
	Ip := os.Getenv(manageIp)
	if Ip == "" {
		fmt.Println("not find env ", manageIp)
		return
	}
	//fmt.Println("manage nic is: ", Nic)
	//fmt.Println("manage ip is: ", Ip)

	h.ManageIP = Ip
	h.ManageNic = Nic
	return
}

func SetCommandStd(cmd *exec.Cmd) (stdout *bytes.Buffer) {
	stdout = &bytes.Buffer{}
	//stderr = &bytes.Buffer{}
	cmd.Stdout = stdout
	//cmd.Stderr = stderr
	return
}

func (h *HostInfo) GetNetNum() {

	PCIMap := make(map[string]string, 10)

	CardInfo := make([]byte, 1000)
	//
	// baseCom := "lspci -v | grep Ethernet"
	lspciCmd := exec.Command("lspci", "-v")
	stdo := SetCommandStd(lspciCmd)
	err := lspciCmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	ethCmd := exec.Command("grep", "Ethernet")
	ethCmd.Stdin = stdo

	//SpeedCmd := exec.Command("grep", "Speed")
	//SpeedCmd.Stdin = stdo

	ethstdo := SetCommandStd(ethCmd)
	ethCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	ethstdo.Read(CardInfo)

	sep := []byte("\n")
	res := bytes.Split(bytes.TrimSpace(CardInfo), sep)
	for index, value := range res {
		if index >= len(res)-1 {
			break
		}
		res := bytes.Split(value, []byte(" Ethernet controller: "))
		PCIMap[string(res[0])] = string(res[1])
	}

	//fmt.Println(PCIMap, len(PCIMap))
	h.NicNum = len(PCIMap)

	for _, v := range h.Nic {
		v.NicModel = PCIMap[v.Pci[5:]]
	}
}

func (h *HostInfo) GetSystemInfo() {

	//cat /proc/version  路径说明
	f, err := ioutil.ReadFile(kernel_path)
	if err != nil {
		fmt.Println(err)
	}
	//Linux version .*-generic
	//(gcc version 7.5.0 (.*))#

	infos := string(f)

	reg_k := regexp.MustCompile("Linux version .*-generic")
	if reg_k == nil {
		fmt.Println("regexp error!")
		return
	}

	kernel_info := reg_k.FindStringSubmatch(infos)[0]
	//fmt.Println("kernel info :", kernel_info)

	reg_s := regexp.MustCompile("(gcc version 7.5.0 (.*?))#")
	if reg_s == nil {
		fmt.Println("regexp error!")
		return
	}

	system_infos := reg_s.FindString(infos)
	system_info := strings.Split(strings.Split(system_infos, `(`)[1], `)`)[0]
	//fmt.Println("system info :", system_info)
	h.Kernel = kernel_info
	h.Os = system_info

}

func judge(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

//网卡信息
func (h *HostInfo) GetNetInfo() {
	//basePath := "/sys/class/net/"
	cmd := exec.Command("ls", net_path)
	buf, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	output := string(buf)
	// fmt.Println("output: ",output)

	//basePath2 := "/sys/devices/virtual/net/"
	cmd2 := exec.Command("ls", virtual_net_path)
	buf2, err := cmd2.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	output2 := string(buf2)
	// fmt.Println("output2: ",output2)
	netList := strings.Split(output2, "\n")
	// fmt.Println("net: ",netList)

	for _, device := range strings.Split(output, "\n") {
		if judge(device, netList) {

		} else {
			info, err := getNetDetail(device)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			//fmt.Println("net: ", device, " infos: ", info)
			nic := NicInfo{}
			nic.NicName = device
			nic.NicDriver = info["driver"]
			nic.DriverVersion = info["version"]
			nic.Pci = info["bus-info"]
			h.Nic = append(h.Nic, &nic)
		}
	}

}

func getNetDetail(net string) (map[string]string, error) {
	baseCmd := "ethtool"
	_, err := exec.LookPath("ethtool")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	//fmt.Println(path, err)

	cmds := exec.Command(baseCmd, "-i", net)

	//fmt.Println(cmds)
	buf, err := cmds.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	output := string(buf)
	NicDetails := make(map[string]string)
	for _, device := range strings.Split(output, "\n") {
		//fmt.Println(device)

		infos := strings.Split(device, ": ")

		switch infos[0] {
		case "driver":
			NicDetails["driver"] = infos[1]
		case "version":
			NicDetails["version"] = infos[1]
		case "bus-info":
			NicDetails["bus-info"] = infos[1]
		default:
			//fmt.Println(infos)
		}

	}
	//fmt.Println("-----------------------")
	//fmt.Println(NicDetails)
	return NicDetails, nil
}

func CreateCTX() (context.Context, context.CancelFunc) {

	timeout := 5 * time.Second
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	return ctx, cancle
}

func NicStatusUp(nic string) (LinkStatus string, err error) {
	cmds := exec.Command("ifconfig", nic, "up")
	_, err = cmds.Output()
	if err != nil {
		return "", err
	}

	//fmt.Println(LinkDetectedget(nic))
	return LinkDetectedget(nic), nil
}

func NicStatusDown(nic string) (LinkStatus string, err error) {
	cmds := exec.Command("ifconfig", nic, "down")
	_, err = cmds.Output()
	if err != nil {
		return "", err
	}

	//fmt.Println(LinkDetectedget(nic))
	return LinkDetectedget(nic), nil
}

func LinkDetectedget(NciName string) (LinkStatus string) {
	LinkDetected := make([]byte, 20)

	nicCmd := exec.Command("ethtool", NciName)
	stdo := SetCommandStd(nicCmd)
	err := nicCmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	LinkCmd := exec.Command("grep", "Link detected")
	LinkCmd.Stdin = stdo

	linkstdo := SetCommandStd(LinkCmd)
	LinkCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	linkstdo.Read(LinkDetected)
	LinkDetected = bytes.Split(LinkDetected, []byte(": "))[1]
	linksta := strings.TrimSpace(string(LinkDetected))
	LinkStatus = strings.TrimSpace(strings.Split(linksta, "\n")[0])
	return LinkStatus
}
