package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func deploy(pc []string) {
	rc := readConf()
	if rc.App.Repo == "" {
		rc.App.Repo = localCommand("git config --get remote.origin.url")
		rc.App.Repo = rc.App.Repo[:len(rc.App.Repo)-5]
	}
	// l_com := "git add . && git commit -m \"updates\" && git push -u origin master"
	c_com := "pkill -f " + rc.App.DomainName + " || true && cd " + rc.GCloud.LiveDir + " && rm -rf " +
		rc.App.DomainName + " || true && curl -o tmp.zip -L " + rc.App.Repo +
		"/zipball/master && unzip tmp.zip && mv *" + strings.Split(rc.App.Repo, "/")[3] + "-* " +
		rc.App.DomainName + " && rm tmp.zip && cd " + rc.App.DomainName + " && go build -o " +
		rc.App.DomainName + " && ./" + rc.App.DomainName + " &; disown || true && pkill bp " +
		"|| true && bp &; disown"
		// localCommand(l_com)
	gitPush()
	cloudCommand([]string{c_com})
}

// func gcloudSCP() {
// 	cs := `gcloud compute scp` +
// 		` --zone ` + rc.GCloud.Zone +
// 		` --project ` + rc.GCloud.Project +
// 		` --recurse . ` +
// 		rc.GCloud.Instance + `:` +
// 		rc.GCloud.LiveDir + rc.App.DomainName
// 	fmt.Println(cs)
// 	localCommand(cs)
// }

func restartProxy() {
	cloudCommand([]string{"pkill", "bp", "||", "true", "&&", "bp", "&;", "disown"})
}

func localCommand(com string) string {
	fmt.Println("Running: ", com)
	var command []string
	// if !strings.Contains(com, "--") {
	// 	if strings.Contains(com, "-") {
	// 		c := strings.Split(com, "-")
	// 		c1 := c[0]
	// 		c2 := c[1]
	// 		command = append(command, c1)
	// 		fmt.Println(command)
	// 		command = append(command, c2)
	// 		fmt.Println(strings.Join(command, " "))
	// 	}
	// } else {
	command = strings.Split(com, " ")
	// }

	var cmd *exec.Cmd = &exec.Cmd{}
	cmd = exec.Command(command[0], command[1:]...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("local command error: ", err, string(o))
	}
	// fmt.Println(o)
	return string(o)
}

func cloudCommand(command []string) string {
	args := []string{`compute`, `ssh`, `--zone`, `us-central1-a`, `main`, `--project`, `mysterygift`, `--`}
	tmx := "tmux send-keys -t dashboard:main '" + strings.Join(command, " ") + "' Enter"

	args = append(args, strings.Split(tmx, " ")...)
	cmd := exec.Command(`gcloud`, args...)
	fmt.Println(cmd.String())
	o, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(o)
}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
