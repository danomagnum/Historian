package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

var activeConf Config
var workingConf Config
var changes bool

var activeContext context.Context

func main() {

	var err error
	////////////////////////
	// Load Config
	////////////////////////
	activeConf, err = ConfigLoad("cfg.json")
	if err != nil {
		activeConf = ConfigNew()
	}
	workingConf = activeConf

	server_addr := fmt.Sprintf("%s:%d", activeConf.General.Host, activeConf.General.Port)
	log.Printf("Starting server at %s", server_addr)

	WebAPIStart()

	////////////////////////
	// Main Loop
	////////////////////////
	for {
		// this context will live until the config changes.
		// eventually the cancel function here will be called by the web interface when
		// the config changes to stop everything and we'll start over at that point.
		activeContext, _ := context.WithCancel(context.Background())

		////////////////////////
		// Init Historians
		////////////////////////
		Historians := make(map[string]Historian)

		for i := range activeConf.Historians.Influx {
			activeConf.Historians.Influx[i].Init(activeContext, Historians)
		}
		for i := range activeConf.Historians.JSON {
			activeConf.Historians.JSON[i].Init(activeContext, Historians)
		}
		for i := range activeConf.Historians.Logging {
			activeConf.Historians.Logging[i].Init(activeContext, Historians)
		}

		////////////////////////
		// Init Data Providers
		////////////////////////
		for i := range activeConf.DataProviders.CIPClass3 {
			activeConf.DataProviders.CIPClass3[i].Init(activeContext, Historians)
		}

		////////////////////////
		// Wait for config change
		////////////////////////
		<-activeContext.Done()

		log.Printf("Active Context Complete. Restart Delay: %v.", activeConf.General.RestartDelay)

		// wait a bit before restarting.
		time.Sleep(activeConf.General.RestartDelay)
		log.Printf("Restarting...")
	}

}
