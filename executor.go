package rockgo

import (
	"time"
	"sync"

	"strconv"
	"log"
	"errors"
	"golang.org/x/sync/syncmap"
)

const (
	Task_Status_Pending  = 0
	Task_Status_Running  = 2
	Task_Status_Finished = 4
	Task_Status_Canceled = 8
)

var ErrTimeOut error = errors.New("time out")
var ErrCanceled error = errors.New("任务已被取消")

//异步类接口
type AsyncTaskIn interface {
	Start() error
	GetInfo() (id int64, tag string, status int)
	Cancel()
	GetResult(duration time.Duration) (interface{}, error)
}

//异步任务基础类
type AsyncTaskObj struct {
	Id     int64
	Tag    string
	status int

	statusSync         *sync.RWMutex
	OnPreFunc          func() error                              //和Start()在同一协成中执行
	DoInBackgroundFunc func(...interface{}) (interface{}, error) //在另一个协成中执行

	resultChan chan interface{}
}

func (ato *AsyncTaskObj) Start() error {
	return ato.StartParams()
}

//启动任务
func (ato *AsyncTaskObj) StartParams(v ...interface{}) error {
	ato.Reset()

	if _, _, status := ato.GetInfo(); status == Task_Status_Running {
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

	if _, _, status := ato.GetInfo(); status != Task_Status_Pending {
		return errors.New("错误:任务启动错误，状态为_" + strconv.Itoa(status) + ato.Tag + "_" + strconv.FormatInt(ato.Id, 10))
	}

	ato.setStatus(Task_Status_Running)
	go ato.doInbackground(ato.resultChan, v...)
	return nil
}

//取消任务
func (ato *AsyncTaskObj) Cancel() {
	ato.setStatus(Task_Status_Finished)
}

func (ato *AsyncTaskObj) GetInfo() (int64, string, int) {
	ato.statusSync.RLock()
	status := ato.status
	ato.statusSync.RUnlock()
	return ato.Id, ato.Tag, status
}

//阻塞等待结果，若超时，返回错误 Task_Err_TimeOut
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
		err = ErrTimeOut
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
	if _, _, status := ato.GetInfo(); status == Task_Status_Finished {
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
	Cancel()
}

type TaskExecutor struct {
	waitTaskChan chan TaskIn
	allTasMap    *syncmap.Map

	maxRunChan chan struct{}
}

func (te *TaskExecutor) Start(maxnum int) error {
	te.waitTaskChan = make(chan TaskIn, 20480)
	te.maxRunChan = make(chan struct{}, maxnum)

	te.allTasMap = &syncmap.Map{}

	go te.handlerTasks()
	return nil
}

func (te *TaskExecutor) AddTask(task TaskIn) {

	id, _, _ := task.GetInfo()
	te.allTasMap.Store(id, task)
	te.waitTaskChan <- task
}

func (te *TaskExecutor) handlerTasks() {
	for {
		item, ok := <-te.waitTaskChan
		if !ok {
			break
		}
		_, _, status := item.GetInfo()
		if status == Task_Status_Finished {
			continue
		}

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
	id, _, _ := taskin.GetInfo()
	te.allTasMap.Delete(id)
}

func (te *TaskExecutor) CancelById(id int64) {
	task, ok := te.allTasMap.Load(id)
	if ok {
		task.(TaskIn).Cancel()
	}
}

func (te *TaskExecutor) CancelByTag(tag string) {
	te.allTasMap.Range(func(key, value interface{}) bool {
		if value == nil {
			return true
		}
		itemTask := value.(TaskIn)
		if _, tagItem, _ := itemTask.GetInfo(); tag == tagItem {
			itemTask.Cancel()
		}
		return true
	})
}

//任务最后状态分为 Task_Status_Canceled,Task_Status_Finished

type TaskObj struct {
	Id         int64
	Tag        string
	status     int
	statusSync *sync.RWMutex

	resChan chan interface{}

	OnPreFunc func() error
	OnDoFunc  func() (interface{}, error)
}

//若任务状态不是Task_Status_Pending，终止执行并返回错误。
//若OnDoFunc()前或执行后，任务被取消，返回错误 ErrCanceled
func (task *TaskObj) Start() error {

	if task.statusSync == nil {
		task.statusSync = &sync.RWMutex{}
	}

	_, _, status := task.GetInfo()
	if status != Task_Status_Pending {
		return errors.New("任务状态异常，终止启动" + strconv.Itoa(status) + task.Tag)
	}

	err := task.OnPreFunc()
	if err != nil {
		return err
	}

	if task.resChan == nil {
		task.resChan = make(chan interface{}, 1)
	} else {
		select {
		case _, ok := <-task.resChan:
			if !ok {
				task.resChan = make(chan interface{}, 1)
			}
		default:
		}
	}

	if _, _, status := task.GetInfo(); status != Task_Status_Pending {
		if status == Task_Status_Canceled {
			return ErrCanceled
		}
		return errors.New("任务状态异常，终止启动" + strconv.Itoa(status) + task.Tag)
	}

	task.setStatus(Task_Status_Running)
	resObj, err := task.OnDoFunc()

	if _, _, status := task.GetInfo(); status == Task_Status_Canceled {
		return ErrCanceled
	}

	if len(task.resChan) > 0 {
		select {
		case <-task.resChan:
		default:
		}
	}
	task.setStatus(Task_Status_Finished)

	if err != nil {
		select {
		case task.resChan <- err:
		default:
			return err
		}
	} else {
		select {
		case task.resChan <- resObj:
		default:
			return errors.New("把结果队列已满，无法存放结果")
		}
	}
	return nil
}

func (task *TaskObj) GetInfo() (int64, string, int) {
	task.statusSync.RLock()
	status := task.status
	task.statusSync.RUnlock()
	return task.Id, task.Tag, status
}

func (task *TaskObj) Cancel() {
	task.setStatus(Task_Status_Canceled)
	close(task.resChan)
}

//阻塞读取结果，若已取消返回错误ErrCanceled,若超时返回错误ErrTimeOut
func (task *TaskObj) GetResult(duration time.Duration) (interface{}, error) {
	if _, _, status := task.GetInfo(); status == Task_Status_Canceled {
		return nil, ErrCanceled
	}

	if duration <= 0 {
		res, ok := <-task.resChan
		if ok {
			return res, nil
		} else {
			return nil, ErrCanceled
		}
	}

	timdDu := time.After(duration)
	select {
	case res, ok := <-task.resChan:
		if ok {
			return res, nil
		} else {
			return nil, ErrCanceled
		}
	case <-timdDu:
		return nil, ErrTimeOut
	}
}

func (task *TaskObj) setStatus(code int) {
	task.statusSync.Lock()
	task.status = code
	task.statusSync.Unlock()
}
