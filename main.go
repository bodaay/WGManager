package main

import (
	"WGManager/utils"
	"WGManager/wg"
	"os"
)

func main() {
	var wgconfig wg.WGConfig
	defaultConfigPath := "wgmanconfig.json"

	UseNAT := false
	NATAdapterName := "eno1"
	MaxInstances := 0
	if len(os.Args) > 1 {
		defaultConfigPath = os.Args[1]
	}
	if !utils.FileExists(defaultConfigPath) {
		err := wgconfig.ParseConfig(wg.DefaultJSONFileString)
		if err != nil {
			panic(err)
		}
		err = wgconfig.InitiateConfig(UseNAT, NATAdapterName, uint16(MaxInstances))
		if err != nil {
			panic(err)
		}
		err = wgconfig.SaveConfigFile(defaultConfigPath)
		if err != nil {
			panic(err)
		}
	} else {
		err := wgconfig.ParseConfigFile(defaultConfigPath)
		if err != nil {
			panic(err)
		}
	}

}
