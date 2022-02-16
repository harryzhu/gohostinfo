package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func init() {
	Hinfo = DefaultHostinfo()
	GetHost()
	GetCPU()
	GetMemory()
	GetDisk()

	json_hinfo, _ := json.Marshal(Hinfo)
	fmt.Println("=======")
	fmt.Println(string(json_hinfo))
	ioutil.WriteFile("hinfo.json", json_hinfo, 0755)
}

type Hostinfo struct {
	Id   string                 `json:"_id"`
	Rev  string                 `json:"_rev"`
	Data map[string]interface{} `json:"data"`
}

var Hinfo *Hostinfo

func DefaultHostinfo() *Hostinfo {
	hinfo := &Hostinfo{}
	hinfo.Id = "pls-set-your-key"
	data := make(map[string]interface{}, 10)
	data["cpu"] = ""
	hinfo.Data = data
	return hinfo
}

func PrintLine(t string) {
	fmt.Println("=======", t, "=======")
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
