package dispatcher

import "github.com/anuj0x16/email-dispatcher/internal/jobs"

type Dispatcher struct {
	JobQueue    chan jobs.EmailJob
	workerQueue chan *worker
	nworkers    int
}

func New(jobQueueSize, nworkers int) *Dispatcher {
	return &Dispatcher{
		JobQueue:    make(chan jobs.EmailJob, jobQueueSize),
		workerQueue: make(chan *worker, nworkers),
		nworkers:    nworkers,
	}
}

func (d *Dispatcher) Start() {
	for i := range d.nworkers {
		worker := newWorker(i+1, d.workerQueue)
		worker.start()
	}

	go func() {
		for job := range d.JobQueue {
			go func() {
				worker := <-d.workerQueue
				worker.job <- job
			}()
		}
	}()
}
