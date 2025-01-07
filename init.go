package main

import (
	"flag"
	"io"
	"log"
	"os"
)

type config struct {
	App    app    `json:"app"`
	GCloud gcloud `json:"gcloud"`
}

type app struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	DomainName string `json:"domain_name"`
	Version    string `json:"version"`
	Env        env    `json:"env"`
	Port       string `json:"port"`
	AlertsOn   bool   `json:"alertsOn"`
	TLSEnabled bool   `json:"tls_enabled"`
	Repo       string `json:"repo"`
}

type env map[string]string

type gcloud struct {
	Command   string `json:"command"`
	Zone      string `json:"zone"`
	Project   string `json:"project"`
	User      string `json:"user"`
	LiveDir   string `json:"livedir"`
	Instance  string `json:"instance"`
	ProxyHome string `json:"proxy_home"`
	ProxyConf string `json:"proxyConf"`
}

type stringFlag struct {
	set   bool
	value string
	Info  string
	Name  string
	do    func([]string)
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

var (
	flagCommands []*stringFlag = []*stringFlag{
		{Name: "init", do: boltInit, Info: "\nEx.: --init appName\n\nInitializes a new bolt project in the directory 'appName'\n"},
		{Name: "new-page", do: createPage, Info: "\nEx.: --new-page pageName\n\nInitializes a new page with the given name\n"},
		{Name: "new-component", do: createComponent, Info: "\nEx.: --new-component componentName\n\nInitializes a new component with the given name\n"},
		{Name: "new-route", do: newRoute, Info: "\nEx.: --new-route routeName,handler\n\nInitializes a new route\n"},
		{Name: "insert-component", do: insertcomponent, Info: "\nEx.: --insert-component componentName,pageName\n\nInserts a component into a page\n"},
		{Name: "streamable", do: genStream, Info: "\nEx.: --streamable ./model.json \n\nCreates a stream of uploadable items with an upload form based on a json model\n"},
		{Name: "deploy", do: deploy, Info: "\nEx.: --deploy\n\nDeploys project to server using values from bolt.conf.json\n"},
		{Name: "autonav", do: autonav, Info: "\nEx.: --autonav page1,page2,page3,pageEtc\n\nInitializes a new navbar component with the given pages\n"},
		{Name: "autosplash", do: autosplash, Info: "\nEx.: --autosplash public/filename.png\n\nInitializes a splash screen component\n"},
		{Name: "remote-service-restart", do: serviceReload, Info: "\nEx.: --remote-service-restart\n\nRestarts a remote service using values from bolt.conf.json\n"},
		{Name: "config", do: configure, Info: "\nEx.: --config\n\nInteractive configuration of bolt\n"},
		// {Name: "add-style", do: , "Adds a style to the stylesheet of the given component, usage: bolt --add-style=component:rulename")},
		// {Name: "build-form", do: , "Genrates an HTML form based on input")},
		// {Name: "install-component", do: , "Installs a component from a git hub repo")
	}
	rc *config
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for _, com := range flagCommands {
		flag.Var(com, com.Name, com.Info)
	}
}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
