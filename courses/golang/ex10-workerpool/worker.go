package workerpool

import (
	"fmt"
	"strconv"
	"time"
)

const (
	StatusIdling = iota
	StatusWorking
)

type Job struct {
	duration float64
}

func NewJob(duration string) (*Job, error) {
	f, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		return nil, err
	}
	return &Job{duration: f}, nil
}

type Worker struct {
	id     int
	status int
	job    chan *Job
	stop   chan byte
}

func NewWorker(id int) *Worker {
	w := &Worker{
		id:     id,
		status: StatusIdling,
		job:    make(chan *Job, 1),
		stop:   make(chan byte, 1),
	}

	fmt.Printf("worker:%d spawning\n", w.id)

	return w
}

func (w *Worker) Work(s *Scheduler) {
	for {
		select {
		case <-w.stop:
			fmt.Printf("worker:%d stopping\n", w.id)
			s.size -= 1
			defer s.wg.Done()
			return
		case job := <-w.job:
			fmt.Printf("worker:%d sleep:%.1f\n", w.id, job.duration)
			w.status = StatusWorking
			time.Sleep(time.Millisecond * time.Duration(int(job.duration*1000)))
			w.status = StatusIdling
			s.IdlingWorkers <- w
		}
	}
}
