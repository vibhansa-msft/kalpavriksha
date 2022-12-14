package main

import (
	"fmt"
	"sync"
)

type workItem struct {
	workerId int
	path     string
	objtype  ObjectType
	status   JobStatusType
}

func startWorkers() {
	config.wgWorkers = sync.WaitGroup{}

	config.jobs = make(chan workItem, config.Parallelism*2)
	config.results = make(chan workItem, config.Parallelism*2)

	for w := 1; w <= config.Parallelism; w++ {
		config.wgWorkers.Add(1)
		go uploadWorker(w)
	}

	go createJobs()

	pendingCount := config.NumberOfDirs * config.NumberOfFiles
	completecount := int64(0)

	for job := range config.results {
		fmt.Printf("Worker %d => %s : %s (%s)\n", job.workerId, job.objtype, job.path, job.status)
		completecount++

		if completecount == pendingCount {
			close(config.results)
		}
	}

	config.wgWorkers.Wait()
}

func createJobs() {
	for d := (int64)(0); d < config.NumberOfDirs; d++ {
		for f := (int64)(0); f < config.NumberOfFiles; f++ {
			config.jobs <- workItem{
				path:    fmt.Sprintf("%d/%d", d, f),
				objtype: EObjectType.FILE(),
				status:  EJobStatusType.WAIT(),
			}
		}
	}
	close(config.jobs)
}

func uploadWorker(w int) {
	defer config.wgWorkers.Done()
	for job := range config.jobs {
		//fmt.Printf("(%d) %s\n", w, job.path)
		job.workerId = w
		job.status = EJobStatusType.SUCCESS()
		config.results <- job
	}
}
