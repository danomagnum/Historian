package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/danomagnum/gologix"
)

type ConfigCIPClass3 struct {
	Name        string
	Address     string // IP Address
	Path        string // CIP Path (ex: "1,0" for slot 0 on the backplane)
	Enable      bool
	DefaultRate time.Duration
	Endpoints   []EndpointCIPClass3
}

func (config *ConfigCIPClass3) Init(ctx context.Context, h map[string]Historian) {
	go config.Run(ctx, h)
}

func (config *ConfigCIPClass3) Run(ctx context.Context, h map[string]Historian) {
	var err error
	client := gologix.NewClient(config.Address)
	client.Path, err = gologix.ParsePath(config.Path)
	if err != nil {
		log.Printf("problem starting CipClass3 Client %s: %v", config.Name, err)
		return
	}

	// split the endpoints up by poll rate
	poll_groups := make(map[time.Duration][]EndpointCIPClass3)
	for i := range config.Endpoints {
		rate := config.Endpoints[i].Rate
		group, ok := poll_groups[rate]
		if !ok {
			poll_groups[rate] = make([]EndpointCIPClass3, 0)
			group = poll_groups[rate]
		}
		group = append(group, config.Endpoints[i])
		poll_groups[rate] = group
	}

	// Connect
	if config.Enable {
		err = client.Connect()
		if err != nil {
			log.Printf("could not connect to %s: %v", config.Name, err)
			return
		}
	}

	subctx, cancel_func := context.WithCancel(ctx)
	// Wait for poll groups
	for k := range poll_groups {
		log.Printf("CIPClass3 %s Rate %v has %d endpoints", config.Name, k, len(poll_groups[k]))
		go config.PollGroup(subctx, client, k, poll_groups[k], h)
	}
	<-ctx.Done()
	cancel_func()
	err = client.Disconnect()
	if err != nil {
		log.Printf("problem disconnecting from %s: %v", config.Name, err)
	}

}

func (config *ConfigCIPClass3) PollGroup(ctx context.Context, client *gologix.Client, rate time.Duration, endpoints []EndpointCIPClass3, h map[string]Historian) {

	tags := make([]string, len(endpoints))
	types := make([]gologix.CIPType, len(endpoints))

	for i := range endpoints {
		tags[i] = endpoints[i].TagName
		types[i] = endpoints[i].TagType
	}

	t := time.NewTicker(rate)
	for {
		select {
		case <-t.C:
			hd := make([]HistorianData, len(endpoints))
			ts := time.Now()
			values, err := client.ReadList(tags, types)
			if err != nil {
				log.Printf("problem reading %s at %v: %v", config.Name, rate, err)
				continue
			}
			for i := range values {
				endpoints[i].Value = values[i]
				hd[i] = HistorianData{
					Timestamp: ts,
					Name:      fmt.Sprintf("%s.%s", config.Name, endpoints[i].TagName),
					Value:     values[i],
				}
				if endpoints[i].Historian != "" {
					h[endpoints[i].Historian].C() <- hd
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

type EndpointCIPClass3 struct {
	Name      string
	TagName   string
	Rate      time.Duration
	TagType   gologix.CIPType
	Value     any
	Historian string
}

func (e EndpointCIPClass3) TypeAsInt() int {
	return int(e.TagType)
}
