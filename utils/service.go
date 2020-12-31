package utils

import (
	"fmt"
	"io/ioutil"
	"path"
)

//InstallLinuxService the service file for the current project
func InstallLinuxService(ServiceName string, ServiceDescription string, ServiceWorkingDirectory string, executableFile string, User string) error {
	var controlService = ExecTask{
		Command: "/usr/bin/systemctl",
		Args:    []string{"Enable", ""},
		Shell:   true,
	}
	//we will only support ubuntu at the moment, later Imran will supprt all os
	osInstallServiceFilesLocation := "/etc/systemd/system"
	serviceFileLocation := path.Join(osInstallServiceFilesLocation, ServiceName)
	// executableLocation := path.Join(webAPILocation, excutableFile)

	serviceFileData := `[Unit]
Description=` + ServiceDescription + `

[Service]
WorkingDirectory=` + ServiceWorkingDirectory + `
ExecStart=` + executableFile + `
RestartSec=always
# Auto restart in 3 Seconds
Restart=3
# Time to wait for process to exit before sending kill signal
TimeoutStopSec=10
KillSignal=SIGINT
SyslogIdentifier=` + ServiceName + `-Log
# This is how I created this user: sudo adduser --no-create-home --disabled-login --shell /bin/false dotnet
User=` + User + `
# add as many Env vairbales required here
#Environment=GO_ENVIRONMENT=Production

[Install]
WantedBy=multi-user.target
`

	err := ioutil.WriteFile(serviceFileLocation, []byte(serviceFileData), 0755)
	if err != nil {
		fmt.Printf("Service File Created Failed\n")
		return err
	}
	cmd := controlService
	cmd.Args[1] = ServiceName
	cmd.Args[0] = "enable"
	_, err = cmd.Execute()
	if err != nil {
		return err
	}
	return nil

}
