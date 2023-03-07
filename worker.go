package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
)

type workItem struct {
	workerId int
	path     string
	objtype  ObjectType
	status   JobStatusType
}

var WaitCount int64 = 0
var totalProcessedCount int64 = 0

func startWorkers() {
	kalpavriksha.wgWorkers = sync.WaitGroup{}

	if config.CreateStub || config.DeleteStub {
		kalpavriksha.jobs = make(chan workItem, 10000000)
		kalpavriksha.results = make(chan workItem, 10000000)
	} else {
		kalpavriksha.jobs = make(chan workItem, config.Parallelism*2)
		kalpavriksha.results = make(chan workItem, config.Parallelism*2)
	}

	atomic.AddInt64(&WaitCount, int64(config.Parallelism))

	for w := 1; w <= config.Parallelism; w++ {
		kalpavriksha.wgWorkers.Add(1)
		if config.CreateStub || config.DeleteStub {
			go createStubWorker(w)
		} else if config.Delete {
			go deleteWorker(w)
		} else if config.SetTier {
			go tierWorker(w)
		} else {
			go uploadWorker(w)
		}
	}

	if config.CreateStub || config.DeleteStub {
		// Push the root directory to the queue
		kalpavriksha.jobs <- workItem{
			path:    "",
			objtype: EObjectType.DIR(),
			status:  EJobStatusType.WAIT(),
		}

		go func() {
			t := time.Tick(time.Duration(60 * time.Second))
			log.Printf("Starting monitor")

			for {
				select {
				case <-t:
					log.Printf("Completed item count: %v", atomic.LoadInt64(&totalProcessedCount))
				}
			}
		}()

		completecount := 0
		tickerCount := 0
		ticker := time.Tick(time.Duration(20 * time.Second))

		for {
			select {
			case _ = <-kalpavriksha.results:
				completecount++
				tickerCount = 0
			case <-ticker:
				tickerCount++
				if atomic.LoadInt64(&WaitCount) == int64(config.Parallelism) {
					if tickerCount > 3 {
						close(kalpavriksha.jobs)
						kalpavriksha.wgWorkers.Wait()
						close(kalpavriksha.results)
						log.Printf("Number of stubs created %d\n", completecount)
						fmt.Printf("Number of stubs created %d\n", completecount)
						return
					}
				} else {
					tickerCount = 0
				}
			}
		}

	} else {
		go createJobs()

		pendingCount := config.NumberOfDirs * config.NumberOfFiles
		completecount := int64(0)

		for job := range kalpavriksha.results {
			completecount++

			log.Printf("Worker %d => %s : %s (%s), job Completion %0.2f\n",
				job.workerId, job.objtype, job.path, job.status,
				float64(completecount)*100/float64(pendingCount))

			if completecount == pendingCount {
				close(kalpavriksha.results)
			}
		}

		kalpavriksha.wgWorkers.Wait()
	}
}

func createJobs() {
	depth := ""
	for i := int64(0); i < config.DirDepth; i++ {
		depth += fmt.Sprintf("%d/", i+1)
	}

	for d := (int64)(0); d < config.NumberOfDirs; d++ {
		for f := (int64)(0); f < config.NumberOfFiles; f++ {
			name := fmt.Sprintf("dir-%d/%sfile-%d", d, depth, f)
			kalpavriksha.jobs <- workItem{
				path:    name,
				objtype: EObjectType.FILE(),
				status:  EJobStatusType.WAIT(),
			}
		}
	}
	close(kalpavriksha.jobs)
}

