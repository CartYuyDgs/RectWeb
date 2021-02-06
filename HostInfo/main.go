package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"log"
	"time"
)

func main() {
	//getCPUInfo()
	getMemInfo()
	getHostInfo()
	getDiskInfo()
	getNetInfo()
}

//cpu使用率
func getCPUInfo() {
	//cpuInfo, err := cpu.Info()
	//if err != nil {
	//	log.Fatalln("get cpuInfo error! ", err)
	//}
	//fmt.Println("Cpu Info: ",cpuInfo)

	per, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalln("get cpuInfo error! ", err)
	}

	fmt.Println("Cpu Per: ", per)
}

//内存
func getMemInfo() {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalln("get memInfo error! ", err)
	}

	fmt.Println("memInfo: ", memInfo.UsedPercent)
}

//磁盘
func getDiskInfo() {
	info, _ := disk.Partitions(true)
	fmt.Println("info: ", info)

	useinfo, _ := disk.Usage("D:")
	fmt.Println("used: ", useinfo.UsedPercent)
}

//主机名
func getHostInfo() {
	hostinfo, err := host.Info()
	if err != nil {
		log.Fatalln("get hostInfo error! ", err)
	}

	fmt.Println("hostname: ", hostinfo.Hostname)
	fmt.Println("hostInfo: ", hostinfo)
}

//网卡信息
func getNetInfo() {
	inters, _ := net.Interfaces()
	fmt.Println("interface: ", inters)
	for _, x := range inters {
		fmt.Println(x.Name, x.Flags)
		fmt.Println(x.Addrs)
		fmt.Println(x.HardwareAddr)
	}

	//info, _ :=net.IOCounters(true)
	//for index, v := range info {
	//	fmt.Printf("%v:%v send:%v recv:%v\n", index, v, v.BytesSent, v.BytesRecv)
	//
	//
	//}

}
