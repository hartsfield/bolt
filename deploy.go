package main

func deploy(pc []string) {
	localCommand([]string{"git", "init"})
	localCommand([]string{"git", "add", "."})
	localCommand([]string{"git", "commit", "-m", `"deployment"`})
	localCommand([]string{"git", "push", "-u", "origin", "master"})
	// checkInsert()
}

// func checkInsert() {
// 	op := cloudCommand([]string{"cat /home/" + rc.GCloud.User + rc.GCloud.ProxyConf})
// 	sop := strings.Split(op, "\n")
// 	for _, o := range sop {
// 		if strings.Contains(o, ":") {
// 			_port := strings.Split(o, ":")[0]
// 			_name := strings.Split(o, ":")[3]
// 			if rc.App.Port == _port || rc.App.Name == _name {
// 				log.Fatalln("Error: Port or name already in use.")
// 			}
// 		}
// 	}
// 	echo := strings.Join([]string{rc.App.Port, rc.GCloud.User, rc.App.AlertsOn, rc.App.Name}, ":")
// 	cloudCommand([]string{`echo "` + echo + `" >> /home/` + rc.GCloud.User + rc.GCloud.ProxyConf})
// 	restartProxy(rc)
// }

func restartProxy(rc *config) {
	cloudCommand([]string{"./" + rc.GCloud.ProxyConf + " &"})
}
