package workerpool

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Scheduler struct {
	size          int
	maxSize       int
	lastID        int
	Jobs          chan *Job
	IdlingWorkers chan *Worker
	Workers       map[int]*Worker
	wg            sync.WaitGroup
}

func NewScheduler(poolSize int) *Scheduler {
	return &Scheduler{
		maxSize:       poolSize,
		Jobs:          make(chan *Job, 10+poolSize),
		IdlingWorkers: make(chan *Worker, poolSize),
		Workers:       make(map[int]*Worker),
	}
}

func (scheduler *Scheduler) AddJob(job *Job) {
	scheduler.Jobs <- job
}

func (scheduler *Scheduler) GetFreeWorker() *Worker {
	select {
	case worker := <-scheduler.IdlingWorkers:
		return worker
	default:
		var worker *Worker
		if scheduler.size < scheduler.maxSize {
			scheduler.size++
			scheduler.lastID++
			worker = NewWorker(scheduler.lastID)

			scheduler.Workers[worker.id] = worker
			scheduler.wg.Add(1)
			go worker.Work(scheduler)
		} else {
			worker = <-scheduler.IdlingWorkers
		}
		return worker
	}
}

func (scheduler *Scheduler) Start() {
	for job := range scheduler.Jobs {
		worker := scheduler.GetFreeWorker()
		worker.job <- job
	}
	for scheduler.size > 0 {
		worker := <-scheduler.IdlingWorkers
		worker.stop <- 1
	}
}

func (scheduler *Scheduler) PrintWorkers() {
	fmt.Printf("  ID -> Status\n")
	for _, worker := range scheduler.Workers {
		switch worker.status {
		case StatusIdling:
			fmt.Printf("*%3d -> IDLING\n", worker.id)
		case StatusWorking:
			fmt.Printf("*%3d -> WORKING\n", worker.id)
		}
	}
}

func Read(scheduler *Scheduler) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		words := strings.Fields(text)
		if len(words) == 0 {
			continue
		}

		command := words[0]
		switch command {
		case "exit":
			return
		case "print":
			scheduler.PrintWorkers()
		default:
			job, err := NewJob(command)
			if err != nil {
				fmt.Println("Invalid input:", err)
				continue
			}
			scheduler.AddJob(job)
		}
	}
	close(scheduler.Jobs)
}

func Run(poolSize int) {
	scheduler := NewScheduler(poolSize)

	go scheduler.Start()
	Read(scheduler)

	scheduler.wg.Wait()
}
