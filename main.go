package main

import (
	"fmt"
	"log"

	sysmon "github.com/mikelsr/sysmon/utils"
)

func main() {
	sys, err := sysmon.LoadSystem()
	if err != nil {
		log.Fatal(err)
	}
	sys.Update()
	fmt.Println(sys)
}
