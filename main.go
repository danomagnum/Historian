package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/danomagnum/gologix"
)

func main() {

	ep := []EndpointCIPClass3{
		{Name: "ShedTemp", TagName: "Program:Shed.Temp1", Rate: time.Second, TagType: gologix.CIPTypeREAL, Value: 0, Historian: "Weather"},
		{Name: "OutsideTemp", TagName: "Program:Shed.Temp2", Rate: time.Second, TagType: gologix.CIPTypeREAL, Value: 0, Historian: "Weather"},
		{Name: "ShedRH", TagName: "Program:Shed.RH1", Rate: time.Second, TagType: gologix.CIPTypeREAL, Value: 0, Historian: "Weather"},
		{Name: "OutsideRH", TagName: "Program:Shed.RH2", Rate: time.Second, TagType: gologix.CIPTypeREAL, Value: 0, Historian: "Weather"},
		{Name: "LivingRoomTemp", TagName: "Program:RpiTempHum1.Temperature", Rate: time.Second, TagType: gologix.CIPTypeREAL, Value: 0, Historian: "Weather"},
		{Name: "LivingRommRH", TagName: "Program:RpiTempHum1.Humidity", Rate: time.Second, TagType: gologix.CIPTypeREAL, Value: 0, Historian: "Weather"},
		{Name: "GaragePrs", TagName: "Program:Garage.Pressure_inHG", Rate: time.Second, TagType: gologix.CIPTypeREAL, Value: 0, Historian: "Weather"},
		{Name: "GarageTemp", TagName: "Program:Garage.Temp", Rate: time.Second, TagType: gologix.CIPTypeREAL, Value: 0, Historian: "Weather"},
	}

	hc := ConfigHistorianInflux{
		Name:   "Weather",
		Server: "http://historian.home:8086",                                                               // server
		Token:  "76GIWIAF7BF7zcQQFDRLsd0t2uplZheI1_6yHq3T8-8y01PUwynCdG11qVzUjo8OiplWdllFUS2D35sjiC8JYA==", // token
		Org:    "home",                                                                                     // organization
		Bucket: "weather",                                                                                  // bucket
	}

	c3 := ConfigCIPClass3{}
	c3.Address = "192.168.2.241"
	c3.Path = "1,0"
	c3.DefaultRate = time.Second * 1
	c3.Enable = true
	c3.Name = "GaragePLC"
	c3.Endpoints = ep

	c := ConfigNew()
	c.General.Host = "localhost"
	c.General.Port = 8000
	c.DataProviders.CIPClass3 = append(c.DataProviders.CIPClass3, c3)
	c.Historians.Influx = append(c.Historians.Influx, hc)

	err := c.Save("cfg.json")
	if err != nil {
		log.Panicf("problem creating config: %v", err)
	}

	conf, err := ConfigLoad("cfg.json")
	if err != nil {
		conf = ConfigNew()
	}

	server_addr := fmt.Sprintf("%s:%d", conf.General.Host, conf.General.Port)
	log.Printf("Starting server at %s", server_addr)

	// this context will live until the config changes.
	// eventually the cancel function here will be called by the web interface when
	// the config changes to stop everything and we'll start over at that point.
	ctx, _ := context.WithCancel(context.Background())

	Historians := make(map[string]Historian)
	//h, err := NewHistorianJSON("data/test.json")
	//h, err := NewHistorianLogging()
	for i := range conf.Historians.Influx {
		conf.Historians.Influx[i].Init(ctx, Historians)
	}

	// Start all the drivers
	CipClass3(ctx, Historians, conf.DataProviders.CIPClass3)

	select {}

}
