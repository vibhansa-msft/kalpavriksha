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
	kalpavriksha.wgWorkers = sync.WaitGroup{}

	kalpavriksha.jobs = make(chan workItem, config.Parallelism*2)
	kalpavriksha.results = make(chan workItem, config.Parallelism*2)

	for w := 1; w <= config.Parallelism; w++ {
		kalpavriksha.wgWorkers.Add(1)
		go uploadWorker(w)
	}

	go createJobs()

	pendingCount := config.NumberOfDirs * config.NumberOfFiles
	completecount := int64(0)

	for job := range kalpavriksha.results {
		completecount++

		fmt.Printf("Worker %d => %s : %s (%s), job Completion %0.2f\n",
			job.workerId, job.objtype, job.path, job.status,
			float64(completecount)*100/float64(pendingCount))

		if completecount == pendingCount {
			close(kalpavriksha.results)
		}
	}

	kalpavriksha.wgWorkers.Wait()
}

func createJobs() {
	for d := (int64)(0); d < config.NumberOfDirs; d++ {
		for f := (int64)(0); f < config.NumberOfFiles; f++ {
			kalpavriksha.jobs <- workItem{
				path:    fmt.Sprintf("dir-%d/file-%d", d, f),
				objtype: EObjectType.FILE(),
				status:  EJobStatusType.WAIT(),
			}
		}
	}
	close(kalpavriksha.jobs)
}

func uploadWorker(w int) {
	defer kalpavriksha.wgWorkers.Done()
	for job := range kalpavriksha.jobs {
		//fmt.Printf("(%d) %s\n", w, job.path)
		job.workerId = w

		job.status = EJobStatusType.INPROGRESS()

		data, err := kalpavriksha.dataSrc.GetData(uint64(config.FileSize))
		if err != nil {
			job.status = EJobStatusType.FAILED()
		} else {

			opt := getUploadOptions(data)
			err := kalpavriksha.storage.UploadData(job.path, data, opt)
			if err != nil {
				job.status = EJobStatusType.FAILED()
			} else {
				job.status = EJobStatusType.SUCCESS()
			}
		}

		kalpavriksha.results <- job
	}
}

func getUploadOptions(data []byte) *UploadOptions {
	if config.UpdateMD5 == true || config.Tier != "none" {
		opt := &UploadOptions{}

		if config.UpdateMD5 {
			opt.MD5Sum = kalpavriksha.dataSrc.GetMd5Sum(data)
		}

		if config.Tier != "none" {
			opt.Tier = &config.BlobTier
		}

		return opt
	}

	return nil
}
