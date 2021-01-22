package Ezio

import (
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	Request *http.Request
}

type worker struct {
	boss    <-chan *Task
	offDuty chan struct{}
	atHome  chan struct{}
}

type result struct {
	resultType int
	latency    time.Duration
}

type httpStatusCodeError struct {
	statusCode int
}

func (e httpStatusCodeError) Error() string {
	return fmt.Sprintf("http status code is %d", e.statusCode)
}
