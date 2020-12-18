package main

import (
	"WGManager/wg"
	"log"
	"os"
)

func main() {
	defaultConfigFilePath := "wgmanconfig.json"
	if len(os.Args) > 1 {
		defaultConfigFilePath = os.Args[1]
	}
	//Load the config file
	var wgc wg.WGConfig
	err := wgc.ParseConfigFile(defaultConfigFilePath)
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
	wgi, err := wgc.FindInstanceByIPAndPort("172.27.40.0/22", 22201)
	if err != nil {
		panic(err)
	}
	if wgi == nil {
		err = wgc.CreateNewInstance("172.27.40.0/22", 22201, []string{"1.1.1.1", "8.8.8.8"}, true, "eno1", 0)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Instance already exist for the ip")
	}
	err = wgc.DeployAllInstances()
	if err != nil {
		panic(err)
	}
	// for _, i := range wgc.WGInstances {
	// 	log.Println(i.WGClients[0])
	// }
}
