package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/winservices"
)

func init() {
	Hinfo = DefaultHostinfo()

}

type Hostinfo struct {
	Id   string                 `json:"_id"`
	Rev  string                 `json:"_rev"`
	Time int64                  `json:"time"`
	IDC  string                 `json:"idc"`
	Data map[string]interface{} `json:"data"`
}

var Hinfo *Hostinfo

func DefaultHostinfo() *Hostinfo {
	hinfo := &Hostinfo{}
	hinfo.Id = "pls-set-your-key"
	hinfo.Time = time.Now().Unix()

	data := make(map[string]interface{}, 10)
	data["cpu"] = ""
	hinfo.Data = data
	return hinfo
}

func WalkOneByOne() {
	SetIDC()
	GetHost()
	GetCPU()
	GetMemory()
	GetDisk()
	GetNet()
	GetDocker()
	GetWinservices()

	json_hinfo, err := json.Marshal(Hinfo)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("======= Final Result =======")
		fmt.Println(string(json_hinfo))
		ioutil.WriteFile("gohostinfo.json", json_hinfo, 0755)
	}
}

func PrintLine(t string) {
	fmt.Println("=======", t, "=======")
}

func SetIDC() {
	Hinfo.IDC = IDC
}

func GetHost() {
	h, _ := host.Info()
	PrintLine("Host")
	fmt.Println(h)

	Hinfo.Id = h.Hostname
	Hinfo.Data["host"] = h
}

func GetMemory() {
	vm, _ := mem.VirtualMemory()
	sm, _ := mem.SwapMemory()
	Hinfo.Data["virtual_memory"] = vm
	Hinfo.Data["swap_memory"] = sm
	PrintLine("VirtualMemory")
	fmt.Println(vm)
	PrintLine("SwapMemory")
	fmt.Println(sm)
}

func GetCPU() {
	c, _ := cpu.Info()
	Hinfo.Data["cpu"] = c
	PrintLine("CPU")
	fmt.Println(c)
}

func GetDisk() {
	du := []*disk.UsageStat{}
	partitions, _ := disk.Partitions(true)
	for _, p := range partitions {
		//fmt.Println(p)
		if p.Mountpoint != "" {
			pu, err := disk.Usage(p.Mountpoint)
			if err != nil {
				fmt.Println(err)
			} else {
				du = append(du, pu)
			}

		}
	}
	Hinfo.Data["disk"] = du
	PrintLine("Disk")
	fmt.Println(du)
}

func GetNet() {
	i, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}

	Hinfo.Data["net"] = i
	PrintLine("Net")
	fmt.Println(i)
}

func GetDocker() {
	dkrs, err := docker.GetDockerStat()
	if err != nil {
		fmt.Println(err)
	}

	Hinfo.Data["docker"] = dkrs
	PrintLine("Docker")
	fmt.Println(dkrs)
}

func GetWinservices() {
	w, err := winservices.ListServices()
	if err != nil {
		fmt.Println(err)
	}

	Hinfo.Data["winservice"] = w
	PrintLine("Winservice")
	fmt.Println(w)
}
