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
		{TagName: "TestInt", Rate: time.Second, TagType: gologix.CIPTypeINT, Value: 0},
		{TagName: "TestDint", Rate: time.Second, TagType: gologix.CIPTypeDINT, Value: 0},
		{TagName: "TestBool", Rate: time.Second, TagType: gologix.CIPTypeBOOL, Value: 0},
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
	c.CIPClass3 = append(c.CIPClass3, c3)

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
	ctx, cancel_func := context.WithCancel(context.Background())

	//h, err := NewHistorianJSON("data/test.json")
	h, err := NewHistorianLogging()
	//h, err := NewHistorianTStore("data/")
	if err != nil {
		log.Panicf("failed to open historian %v", err)
	}

	go h.Run(ctx)

	// Start all the drivers
	CipClass3(ctx, h, conf.CIPClass3)

	select {}

	cancel_func()

}
