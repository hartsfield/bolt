package main

import (
	"log"
	"strings"
)

func deploy(pc []string) {
	localCommand([]string{"git", "init"})
	localCommand([]string{"git", "add", "."})
	localCommand([]string{"git", "commit", "-m", `"deployment"`})
	localCommand([]string{"git", "push", "-u", "origin", "master"})
	checkInsert(pc[0], pc[1], pc[2], pc[3])
}
func checkInsert(name, port, hasTLS, alertsOn string) {
	op := cloudCommand([]string{"cat ~/gp/prox.conf"})
	sop := strings.Split(op, "\n")
	for _, o := range sop {
		if strings.Contains(o, ":") {
			_port := strings.Split(o, ":")[0]
			_name := strings.Split(o, ":")[3]
			if port == _port || name == _name {
				log.Fatalln("Error: Port or name already in use.")
			}
		}
	}
	echo := strings.Join([]string{port, hasTLS, alertsOn, name}, ":")
	cloudCommand([]string{`echo "` + echo + `" >> prox.conf`})
	restartProxy()
}

func restartProxy() {
	cloudCommand([]string{"./go_proxy/restart-service &"})
}
