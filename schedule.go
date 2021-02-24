package Ezio

import (
	"flag"
	"time"
)

var bufferLength int64
var routineNumber int64
var workers []worker

func Run() {
	initEnv()
	go schedule()
}

func Stop() {
	stopSignal <- struct{}{}
	<-stopResponse
}

func WaitUntilFinish() {
	close(queue)
	<-stopResponse
}

func schedule() {
	sendInterval := time.Second / time.Duration(qps)
	tickerChan := time.NewTicker(sendInterval).C

	for {
		<-tickerChan
		select {
		case <-stopSignal:
			joinRoutines()
			stopMonitor()
			stopResponse <- struct{}{}
			return
		case task, ok := <-queue:
			if !ok {
				joinRoutines()
				stopMonitor()
				stopResponse <- struct{}{}
				return
			}
			sendTask(task)
		default:
			statisticsChan <- &result{resultType: hungry}
		}
	}
}

func sendTask(task *Task) {
	select {
	case dispatcher <- task:
	default:
		statisticsChan <- &result{resultType: busy}
	}
}

func Append(task *Task) {
	queue <- task
}

func initEnv() {
	if !flag.Parsed() {
		flag.Parse()
	}

	bufferLength = qps * 5
	routineNumber = bufferLength * 10
	queue = make(chan *Task, bufferLength)
	dispatcher = make(chan *Task)
	stopSignal = make(chan struct{})
	stopResponse = make(chan struct{})
	statisticsChan = make(chan *result, qps)
	stopMonitorSignal = make(chan struct{})
	stopMonitorResponse = make(chan struct{})
	workers = make([]worker, routineNumber)
	for i := range workers {
		workers[i].boss = dispatcher
		workers[i].offDuty = make(chan struct{})
		workers[i].atHome = make(chan struct{})
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
