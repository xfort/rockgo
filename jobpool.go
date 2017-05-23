package rockgo

import (
	"github.com/pingcap/tidb/_vendor/src/github.com/juju/errors"
	"time"
	"fmt"
	"log"
	"sync"
)

const Job_Status_IDLE = 10
const Job_Status_Running = 20

const Job_Status_Stopped = 30
const Job_Status_Destroyed = 31

type JobIn interface {
	Start(resCallback ResCallback)
	Stop()
	GetInfo() (jobId int64, statusCode int, jobTag int)
}

type ResCallback func()

func (jobpool *JobPoolExecutor) jobCallback() {
	<-jobpool.coreChan
}

type JobPoolExecutor struct {
	jobQueue chan JobIn

	coreChan    chan struct{}
	maxCoreSize int //并发数
	maxLength   int //队列最大数

	jobsMap *JobSyncMap

	StatusCode int

	statusMutex *sync.RWMutex
}

func (jobpool *JobPoolExecutor) Init(maxCoreSize int, maxLength int) {

	jobpool.coreChan = make(chan struct{}, maxCoreSize)

	jobpool.jobQueue = make(chan JobIn, maxLength)

	jobpool.jobsMap = &JobSyncMap{}

	jobpool.statusMutex = new(sync.RWMutex)

	go jobpool.start()
}

func (jobpool *JobPoolExecutor) start() {
	for {
		if jobpool.StatusCode == Job_Status_Destroyed || jobpool.jobQueue == nil {
			fmt.Println("工作池chan已销毁")
			break
		}

		jobItem, ok := <-jobpool.jobQueue
		if !ok {
			fmt.Println("工作池chan已关闭")
			break
		}

		if jobpool.StatusCode == Job_Status_Destroyed {
			fmt.Println("工作池chan已销毁")
			break
		}

		jobId, statusCode, jobTag := jobItem.GetInfo()
		if statusCode == Job_Status_Stopped || jobId < 0 || jobTag < 0 {

			fmt.Println("任务已终止", jobId, statusCode)

			continue
		}

		var tmpObj struct{}
		jobpool.coreChan <- tmpObj

		go jobItem.Start(jobpool.jobCallback)
	}
}

func (jobpool *JobPoolExecutor) Put(job JobIn) error {
	if jobpool.StatusCode == Job_Status_Destroyed || jobpool.jobQueue == nil {
		return errors.New("jobpool已被销毁")
	}

	go jobpool.putJob(job)
	return nil
}

func (jobpool *JobPoolExecutor) putJob(job JobIn) {
	jobpool.statusMutex.RLock()
	if jobpool.StatusCode == Job_Status_Destroyed {
		jobpool.statusMutex.RUnlock()
		return
	}
	jobpool.statusMutex.RUnlock()

	if jobpool.jobQueue != nil {
		jobpool.jobQueue <- job

		if jobpool.jobsMap != nil {
			jobId, _, _ := job.GetInfo()
			jobpool.jobsMap.Store(jobId, job)
		}
	}
}

func (jobpool *JobPoolExecutor) GetAllJob() *JobSyncMap {
	return jobpool.jobsMap
}

func (jobpool *JobPoolExecutor) StopAll() {
	if jobpool.StatusCode == Job_Status_Destroyed {
		return
	}
	jobpool.jobsMap.Range(func(key int64, value JobIn) bool {
		value.Stop()
		//jobpool.jobsMap.Delete(key)
		return true
	})
}

func (jobpool *JobPoolExecutor) StopJobId(jobId int64) {
	if jobpool.StatusCode == Job_Status_Destroyed {
		return
	}
	jobtask, ok := jobpool.jobsMap.Load(jobId)
	if !ok || jobtask == nil {
		return
	}
	jobtask.Stop()
}

func (jobpool *JobPoolExecutor) StopJobTag(jobTag int) {

	if jobpool.StatusCode == Job_Status_Destroyed {
		return
	}

	jobpool.jobsMap.Range(func(key int64, value JobIn) bool {
		_, _, itemTag := value.GetInfo()
		if itemTag == jobTag {
			value.Stop()
		}
		return true
	})
}

func (jobpool *JobPoolExecutor) Destroy() {
	jobpool.statusMutex.RLock()
	if jobpool.StatusCode == Job_Status_Destroyed {
		jobpool.statusMutex.RUnlock()
		return
	}
	jobpool.statusMutex.RUnlock()

	jobpool.statusMutex.Lock()
	jobpool.StatusCode = Job_Status_Destroyed
	jobpool.statusMutex.Unlock()

	jobpool.StopAll()

	close(jobpool.jobQueue)

	jobpool.jobQueue = nil
	//jobpool.jobsMap = nil
}

type JobTask struct {
	Id         int64
	Tag        int
	StatusCode int
}

func (jobtask *JobTask) Start(rescallback ResCallback) {
	log.Println("任务", jobtask.Id, "开始")

	time.Sleep(1 * time.Second)
	rescallback()
	fmt.Println("任务", jobtask.Id, "结束")
}

func (task *JobTask) Stop() {
	task.StatusCode = Job_Status_Stopped
}

func (task *JobTask) GetInfo() (id int64, status int, tag int) {
	return task.Id, task.StatusCode, task.Tag
}
