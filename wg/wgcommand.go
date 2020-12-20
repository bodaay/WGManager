package wg

import (
	"WGManager/utils"
	"fmt"
)

var controlWgInstance = utils.ExecTask{
	Command: "/usr/bin/systemctl",
	Args:    []string{"Action", ""},
	Shell:   true,
}

func restartWGInstance(instanceName string) (string, error) {
	cmd := controlWgInstance
	cmd.Args[1] = fmt.Sprintf("wg-quick@%s", instanceName)
	cmd.Args[0] = "restart"
	res, err := cmd.Execute()
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}

func stopWGInstance(instanceName string) (string, error) {
	cmd := controlWgInstance
	cmd.Args[1] = fmt.Sprintf("wg-quick@%s", instanceName)
	cmd.Args[0] = "stop"
	res, err := cmd.Execute()
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}

func startWGInstance(instanceName string) (string, error) {
	cmd := controlWgInstance
	cmd.Args[1] = fmt.Sprintf("wg-quick@%s", instanceName)
	cmd.Args[0] = "start"
	res, err := cmd.Execute()
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}
