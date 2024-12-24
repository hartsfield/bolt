package main

import (
	"fmt"
	"strings"
)

func deploy(pc []string) {
	if pc[0] == "init" {
		gcloudSCP()
		return
	}
	cloudCommand([]string{
		"pkill", rc.App.Command, "||", "true", "&&",
		"cd", rc.GCloud.LiveDir + rc.App.DomainName, "&&",
		"git", "pull", "||", "true", "&&",
		"go", "build", "-o", rc.App.Command, "&&",
		"mv", rc.App.Command, "~/bin", "&&",
		rc.App.Command, "&; disown", "||", "true", "&&", "pkill", "bp",
		"||", "true", "&&", "bp", "&;", "disown",
	})
}

func gcloudSCP() {
	cs := `gcloud compute scp` +
		` --zone ` + rc.GCloud.Zone +
		` --project ` + rc.GCloud.Project +
		` --recurse . ` +
		rc.GCloud.Instance + `:` +
		rc.GCloud.LiveDir + rc.App.DomainName
	fmt.Println(cs)
	localCommand(strings.Split(cs, " "))
}

func restartProxy() {
	cloudCommand([]string{"pkill", "bp", "||", "true", "&&", "bp", "&;", "disown"})
}
