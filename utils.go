package Ezio

import (
	"math"
	"net/http"
	"time"
)

const (
	success = iota
	failure
	busy
	slow
)

var (
	queue               chan *Task
	dispatcher          chan *Task
	stopSignal          chan struct{}
	stopResponse        chan struct{}
	statisticsChan      chan *result
	stopMonitorSignal   chan struct{}
	stopMonitorResponse chan struct{}
)

var tr = &http.Transport{
	MaxIdleConns:        20,
	MaxIdleConnsPerHost: 2,
	MaxConnsPerHost:     10,
	IdleConnTimeout:     10 * time.Second,
}

var client = http.Client{
	Transport: tr,
	Timeout:   time.Minute,
}
var (
	totalSuccess int64 = 0
	totalFailure int64 = 0
	totalBusy    int64 = 0
	totalSlow    int64 = 0

	currentSuccess int64 = 0
	currentFailure int64 = 0
	currentBusy    int64 = 0
	currentSlow    int64 = 0

	totalLatency    = time.Duration(0)
	totalMaxLatency = time.Duration(math.MinInt64)
	totalMinLatency = time.Duration(math.MaxInt64)

	currentLatency    = time.Duration(0)
	currentMaxLatency = time.Duration(math.MinInt64)
	currentMinLatency = time.Duration(math.MaxInt64)
)
