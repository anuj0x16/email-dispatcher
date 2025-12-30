package dispatcher

import "github.com/anuj0x16/email-dispatcher/internal/jobs"

type worker struct {
	id          int
	job         chan jobs.EmailJob
	workerQueue chan *worker
	quitChan    chan bool
}

func newWorker(id int, workerQueue chan *worker) *worker {
	return &worker{
		id:          id,
		job:         make(chan jobs.EmailJob),
		workerQueue: workerQueue,
		quitChan:    make(chan bool),
	}
}

func (w *worker) start() {
	go func() {
		for {
			w.workerQueue <- w

			select {
			case _ = <-w.job:
				// TODO: send an email

			case <-w.quitChan:
				// TODO: worker has been requested to stop
			}
		}
	}()
}

func (w *worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}
