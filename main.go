package main

import (
	"WGManager/utils"
	"WGManager/webapi"
	"WGManager/wg"
	"log"
	"os"
)

func main() {
	defaultConfigFilePath := "wgmanconfig.json"
	if len(os.Args) > 1 {
		defaultConfigFilePath = os.Args[1]
	}
	runningAsRoot, err := utils.CheckIfAdminOrRoot()
	if err != nil {
		panic(err)
	}
	if !runningAsRoot {
		log.Fatalln("You must run this app as Admin or Root!")
	}
	ls := utils.ExecTask{
		Command: "ls",
		Args:    []string{"-l"},
		Shell:   false,
	}
	_, err = ls.Execute()
	if err != nil {
		panic(err)
	}
	//Load the config file
	var wgc wg.WGConfig
	err = wgc.ParseConfigFile(defaultConfigFilePath)
	if err != nil {
		newconfig, err := wgc.CreateDefaultconfig(defaultConfigFilePath)
		if err != nil {
			panic(err)
		}
		wgc = *newconfig
	}
	//Search the path for instances configuration files
	err = wgc.LoadInstancesFiles()
	if err != nil {
		panic(err)
	}
	go webapi.StartClient(&wgc)
	webapi.StartAdminClient(&wgc)
}
