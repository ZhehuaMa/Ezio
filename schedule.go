package Ezio

import "time"

var bufferLength int64
var routineNumber int64
var workers []worker

func Run() {
	initEnv()
	go schedule()
}

func Stop() {
	stopSignal <- struct{}{}
	<-stopSignal
}

func schedule() {
	defer close(stopSignal)
	sendInterval := time.Second / time.Duration(qps)
	tickerChan := time.NewTicker(sendInterval).C

	for {
		<-tickerChan
		select {
		case <-stopSignal:
			joinRoutines()
			stopMonitor()
			return
		case dispatcher <- <-queue:
		default:
			statisticsChan <- &result{resultType: busy}
		}
	}
}

func Append(task *Task) {
	queue <- task
}

func initEnv() {
	bufferLength = qps * 5
	routineNumber = bufferLength
	queue = make(chan *Task, bufferLength)
	stopSignal = make(chan struct{})
	statisticsChan = make(chan *result, qps)
	stopMonitorSignal = make(chan struct{})
	workers = make([]worker, routineNumber)
	for i := range workers {
		workers[i].boss = dispatcher
		workers[i].offDuty = make(chan struct{})
	}

	startMonitor()
	startRoutines()
}

func startRoutines() {
	for i := range workers {
		go workers[i].work()
	}
}

func joinRoutines() {
	for i := range workers {
		workers[i].rest()
	}
	for i := range workers {
		workers[i].ackRest()
	}
}
