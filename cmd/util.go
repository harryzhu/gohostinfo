package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func init() {
	Hinfo = DefaultHostinfo()

}

type Hostinfo struct {
	ID    string                 `json:"_id"`
	Time  int64                  `json:"time"`
	IDC   string                 `json:"idc"`
	Group string                 `json:"group"`
	Tags  []string               `json:"tags"`
	Data  map[string]interface{} `json:"data"`
}

var (
	Hinfo *Hostinfo
)

func DefaultHostinfo() *Hostinfo {
	hinfo := &Hostinfo{}
	hinfo.ID = "pls-set-your-key"
	hinfo.Time = time.Now().Unix()

	data := make(map[string]interface{}, 5)
	hinfo.Data = data
	return hinfo
}

func WalkOneByOne() {
	SetKey()
	SetTags()
	GetHost()
	GetCPU()
	GetMemory()
	GetDisk()
	GetNet()
	if WithDocker {
		GetDocker()
	}

	jsonHinfo, err := json.Marshal(Hinfo)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("======= Final Result =======")
		fmt.Println(string(jsonHinfo))
		err := ioutil.WriteFile(File, jsonHinfo, 0755)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK. Result was saved in: ", File)
		}
	}
}

func PrintLine(t string) {
	fmt.Println("=======", t, "=======")
}

func SetKey() {
	Hinfo.IDC = strings.ToLower(IDC)
	Hinfo.Group = strings.ToLower(Group)
}

func SetTags() {
	var arrTags []string
	if Tags != "" {
		Tags = strings.ReplaceAll(Tags, ",", ";")
		Tags = strings.ReplaceAll(Tags, " ", "-")
		Tags = strings.Trim(Tags, ";")

		tags := strings.Split(Tags, ";")
		for _, tag := range tags {
			if tag != "" {
				arrTags = append(arrTags, strings.ToLower(tag))
			}
		}
	}

	if len(arrTags) == 0 {
		arrTags = []string{}
	}

	Hinfo.Tags = arrTags
}

func GetHost() {
	h, err := host.Info()
	if err != nil {
		panic("cannot get the hostname, will abort")
	}

	Hinfo.ID = strings.ToLower(strings.Join([]string{"gohostinfo", h.Hostname}, "-"))
	Hinfo.Data["host"] = h

	PrintLine("Host")
	fmt.Println(h)
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
