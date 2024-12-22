package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var questions []string = []string{
	"App Name: ", "Version: ", "Command: ", "Domain: ", "Port: ", "Repo: ",
	"Login: ", "Zone: ", "Instance: ", "Project: ", "Live Directory: ",
}

var answers []string = []string{
	"Bolt App", "0.01", "boltapp", "domain_name.com", "9123", "", "", "",
	"", "", "~/live",
}

func configure(answerString []string) {
	newconf := &config{}
	var tlsBoolValue bool
	var err error
	if len(answerString) == len(questions)+1 {
		answers = answerString
		tlsBoolValue, err = strconv.ParseBool(answerString[len(answerString)])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		for i, q := range questions {
			fmt.Println(q)
			a, err := reader.ReadString('\n')
			if err != nil {
				log.Println(err)
			}
			answers[i] = strings.TrimSpace(a)
		}
		fmt.Print("TLS Enabled (default: true): ")
		tlsenabled, _ := reader.ReadString('\n')
		tlsBoolValue, err = strconv.ParseBool(tlsenabled[:len(tlsenabled)-1])
		if err != nil {
			log.Fatal(err)
		}
	}
	newconf.App.Name = answers[0]
	newconf.App.Command = answers[1]
	newconf.App.DomainName = answers[2]
	newconf.App.Port = answers[3]
	newconf.App.Version = answers[4]
	newconf.App.Repo = answers[5]
	newconf.GCloud.Zone = answers[6]
	newconf.GCloud.Instance = answers[7]
	newconf.GCloud.Project = answers[8]
	newconf.GCloud.LiveDir = answers[9]
	newconf.GCloud.User = answers[10]
	newconf.App.TLSEnabled = tlsBoolValue

	b, err := json.Marshal(&newconf)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(b))
}
