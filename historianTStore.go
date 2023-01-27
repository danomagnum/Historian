package main

import (
	"log"

	"github.com/nakabonne/tstorage"
)

func NewHistorianTStore(filename string) (*HistorianTStore, error) {
	h := new(HistorianTStore)
	s, err := tstorage.NewStorage(tstorage.WithDataPath(filename))
	if err != nil {
		return nil, err
	}
	h.TS = s
	h.c = make(chan []HistorianData, 1024)

	return h, nil
}

// this only stores float64s!!!
type HistorianTStore struct {
	TS tstorage.Storage
	c  chan []HistorianData
}

func (h *HistorianTStore) Close() {
	h.TS.Close()
}
func (h *HistorianTStore) Add() chan<- []HistorianData {
	return h.c
}

func (h *HistorianTStore) Run(hd []HistorianData) {
	for i := range hd {
		val64, ok := h.ToFloat64(hd[i].Value)
		if ok {
			err := h.TS.InsertRows([]tstorage.Row{
				{
					Metric:    hd[i].Name,
					DataPoint: tstorage.DataPoint{Timestamp: hd[i].Timestamp.Unix(), Value: val64},
				},
			})

			if err != nil {
				log.Printf("Problem writing to storage: %v", err)
			}
		}
	}
}

func (h *HistorianTStore) ToFloat64(val any) (float64, bool) {
	switch x := val.(type) {
	case byte:
		return float64(x), true
	case int8:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case float32:
		return float64(x), true
	case float64:
		return x, true
	}
	return 0, false
}
