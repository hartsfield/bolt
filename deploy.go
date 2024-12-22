package main

import (
	"fmt"
	"os"
	"strings"
)

func deploy(pc []string) {
	os.Setenv("GOARCH", "amd64")
	localCommand(strings.Split("go build -o "+rc.App.DomainName, " "))
	gcloudSCP()
}

func gcloudSCP() {
	cs := `gcloud compute scp` +
		` --zone ` + rc.GCloud.Zone +
		` --project ` + rc.GCloud.Project +
		` --recurse ` + rc.App.DomainName + ` internal/bolt.conf.json ` +
		rc.GCloud.Instance + `:` +
		rc.GCloud.LiveDir + rc.App.DomainName
	fmt.Println(cs)
	localCommand(strings.Split(cs, " "))
}

func restartProxy(rc *config) {
	cloudCommand([]string{"./" + rc.GCloud.ProxyConf + " &"})
}
