package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var APPROOT string = "."

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
	} else {
		APPROOT = dir
	}
	log.Println("APPROOT: ", APPROOT)
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

	miscPath := "misc"
	_, err := os.Stat(miscPath)
	if err != nil {
		miscPath = filepath.Join(APPROOT, "misc")
	}
	LoadMiscDir(miscPath)

	jsonHinfo, err := json.Marshal(Hinfo)
	if err != nil {
		log.Println(err)
	} else {
		Echo("Final Result", string(jsonHinfo))

		err := ioutil.WriteFile(File, jsonHinfo, 0755)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("OK. Result was saved in: ", File)
		}
	}
}

func Echo(title string, t interface{}) {
	log.Println("=======", title, "=======")
	if Quiet != true {
		log.Println(t)
	}
}

func LoadMiscDir(pth string) error {
	Echo("LoadMiscDir", pth)
	err := filepath.Walk(pth, func(pth string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if filepath.Ext(pth) != ".txt" {
			return nil
		}

		if strings.Index(filepath.Base(pth), "gohostinfo") == -1 {
			return nil
		}
		log.Println("checking: ", pth)
		cnt, err := ioutil.ReadFile(pth)
		if err != nil {
			log.Println(err)
		}
		content := strings.ReplaceAll(string(cnt), "\r\n", "\n")
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			kv := strings.Split(line, "=")
			if len(kv) == 2 {
				if kv[0] != "" {
					log.Println(i, ")", kv[0], " = ", kv[1])
					k := strings.Trim(kv[0], " ")
					v := strings.Trim(kv[1], " ")
					if strings.Index(k, "#") == 0 {
						continue
					}
					Hinfo.Data[k] = v
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return err
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
		log.Fatal("cannot get the hostname, will abort")
	}

	Hinfo.ID = strings.ToLower(strings.Join([]string{"gohostinfo", h.Hostname}, "-"))
	Hinfo.Data["host"] = h

	log.Println(Hinfo.ID)
	Echo("Host", h)
}

func GetMemory() {
	vm, _ := mem.VirtualMemory()
	sm, _ := mem.SwapMemory()
	Hinfo.Data["virtual_memory"] = vm
	Hinfo.Data["swap_memory"] = sm
	Echo("VirtualMemory", vm)
	Echo("SwapMemory", sm)
}

func GetCPU() {
	c, _ := cpu.Info()
	Hinfo.Data["cpu"] = c
	Echo("CPU", c)
}

func GetDisk() {
	du := []*disk.UsageStat{}
	partitions, _ := disk.Partitions(true)
	for _, p := range partitions {
		if p.Mountpoint != "" {
			pu, err := disk.Usage(p.Mountpoint)
			if err != nil {
				log.Println(err)
			} else {
				du = append(du, pu)
			}

		}
	}
	Hinfo.Data["disk"] = du

	Echo("Disk", du)
}

func GetNet() {
	i, err := net.Interfaces()
	if err != nil {
		log.Println(err)
	}

	Hinfo.Data["net"] = i

	Echo("Net", i)
}

func GetDocker() {
	dkrs, err := docker.GetDockerStat()
	if err != nil {
		log.Println(err)
	}

	Hinfo.Data["docker"] = dkrs

	Echo("Docker", dkrs)
}
