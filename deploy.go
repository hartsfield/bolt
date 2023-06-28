package main

import (
	"log"
	"os/exec"
	"strings"
)

func deploy(proxconfig string) {
	pc := strings.Split(proxconfig, " ")
	localCommand([]string{"git", "init"})
	localCommand([]string{"git", "add", "."})
	localCommand([]string{"git", "commit", "-m", `"deployment"`})
	localCommand([]string{"git", "push", "-u", "origin", "master"})
	checkInsert(pc[0], pc[1], pc[2], pc[3])
}

func localCommand(command []string) string {
	cmd := exec.Command(command[0], command[1:]...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(o)
}

func cloudCommand(command []string) string {
	args := []string{`compute`, `ssh`, `--zone`, `us-central1-a`, `instance-2`, `--project`, `mysterygift`, `--`}
	args = append(args, command...)
	cmd := exec.Command(`gcloud`, args...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(o)
}

func checkInsert(name, port, hasTLS, alertsOn string) {
	op := cloudCommand([]string{"cat prox.conf"})
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
