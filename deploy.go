package main

import (
	"fmt"
	"os"
	"strings"
)

func deploy(pc []string) {
	os.Setenv("GOARCH", "amd64")
	localCommand(strings.Split("go build -o "+rc.App.DomainName, " "))
	// localCommand([]string{"git", "init"})
	// localCommand([]string{"git", "add", "."})
	// localCommand([]string{"git", "commit", "-m", `"deployment"`})
	// localCommand([]string{"git", "push", "-u", "origin", "master"})
	gcloudSCP()
	// checkInsert()
}

//	func checkInsert() {
//		op := cloudCommand([]string{"cat /home/" + rc.GCloud.User + rc.GCloud.ProxyConf})
//		sop := strings.Split(op, "\n")
//		for _, o := range sop {
//			if strings.Contains(o, ":") {
//				_port := strings.Split(o, ":")[0]
//				_name := strings.Split(o, ":")[3]
//				if rc.App.Port == _port || rc.App.Name == _name {
//					log.Fatalln("Error: Port or name already in use.")
//				}
//			}
//		}
//		echo := strings.Join([]string{rc.App.Port, rc.GCloud.User, rc.App.AlertsOn, rc.App.Name}, ":")
//		cloudCommand([]string{`echo "` + echo + `" >> /home/` + rc.GCloud.User + rc.GCloud.ProxyConf})
//		restartProxy(rc)
//	}
func gcloudSCP() {
	cs := `gcloud compute scp` +
		` --zone ` + rc.GCloud.Zone +
		` --project ` + rc.GCloud.Project +
		` --recurse ` + rc.App.DomainName + ` internal/ bolt.conf.json ` +
		rc.GCloud.Instance + `:` +
		rc.GCloud.LiveDir + rc.App.DomainName
	fmt.Println(cs)
	localCommand(strings.Split(cs, " "))
}
func restartProxy(rc *config) {
	cloudCommand([]string{"./" + rc.GCloud.ProxyConf + " &"})
}
