package sysmon

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func CreateDB(conf *Conf) error {
	form := url.Values{}
	form.Add("q", fmt.Sprintf("CREATE DATABASE %s", conf.InfluxDB.DB))

	req, err := http.PostForm(fmt.Sprintf("%s/query", conf.BaseURL()), form)
	return err
}

func PointHeader(sys *System, measurement string) string {
	return fmt.Sprintf("%s,hostname=%s,kernel_version=%s,uptime=%d",
		measurement, sys.HostName, sys.KernelVersion, sys.Uptime)
}

func PostStatus(conf *Conf, sys *System) error {
	_, err := http.Post(conf.URI(), "text/plain", strings.NewReader(RequestBody(sys)))
	return err
}

func RequestBody(sys *System) string {
	body := ""
	// cpu
	for i := 0; i < sys.CPUThreads; i++ {
		body += fmt.Sprintf("%s,core=%d value=%.2f\n",
			PointHeader(sys, "cpu_load"), i+1, sys.CPUUsages[i])
	}
	// mem
	body += fmt.Sprintf("%s,mem_total=%d value=%d\n",
		PointHeader(sys, "mem_used"), sys.MemTotal, sys.MemUsed)
	// disk
	for i, p := range sys.Partitions {
		body += fmt.Sprintf("%s,partition=%s,total=%d value=%d\n",
			PointHeader(sys, "disk_used"), p,
			sys.DiskUsages[i][0], sys.DiskUsages[i][1])
	}
	return body
}
