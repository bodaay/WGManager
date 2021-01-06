package main

import (
	"WGManager/utils"
	"WGManager/webapi"
	"WGManager/wg"
	"log"
	"os"
	"path"
	"path/filepath"
)

//WGManagerVersion WGManager Software Version
const WGManagerVersion = "1.0"

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
	//Creating service file for the project if it doesnt exist
	servicefilename := "wgmanager.service"
	//fmt.Printf("Arg length %d\n", len(os.Args))
	//Get current executable name and path
	//TODO: this logic sucks, fix it
	appExec, err := os.Executable()
	if err != nil {
		panic(err)
	}
	serviceFileLocation := path.Join("/etc/systemd/system", servicefilename)
	log.Println(serviceFileLocation)
	if !utils.FileExists(serviceFileLocation) {
		err := utils.InstallLinuxService(servicefilename, "WgManager Core Service", filepath.Dir(appExec), appExec, "root")
		if err != nil {
			log.Println(err) //don't want to panic if this is the case
		}
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
	go webapi.StartClient(&wgc, WGManagerVersion)
	webapi.StartAdminClient(&wgc)
}
