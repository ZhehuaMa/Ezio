package Ezio

import (
	"github.com/golang/glog"
	"time"
)

var format =
	"Last %f secs: QPS %f, success %d, failure %d, busy %d, AvgLatency %f ms, MinLatency %f ms, MaxLatency %f"

func startMonitor() {
	go monitor()
}

func stopMonitor() {
	stopMonitorSignal <- struct{}{}
	<-stopMonitorSignal
}

func monitor() {
	defer close(stopMonitorSignal)

	ticker := time.NewTicker(10 * time.Second).C

	globalStartTime := time.Now()
	startTime := globalStartTime

	printStatistics := func() {
		endTime := time.Now()
		interval := endTime.Sub(startTime)
		globalInterval := endTime.Sub(globalStartTime)

		currentQps := float64(currentSuccess) / interval.Seconds()
		globalQps := float64(totalSuccess) / globalInterval.Seconds()

		avgLatency := .0
		if currentSuccess != 0 {
			avgLatency = float64(currentLatency) / float64(currentSuccess)
		}
		globalAvgLatency := .0
		if totalSuccess != 0{
			globalAvgLatency = float64(totalLatency) / float64(totalSuccess)
		}

		glog.Infof(
			format,
			interval.Seconds(),
			currentQps,
			currentSuccess,
			currentFailure,
			currentBusy,
			avgLatency,
			currentMinLatency.Seconds() * 1000,
			currentMaxLatency.Seconds() * 1000,
			)
		glog.Infof(
			format,
			globalInterval.Seconds(),
			globalQps,
			totalSuccess,
			totalFailure,
			totalBusy,
			globalAvgLatency,
			totalMinLatency.Seconds() * 1000,
			totalMaxLatency.Seconds() * 1000,
		)

		resetCurrentLatencies()
		startTime = time.Now()
	}

	for {
		select {
		case res := <-statisticsChan:
			switch res.resultType {
			case success:
				totalSuccess++
				currentSuccess++
				totalLatency += res.latency
				currentLatency += res.latency
				updateLatencies(res.latency)
			case failure:
				totalFailure++
				currentFailure++
			case busy:
				totalBusy++
				currentBusy++
			}
		case <-ticker:
			printStatistics()
		case <-stopMonitorSignal:
			printStatistics()
			return
		}
	}
}

func updateLatencies(latency time.Duration) {
	if totalMaxLatency < latency {
		totalMaxLatency = latency
	}
	if currentMaxLatency < latency {
		currentMaxLatency = latency
	}
	if totalMinLatency > latency {
		totalMinLatency = latency
	}
	if currentMinLatency > latency {
		currentMinLatency = latency
	}
}

func resetCurrentLatencies() {
	currentLatency = 0
	currentMinLatency = 0
	currentMaxLatency = 0
}
