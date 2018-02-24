package main

import (
	"log"
	"time"

	sysmon "github.com/mikelsr/sysmon/utils"
)

func main() {
	sys, err := sysmon.LoadSystem()
	if err != nil {
		log.Fatal(err)
	}

	conf, err := sysmon.LoadConf()
	if err != nil {
		log.Fatal(err)
	}
	for {
		err = sysmon.CreateDB(conf)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(3))
	}

	for {
		sys.Measure()
		// fmt.Println(sysmon.RequestBody(sys))
		_ = sysmon.PostStatus(conf, sys)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		time.Sleep(time.Second * time.Duration(conf.Interval))
	}
}
