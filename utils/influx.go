package sysmon

import (
	"fmt"
	// "log"
	"net/http"
	"net/url"
	"strings"
)

func CreateDB(conf *Conf) error {
	// log.Printf("Creating database %s\n", conf.InfluxDB.DB)
	form := url.Values{}
	form.Add("q", fmt.Sprintf("CREATE DATABASE %s", conf.InfluxDB.DB))

	_, err := http.PostForm(fmt.Sprintf("%s/query", conf.BaseURL()), form)
	return err
}

func BuildPoint(sys *System, measurement string, body string) string {
	return fmt.Sprintf("%s,hostname=%s,kernel_version=%s%s,uptime=%d\n",
		measurement, sys.HostName, sys.KernelVersion, body, sys.Uptime)
}

func PostStatus(conf *Conf, sys *System) error {
	// log.Printf("Sendig data to %s\n", conf.BaseURL())
	_, err := http.Post(conf.URI(), "text/plain", strings.NewReader(RequestBody(sys)))
	return err
}

func RequestBody(sys *System) string {
	body := ""
	// cpu
	for i := 0; i < sys.CPUThreads; i++ {
		body += BuildPoint(sys, "cpu_load", fmt.Sprintf(",core=%d value=%.2f",
			i+1, sys.CPUUsages[i]))
	}
	// mem
	body += BuildPoint(sys, "mem_used", fmt.Sprintf(",mem_total=%d value=%d",
		sys.MemTotal, sys.MemUsed))

	// disk
	for i, p := range sys.Partitions {
		body += BuildPoint(sys, "disk_used", fmt.Sprintf(",partition=%s,total=%d value=%d",
			p, sys.DiskUsages[i][0], sys.DiskUsages[i][1]))
	}
	return body
}
