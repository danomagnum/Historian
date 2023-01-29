package main

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type ConfigHistorianInflux struct {
	Name   string
	Server string
	Token  string
	Org    string
	Bucket string
}

func NewHistorianInflux(server, token, org, bucket string) (*HistorianInflux, error) {
	h := new(HistorianInflux)
	h.c = make(chan []HistorianData, 1024)
	h.Client = influxdb2.NewClient(server, token)
	h.WriteAPI = h.Client.WriteAPI(org, bucket)
	h.Org = org
	h.Bucket = bucket
	h.Server = server
	h.Token = token

	return h, nil
}

// this only stores float64s!!!
type HistorianInflux struct {
	Server   string
	Token    string
	Org      string
	Bucket   string
	WriteAPI api.WriteAPI
	c        chan []HistorianData
	Client   influxdb2.Client
}

func (h *HistorianInflux) Close() {
}

func (h *HistorianInflux) C() chan<- []HistorianData {
	return h.c
}

func (h *HistorianInflux) Run(ctx context.Context) {

	for {
		select {
		case hd := <-h.c:
			// new data came in so grab it and put it in the format we need for processing
			for i := range hd {
				v := map[string]any{"Value": hd[i].Value}
				p := influxdb2.NewPoint(hd[i].Name, nil, v, hd[i].Timestamp)
				h.WriteAPI.WritePoint(p)
			}

		case <-ctx.Done():
			return
		}
	}
}