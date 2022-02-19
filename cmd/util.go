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

var (
	APPROOT string = "."
	Hinfo   *Hostinfo
)

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
	} else {
		APPROOT = dir
	}

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

func DefaultHostinfo() *Hostinfo {
	return &Hostinfo{
		ID:    "pls-set-your-key",
		Time:  time.Now().Unix(),
		IDC:   "",
		Group: "",
		Tags:  []string{},
		Data:  make(map[string]interface{}, 5),
	}
}

func WalkOneByOne() {
	Echo("APPROOT", APPROOT)

	SetKeys()
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
		err := ioutil.WriteFile(File, jsonHinfo, 0755)
		if err != nil {
			log.Println(err)
		} else {
			fabs, _ := filepath.Abs(File)
			Echo("OK. Save to File", fabs)
		}
	}
}

func Echo(title string, t interface{}) {
	if Quiet != true {
		log.Println("=======", title, "=======")
		log.Println(t)
	}
}

func LoadMiscDir(pth string) error {
	miscNotice := `INFO: user can put "gohostinfo.user.***.json" file in "misc" folder 
	to customize Key = Value. You have to keep the json format string in *ONLY* one line,
	format(josn) should be like : 
	{"key1":[{"k11":"val11"},{"k12":"val12"}],"key2":[{"k21":"val21"}]}
	they will be added into final result:
	  result["data"]["key1"]=[{"k11":"val11"},{"k12":"val12"}]
	  result["data"]["key2"]=[{"k21":"val21"}]`

	Echo("LoadMiscDir", miscNotice)

	err := filepath.Walk(pth, func(pth string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if filepath.Ext(pth) != ".json" {
			return nil
		}

		if strings.Index(filepath.Base(pth), "gohostinfo.user.") == -1 {
			return nil
		}

		Echo("Checking misc folder", pth)
		fpth, err := os.Open(pth)
		if err != nil {
			log.Println(pth, err)
			return err
		}
		defer fpth.Close()

		cnt, err := ioutil.ReadAll(fpth)
		if err != nil {
			log.Println(err)
		}

		var data map[string][]map[string]interface{}

		err = json.Unmarshal(cnt, &data)
		if err != nil {
			log.Println(err)
			return err
		}

		for key, val := range data {
			Echo("Load Key", key)
			Hinfo.Data[key] = val
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return err
}

func SetKeys() {
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

	Echo("Get Hostname", Hinfo.ID)
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
	diskStat := []*disk.UsageStat{}
	almostFullDisks := []*disk.UsageStat{}
	partitions, _ := disk.Partitions(true)
	for _, p := range partitions {
		if p.Mountpoint != "" {
			pu, err := disk.Usage(p.Mountpoint)
			if err != nil {
				log.Println(err)
			} else {
				diskStat = append(diskStat, pu)
				if pu.UsedPercent > 80 {
					almostFullDisks = append(almostFullDisks, pu)
				}
			}

		}
	}
	Hinfo.Data["disk"] = diskStat
	Hinfo.Data["disk_full"] = almostFullDisks

	Echo("Disk", diskStat)
	Echo("Disk Almost Full", almostFullDisks)
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
