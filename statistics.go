package Ezio

import (
	"github.com/golang/glog"
	"math"
	"time"
)

var format = "Last %f secs: QPS %f, success %d, failure %d, busy %d, slow %d, AvgLatency %f ms, MinLatency %f ms, MaxLatency %f ms"

func startMonitor() {
	go monitor()
}

func stopMonitor() {
	stopMonitorSignal <- struct{}{}
	<-stopMonitorResponse
}

func monitor() {
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
			avgLatency = currentLatency.Seconds() * 1000.0 / float64(currentSuccess)
		}
		globalAvgLatency := .0
		if totalSuccess != 0 {
			globalAvgLatency = totalLatency.Seconds() * 1000.0 / float64(totalSuccess)
		}

		glog.Infoln("-------------------------------------------------------------------------------" +
			"---------------------------------------------------------------------------------")
		glog.Infof(
			format,
			interval.Seconds(),
			currentQps,
			currentSuccess,
			currentFailure,
			currentBusy,
			currentSlow,
			avgLatency,
			currentMinLatency.Seconds()*1000,
			currentMaxLatency.Seconds()*1000,
		)
		glog.Infof(
			format,
			globalInterval.Seconds(),
			globalQps,
			totalSuccess,
			totalFailure,
			totalBusy,
			totalSlow,
			globalAvgLatency,
			totalMinLatency.Seconds()*1000,
			totalMaxLatency.Seconds()*1000,
		)
		glog.Infoln("-------------------------------------------------------------------------------" +
			"---------------------------------------------------------------------------------")

		resetCurrentStatistics()
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
			case slow:
				totalSlow++
				currentSlow++
			}
		case <-ticker:
			printStatistics()
		case <-stopMonitorSignal:
			printStatistics()
			stopMonitorResponse <- struct{}{}
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

func resetCurrentStatistics() {
	currentLatency = 0
	currentMinLatency = math.MaxInt64
	currentMaxLatency = math.MinInt64
	currentSuccess = 0
	currentFailure = 0
	currentBusy = 0
	currentSlow = 0
}
