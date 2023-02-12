package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type SystemState int

const (
	StateStartup SystemState = iota
	StateRunning
	StateRestarting
)

func (s SystemState) String() string {
	switch s {
	case StateStartup:
		return "Starting Up"
	case StateRunning:
		return "Running"
	case StateRestarting:
		return "Restarting"
	}
	return "Unknown"
}

type System struct {
	ActiveContext       context.Context
	ActiveContextCancel context.CancelFunc
	ActiveConfig        Config
	WorkingConfig       Config
	Changes             bool
	State               SystemState
}

var system System

func main() {
	system.State = StateStartup

	var err error
	////////////////////////
	// Load Config
	////////////////////////
	system.ActiveConfig, err = ConfigLoad("cfg.json")
	if err != nil {
		system.ActiveConfig = ConfigNew()
	}
	system.WorkingConfig = system.ActiveConfig

	server_addr := fmt.Sprintf("%s:%d", system.ActiveConfig.General.Host, system.ActiveConfig.General.Port)
	log.Printf("Starting server at %s", server_addr)

	WebAPIStart()

	////////////////////////
	// Main Loop
	////////////////////////
	for {
		// this context will live until the config changes.
		// eventually the cancel function here will be called by the web interface when
		// the config changes to stop everything and we'll start over at that point.
		system.ActiveContext, system.ActiveContextCancel = context.WithCancel(context.Background())

		////////////////////////
		// Init Historians
		////////////////////////
		Historians := make(map[string]Historian)

		for i := range system.ActiveConfig.Historians.Influx {
			system.ActiveConfig.Historians.Influx[i].Init(system.ActiveContext, Historians)
		}
		for i := range system.ActiveConfig.Historians.JSON {
			system.ActiveConfig.Historians.JSON[i].Init(system.ActiveContext, Historians)
		}
		for i := range system.ActiveConfig.Historians.Logging {
			system.ActiveConfig.Historians.Logging[i].Init(system.ActiveContext, Historians)
		}

		////////////////////////
		// Init Data Providers
		////////////////////////
		for i := range system.ActiveConfig.DataProviders.CIPClass3 {
			system.ActiveConfig.DataProviders.CIPClass3[i].Init(system.ActiveContext, Historians)
		}

		system.State = StateRunning

		////////////////////////
		// Wait for config change
		////////////////////////
		<-system.ActiveContext.Done()
		system.State = StateRestarting

		log.Printf("Active Context Complete. Restart Delay: %v.", system.ActiveConfig.General.RestartDelay)

		// wait a bit before restarting.
		time.Sleep(system.ActiveConfig.General.RestartDelay)
		system.ActiveConfig = system.WorkingConfig
		system.Changes = false
		log.Printf("Restarting...")
	}

}
