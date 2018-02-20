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

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// --- Usages ---
// TODO: Network

func (sys *System) CPUUsage() {
	usages, err := cpu.Percent(1, true)
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
	sys.MemAv = vmem.Available
	sys.MemUsed = vmem.Used
}

// --- Conf ---

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

	return sys, nil
}

// --- Update ---

func (sys *System) Update() {
	sys.CPUUsage()
	sys.DiskUsage()
	sys.MemUsage()
}

// --- Misc ---

func (sys *System) String() string {

	str := fmt.Sprintf("Kernel: %s\n", sys.KernelVersion)

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
	str = fmt.Sprintf("%sMem:\n\tAvailable: %v\n\tUsed: %v\n", str,
		sys.MemAv, sys.MemUsed)

	return str
}

type System struct {
	// Info
	BootTime      uint64
	CPUInfo       []cpu.InfoStat
	CPUThreads    int
	Partitions    []string `json:"partitions"`
	KernelVersion string
	Uptime        uint64

	// Usages
	CPUUsages  []float64
	DiskUsages [][2]uint64
	MemAv      uint64
	MemUsed    uint64
}
