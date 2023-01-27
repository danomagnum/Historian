package main

import (
	"context"
	"log"
	"time"
)

func NewHistorianLogging() (*HistorianLogging, error) {
	h := new(HistorianLogging)
	h.c = make(chan []HistorianData, 1024)
	h.DataCache = map[string][]historianLoggingDataEntry{}
	h.Timeout = time.Minute

	return h, nil
}

// this only stores float64s!!!
type HistorianLogging struct {
	DataCache map[string][]historianLoggingDataEntry
	c         chan []HistorianData
	Timeout   time.Duration
}

func (h *HistorianLogging) Close() {
}

func (h *HistorianLogging) C() chan<- []HistorianData {
	return h.c
}

func (h *HistorianLogging) Run(ctx context.Context) {

	t := time.NewTicker(h.Timeout)
	for {
		select {
		case hd := <-h.c:
			// new data came in so grab it and put it in the format we need for processing
			for i := range hd {
				dc, ok := h.DataCache[hd[i].Name]
				if !ok {
					dc = make([]historianLoggingDataEntry, 0)
				}
				dc = append(dc, historianLoggingDataEntry{Timestamp: hd[i].Timestamp, Value: hd[i].Value})
				h.DataCache[hd[i].Name] = dc
			}

		case <-t.C:
			// time to write the data out
			log.Printf("Data Cache: %+v", h.DataCache)
			h.DataCache = map[string][]historianLoggingDataEntry{}

		case <-ctx.Done():
		}
	}
}

type historianLoggingDataEntry struct {
	Timestamp time.Time
	Value     any
}
