package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	fsNetworkHeight *prometheus.GaugeVec
}

type config struct {
	url        string
	mcNetwork  string
	checkWait  time.Duration
	timeout    time.Duration
	listenHost string
	listenPort int
}

type FsRequest struct {
	Method  string      `json:"method"`
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Params  interface{} `json:"params,omitempty"`
}

type FsReturn struct {
	Method  string  `json:"method"`
	JsonRPC string  `json:"jsonrpc"`
	ID      int     `json:"id"`
	Err     FsError `json:"error"`
	// Result  map[string]interface{} `json:"result"`
}

type FsError struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// get_network_status structs
// {
// 		network_status: {
// 			network_block_height: "11111111"
// 		}
// }

type FsGetNetworkStatusReturn struct {
	FsReturn
	Result FsGetNetworkStatusResult `json:"result"`
}

type FsGetNetworkStatusResult struct {
	NetworkStatus FsGetNetworkStatusNetworkStatus `json:"network_status"`
}

type FsGetNetworkStatusNetworkStatus struct {
	NetworkBlockHeight string `json:"network_block_height"`
}
