package Ezio

import (
	"io"
	"io/ioutil"
	"net/url"
	"time"
)
import "github.com/golang/glog"

func (w *worker) work() {
	for {
		select {
		case <-w.offDuty:
			w.atHome <- struct{}{}
			return
		case task := <-w.boss:
			startTime := time.Now()
			err := doTask(task)
			diff := time.Now().Sub(startTime)
			if err != nil {
				statisticsChan <- &result{resultType: failure}
			} else {
				statisticsChan <- &result{
					resultType: success,
					latency:    diff,
				}
			}
		}
	}
}

func (w *worker) rest() {
	go func() {
		w.offDuty <- struct{}{}
	}()
}

func (w *worker) ackRest() {
	<-w.atHome
}

func doTask(task *Task) error {
	response, netError := client.Do(task.Request)
	if netError != nil {
		switch err := netError.(type) {
		case *url.Error:
			if err.Timeout() {
				glog.Errorf("Timeout: %s", err.Error())
			} else {
				glog.Errorf("Connection error: %s", err.Error())
			}
		default:
			glog.Errorf("Network error: %s", err.Error())
		}
		return netError
	}
	defer response.Body.Close()
	io.Copy(ioutil.Discard, response.Body)

	if response.StatusCode/100 != 2 {
		return httpStatusCodeError{statusCode: response.StatusCode}
	}
	return nil
}