// Workers for upload task
func uploadWorker(w int) {
	defer kalpavriksha.wgWorkers.Done()
	for job := range kalpavriksha.jobs {
		log.Printf("(%d) %s\n", w, job.path)
		job.workerId = w

		job.status = EJobStatusType.INPROGRESS()

		data, err := kalpavriksha.dataSrc.GetData()
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

// Workers for delete task
func deleteWorker(w int) {
	defer kalpavriksha.wgWorkers.Done()
	for job := range kalpavriksha.jobs {
		log.Printf("(%d) %s\n", w, job.path)
		job.workerId = w

		job.status = EJobStatusType.INPROGRESS()

		err := kalpavriksha.storage.Delete(job.path, nil)
		if err != nil {
			job.status = EJobStatusType.FAILED()
		} else {
			job.status = EJobStatusType.SUCCESS()
		}

		kalpavriksha.results <- job
	}
}

// worker to change tier of given data set
func tierWorker(w int) {
	defer kalpavriksha.wgWorkers.Done()
	for job := range kalpavriksha.jobs {
		log.Printf("(%d) %s\n", w, job.path)
		job.workerId = w

		job.status = EJobStatusType.INPROGRESS()

		err := kalpavriksha.storage.SetTier(job.path, config.BlobTier)
		if err != nil {
			job.status = EJobStatusType.FAILED()
		} else {
			job.status = EJobStatusType.SUCCESS()
		}

		kalpavriksha.results <- job
	}
}

// Workers for delete task
func createStubWorker(w int) {
	defer kalpavriksha.wgWorkers.Done()
	for job := range kalpavriksha.jobs {
		atomic.AddInt64(&WaitCount, -1)

		job.workerId = w
		job.status = EJobStatusType.INPROGRESS()

		// List the items
		pager := kalpavriksha.storage.ListBlobs(job.path)

		listCnt := uint64(0)
		// Iterate blob prefixes
		for pager.More() {
			resp, err := pager.NextPage(context.TODO())
			if err == nil {
				listCnt += uint64(len(resp.Segment.BlobItems))
				if listCnt > 100000 {
					atomic.AddInt64(&totalProcessedCount, int64(listCnt))
					listCnt = 0
				}

				//if resp.Marker != nil && resp.NextMarker != nil {
				//	//log.Printf("(%d) Path : %s, Current Marker : %s, Next Marker : %s\n",
				//	//	job.workerId, job.path, *resp.Marker, *resp.NextMarker)
				//	log.Printf("(%d) Path : %s, Current Count: %d\n", job.workerId, job.path, listCnt)
				//}

				for _, item := range resp.Segment.BlobPrefixes {
					dirPath := *item.Name
					if dirPath[len(dirPath)-1] == '/' {
						dirPath = dirPath[:len(dirPath)-1]
					}

					// Get properties of directory
					if config.CreateStub {
						err = kalpavriksha.storage.CreateStub(dirPath)
						if err == nil {
							log.Printf("(%d) Stub creatd for %s", job.workerId, dirPath)
						} else if bloberror.HasCode(err, bloberror.BlobAlreadyExists) {
							log.Printf("(%d) Stub already exists for %s\n", job.workerId, dirPath)
						} else {
							log.Printf("(%d) Failed to create stub unknown error : %s\n", job.workerId, err.Error())
						}
					} else if config.DeleteStub {
						err = kalpavriksha.storage.Delete(dirPath, nil)
						if err == nil {
							log.Printf("(%d) Stub deleted for %s", job.workerId, dirPath)
						}
					}

					w := workItem{
						path:     dirPath + "/",
						workerId: job.workerId,
						objtype:  EObjectType.DIR(),
						status:   EJobStatusType.SUCCESS(),
					}
					if err != nil {
						w.status = EJobStatusType.FAILED()
					}

					kalpavriksha.results <- w

					// Insert this directory for further iteration to main queue
					go func() {
						kalpavriksha.jobs <- workItem{
							path:    dirPath + "/",
							objtype: EObjectType.DIR(),
							status:  EJobStatusType.WAIT(),
						}
					}()
				}
			} else {
				log.Printf("(%d) Failed to get list of blobs %s, Marker : %s\n", job.workerId, job.path, *resp.NextMarker)
				time.Sleep(5 * time.Second)
			}
		}
		atomic.AddInt64(&totalProcessedCount, int64(listCnt))
		atomic.AddInt64(&WaitCount, 1)
	}
}
