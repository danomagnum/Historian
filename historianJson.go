package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func NewHistorianJSON(filename string) (*HistorianJSON, error) {
	h := new(HistorianJSON)
	var err error

	h.File, err = os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("problem opening historian json file %s: %w", filename, err)
	}

	h.JsonEncoder = *json.NewEncoder(h.File)
	h.c = make(chan []HistorianData, 1024)

	return h, nil
}

type HistorianJSON struct {
	File        io.WriteCloser
	JsonEncoder json.Encoder
	c           chan []HistorianData
}

func (h *HistorianJSON) Close() {
	h.File.Close()
}

func (h *HistorianJSON) C() chan<- []HistorianData {
	return h.c
}

func (h *HistorianJSON) Run(ctx context.Context) {

	for {
		select {
		case hd := <-h.c:
			for i := range hd {
				h.JsonEncoder.Encode(hd[i])
			}
		case <-ctx.Done():

		}
	}
}
