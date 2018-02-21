package sysmon

/*

IDEA: Funciones que carguen cada uso en el struct
      Función de actualización que llama a todas las anteriores
      Función para reportar estado

*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// --- Usages ---
// TODO: Network

func (sys *System) CPUUsage() {
	usages, err := cpu.Percent(time.Second, true)
	if err != nil {
		log.Fatal(err)
	}
	sys.CPUUsages = usages
}

func (sys *System) DiskUsage() {
	usages := make([][2]uint64, len(sys.Partitions))
	for i, p := range sys.Partitions {
		usage, err := disk.Usage(p)
		if err != nil {
			log.Fatal(err)
		}
		usages[i][0] = usage.Total
		usages[i][1] = usage.Used
	}
	sys.DiskUsages = usages
}

func (sys *System) MemUsage() {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}
	sys.MemTotal = vmem.Total
	sys.MemUsed = vmem.Used
}

// --- Conf ---

func LoadConf() (*Conf, error) {
	conf := new(Conf)

	// Load configuration from json file
	confFile, err := ioutil.ReadFile("conf.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(confFile, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func LoadSystem() (*System, error) {

	sys := new(System)

	// Load configuration from json file
	confFile, err := ioutil.ReadFile("system.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(confFile, sys)
	if err != nil {
		return nil, err
	}

	// Load system information

	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	sys.CPUInfo = info
	sys.CPUThreads = len(info)

	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	sys.KernelVersion = hostInfo.KernelVersion

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	sys.HostName = hostName

	return sys, nil
}

// --- Update ---

func (sys *System) Measure() {
	sys.CPUUsage()
	sys.DiskUsage()
	sys.MemUsage()
}

// --- Misc ---

func (sys *System) String() string {
	str := fmt.Sprintf("Hostname: %s\n", sys.HostName)
	str = fmt.Sprintf("%sKernel: %s\n", str, sys.KernelVersion)

	// cpu
	str = fmt.Sprintf("%sCPU:\n\tThreads:%v\n\tUsages:%v\n", str,
		sys.CPUThreads, sys.CPUUsages)

	// disk
	str = fmt.Sprintf("%sDisk:\n", str)
	for i, p := range sys.Partitions {
		str = fmt.Sprintf("%s\t%s:\n\t\tTotal: %v\n\t\tFree: %v\n", str, p,
			sys.DiskUsages[i][0], sys.DiskUsages[i][1])
	}

	// mem
	str = fmt.Sprintf("%sMem:\n\tTotal: %v\n\tUsed: %v\n", str,
		sys.MemTotal, sys.MemUsed)

	return str
}

func (conf *Conf) BaseURL() string {
	var protocol string
	if conf.InfluxDB.TLS {
		protocol = "https"
	} else {
		protocol = "http"
	}

	return fmt.Sprintf("%s://%s:%d", protocol, conf.InfluxDB.Host, conf.InfluxDB.Port)
}

func (conf *Conf) URI() string {
	return fmt.Sprintf("%s/write?db=%s",
		conf.BaseURL(), conf.InfluxDB.DB)
}

type Conf struct {
	InfluxDB struct {
		DB   string `json:"db"`
		Host string `json:"host"`
		Port int    `json:"port"`
		TLS  bool   `json:"tls"`
	} `json: "influxdb"`
	Interval int `json:"interval"`
}

type System struct {
	// Info
	BootTime      uint64
	CPUInfo       []cpu.InfoStat
	CPUThreads    int
	HostName      string
	Partitions    []string `json:"partitions"`
	KernelVersion string
	Uptime        uint64

	// Usages
	CPUUsages  []float64
	DiskUsages [][2]uint64
	MemUsed    uint64
	MemTotal   uint64
}
