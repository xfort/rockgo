package rockgo

import (
	"github.com/pingcap/tidb/_vendor/src/github.com/juju/errors"
	"time"
	"fmt"
	"log"
)

const Job_Status_IDLE = 10
const Job_Status_Running = 20

const Job_Status_Stopped = 30
const Job_Status_Destroyed = 31

type JobIn interface {
	Start(resCallback ResCallback)
	Stop()
	GetInfo() (jobId int, statusCode int, jobTag int)
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
}

func (jobpool *JobPoolExecutor) Init(maxCoreSize int, maxLength int) {

	jobpool.jobQueue = make(chan JobIn, maxLength)
	jobpool.coreChan = make(chan struct{}, maxCoreSize)

	jobpool.jobsMap = &JobSyncMap{}

	go jobpool.start()
}

func (jobpool *JobPoolExecutor) start() {
	for {
		if jobpool.StatusCode == Job_Status_Destroyed {
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
	if jobpool.StatusCode == Job_Status_Destroyed {
		return errors.New("jobpool已被销毁")
	}
	go jobpool.putJob(job)
	return nil
}

func (jobpool *JobPoolExecutor) putJob(job JobIn) {
	jobpool.jobQueue <- job
	jobId, _, _ := job.GetInfo()
	jobpool.jobsMap.Store(jobId, job)

	log.Println("放入任务队列", jobId)
}

func (jobpool *JobPoolExecutor) StopAll() {
	jobpool.jobsMap.Range(func(key int, value JobIn) bool {
		jobpool.jobsMap.Delete(key)
		return true
	})

}

func (jobpool *JobPoolExecutor) StopJobId(jobId int) {
	jobtask, ok := jobpool.jobsMap.Load(jobId)
	if !ok || jobtask == nil {
		return
	}
	jobtask.Stop()
	//jobpool.jobsMap
}

func (jobpool *JobPoolExecutor) StopJobTag(jobTag int) {
	jobpool.jobsMap.Range(func(key int, value JobIn) bool {
		_, _, itemTag := value.GetInfo()
		if itemTag == jobTag {
			value.Stop()
		}
		return true
	})
}

func (jobpool *JobPoolExecutor) DestroyNow() {
	if jobpool.StatusCode == Job_Status_Destroyed {
		return
	}
	close(jobpool.jobQueue)

	jobpool.StatusCode = Job_Status_Destroyed
}

type JobTask struct {
	Id         int
	Tag        int
	StatusCode int
}

func (jobtask *JobTask) Start(rescallback ResCallback) {
	log.Println("任务", jobtask.Id, "开始")
	time.Sleep(3 * time.Second)
	rescallback()
	fmt.Println("任务", jobtask.Id, "结束")
}

func (task *JobTask) Stop() {
	task.StatusCode = Job_Status_Stopped
}

func (task *JobTask) GetInfo() (id int, status int, tag int) {
	return task.Id, task.StatusCode, task.Tag
}
