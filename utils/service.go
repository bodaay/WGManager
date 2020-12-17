package utils

import (
	"time"

	"github.com/kardianos/service"
)

var Logger service.Logger

var ServiceIntervalInSeconds uint

var RunProcess func() error

var ServiceStopNow bool

var ReadyToStop bool

type Program struct {
	exit chan struct{}
}

func (p *Program) Start(s service.Service) error {

	if service.Interactive() {
		Logger.Info("Running in terminal.")
	} else {
		Logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *Program) run() error {
	go RunProcess() //this function will be set by main
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			// Logger.Infof("ZenExporter Still running at %v...", tm)
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
}
func (p *Program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	Logger.Info("I'm Stopping!")
	ServiceStopNow = true
	//we will wait for 5 seconds only, if not, we will just ignore that shit and force stop
	MaxWaitTime := 5 * 1000
	totalWaitTime := 0
	for {
		time.Sleep(10 * time.Millisecond)
		if ReadyToStop || totalWaitTime >= MaxWaitTime {
			break
		}
		totalWaitTime += 10
	}
	Logger.Info("Service Stopped")
	close(p.exit)
	return nil
}
