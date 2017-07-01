package rockgo

import (
	"time"
	"sync"

	"strconv"
	"log"
	"errors"
)

const (
	Task_Status_Pending  = 0
	Task_Status_Running  = 2
	Task_Status_Finished = 4
)

const (
	Task_Err_TimeOut = "timeout"
)

type AsyncTaskIn interface {
	Start(...interface{}) error
	Cancel()
	GetStatus() int
	GetResult(duration time.Duration) (interface{}, error)
}

type AsyncTaskObj struct {
	Id     int64
	Tag    string
	status int

	statusSync         *sync.RWMutex
	OnPreFunc          func() error
	DoInBackgroundFunc func(...interface{}) (interface{}, error)

	resultChan chan interface{}
}

func (ato *AsyncTaskObj) Start(v ...interface{}) error {
	ato.Reset()

	if ato.GetStatus() == Task_Status_Running {
		return errors.New("错误:任务正在运行，无法启动_" + ato.Tag + "_" + strconv.FormatInt(ato.Id, 10))
	}

	ato.setStatus(Task_Status_Pending)

	if ato.OnPreFunc != nil {
		err := ato.OnPreFunc()
		if err != nil {
			return err
		}
	}

	if ato.DoInBackgroundFunc == nil {
		return errors.New("错误:任务内未实现具体后台逻辑，DoInBackgroundFunc=nil")
	}

	if ato.GetStatus() != Task_Status_Pending {
		return errors.New("错误:任务启动错误，状态为_" + strconv.Itoa(ato.GetStatus()) + ato.Tag + "_" + strconv.FormatInt(ato.Id, 10))
	}

	ato.setStatus(Task_Status_Running)
	go ato.doInbackground(ato.resultChan, v...)
	return nil
}

//取消任务
func (ato *AsyncTaskObj) Cancel() {
	ato.setStatus(Task_Status_Finished)
}

func (ato *AsyncTaskObj) GetStatus() int {
	ato.statusSync.RLock()
	status := ato.status
	ato.statusSync.RUnlock()
	return status
}

func (ato *AsyncTaskObj) GetResult(duration time.Duration) (interface{}, error) {
	var resobj interface{}
	var err error
	if duration <= 0 {
		resobj = <-ato.resultChan
		resTmp, ok := resobj.(error)
		if ok {
			return nil, resTmp
		}
		return resTmp, nil
	}

	timeOut := time.After(duration)
	select {
	case resobj = <-ato.resultChan:
	case <-timeOut:
		err = errors.New(Task_Err_TimeOut)
	}
	return resobj, err
}

func (ato *AsyncTaskObj) setStatus(code int) {
	ato.statusSync.Lock()
	ato.status = code
	ato.statusSync.Unlock()
}

func (ato *AsyncTaskObj) doInbackground(reschan chan interface{}, v ...interface{}) {
	var err error
	var resObj interface{}
	if ato.GetStatus() == Task_Status_Finished {
		err = errors.New("任务已被取消_" + ato.Tag + "_" + strconv.FormatInt(ato.Id, 10))
		select {
		case reschan <- err:
		default:
			log.Println("异常：无法把结果放到任务结果队列")
			//errchan <- errors.New("异常：无法把结果放到任务结果队列")
		}
		return
	}

	ato.setStatus(Task_Status_Running)
	resObj, err = ato.DoInBackgroundFunc(v...)
	if err != nil {
		resObj = err
	}

	select {
	case reschan <- resObj:
	default:
		log.Println("异常：接收任务结果的channel已满，无法发送数据")
	}
	ato.setStatus(Task_Status_Finished)
}

//运行之前重置数据
func (ato *AsyncTaskObj) Reset() {

	if ato.statusSync == nil {
		ato.statusSync = &sync.RWMutex{}
	}

	if ato.resultChan == nil {
		ato.resultChan = make(chan interface{}, 1)
	} else {
		resChanLen := len(ato.resultChan)
		if resChanLen > 0 {
			for index := 0; index < resChanLen; index++ {
				<-ato.resultChan
			}
		}
	}
}

type TaskIn interface {
	Start() error
	GetInfo() (id int64, tag string, status int)
}

type TaskExecutor struct {
	waitTaskChan   chan TaskIn
	runningTaskMap map[int64]TaskIn

	maxRunChan  chan struct{}
	taskTimeout <-chan time.Time
}

func (te *TaskExecutor) Start(maxnum int, timeout time.Duration) error {
	te.waitTaskChan = make(chan TaskIn, 20480)
	te.maxRunChan = make(chan struct{}, maxnum)
	te.taskTimeout = time.After(timeout)
	go te.handlerTasks()
	return nil
}

func (te *TaskExecutor) AddTask(task TaskIn) {
	te.waitTaskChan <- task
}

func (te *TaskExecutor) handlerTasks() {
	for {
		item, ok := <-te.waitTaskChan
		if !ok {
			break
		}
		id, _, status := item.GetInfo()
		if status == Task_Status_Finished {
			continue
		}
		te.runningTaskMap[id] = item
		te.maxRunChan <- struct{}{}

		go te.startTask(item)
	}
}

func (te *TaskExecutor) startTask(taskin TaskIn) {
	taskin.Start()
	select {
	case <-te.maxRunChan:
	default:
	}
}
