package rockgo

import (
	"errors"
	"context"
	"time"
	"sync"
	"fmt"
	"log"
)

const (
	Task_Status_Pending           = 0
	Task_Status_Running           = 1 //正常运行中
	Task_Status_Finished_Normal   = 2 //正常结束
	Task_Status_Finished_Panic    = 3 //出现异常，recover()返回错误
	Task_Status_Finished_Canceled = 4 //被取消结束
)

type TaskIn interface {
	GetInfo() (id int, tag string)
	GetStatus() int
	StartContext(ctx context.Context, v ...interface{}) (interface{}, error)
	Cancel() error
	IsDone() (bool, error)
	GetResult(time.Duration) (interface{}, error)
}

type DoWorkFunc func(context.Context, ...interface{}) (interface{}, error)

type TaskObj struct {
	Id         int
	Tag        string
	status     int
	statusSync *sync.RWMutex

	CtxCancelFunc context.CancelFunc
	resChan       chan interface{}

	DoFunc DoWorkFunc

	doneChan <-chan struct{}
}

func (taskobj *TaskObj) GetInfo() (int, string) {
	return taskobj.Id, taskobj.Tag
}

func (taskobj *TaskObj) GetStatus() int {
	taskobj.statusSync.RLock()
	status := taskobj.status
	taskobj.statusSync.RUnlock()
	return status
}

func (taskobj *TaskObj) setStatus(status int) {
	taskobj.statusSync.Lock()
	taskobj.status = status
	taskobj.statusSync.Unlock()
}

func (taskobj *TaskObj) StartContext(parentCtx context.Context, v ...interface{}) (error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("任务发生异常错误", err, taskobj.Id, taskobj.Tag)
			taskobj.setStatus(Task_Status_Finished_Panic)
		}
	}()

	if taskobj.statusSync == nil {
		taskobj.statusSync = &sync.RWMutex{}
	}

	if status := taskobj.GetStatus(); status != Task_Status_Pending {
		return taskobj.err("任务状态异常,不是初始状态,停止运行")
	}

	taskobj.setStatus(Task_Status_Running)

	if taskobj.resChan == nil {
		taskobj.resChan = make(chan interface{}, 1)
	} else {
		chanLne := len(taskobj.resChan)
		if chanLne > 0 {
			for index := 0; index < chanLne; index++ {
				select {
				case _, ok := <-taskobj.resChan:
					if !ok {
						break
					}
				default:
				}
			}
		}
		select {
		case _, ok := <-taskobj.resChan:
			if !ok {
				taskobj.resChan = make(chan interface{}, 1)
			}
		default:
		}
	}
	taskobj.doneChan = parentCtx.Done()

	go taskobj.doworkfunc(parentCtx, taskobj.resChan, v...)

	return nil
}

func (taskobj *TaskObj) Cancel() error {
	if taskobj.CtxCancelFunc != nil {
		taskobj.CtxCancelFunc()
	}
	if taskobj.GetStatus() == Task_Status_Running {
		taskobj.setStatus(Task_Status_Finished_Canceled)
	}
	return nil
}

func (taskobj *TaskObj) IsDone() (bool, error) {
	if taskobj.GetStatus() >= Task_Status_Finished_Normal {
		return true, nil
	}
	return false, nil
}

func (taskobj *TaskObj) GetResult(dura time.Duration) (interface{}, error) {
	timeDura := time.After(dura)
	var resObj interface{}
	var err error
	var ok bool
	select {
	case <-taskobj.doneChan:
		err = taskobj.err("任务被取消或者超时")
		if taskobj.GetStatus() == Task_Status_Running {
			taskobj.setStatus(Task_Status_Finished_Canceled)
		}
	case resObj, ok = <-taskobj.resChan:
		if !ok {
			err = taskobj.err("任务结果chan被提前关闭")
			if taskobj.GetStatus() == Task_Status_Running {
				taskobj.setStatus(Task_Status_Finished_Normal)
			}
		} else {
			resObjTmp, ok := resObj.(error)
			if ok {
				err = resObjTmp
				resObj = nil
			}
			if taskobj.GetStatus() == Task_Status_Running {
				taskobj.setStatus(Task_Status_Finished_Normal)
			}
		}
	case <-timeDura:
		err = taskobj.err("任务超时," + dura.String())
	}
	return resObj, err
}

func (taskobj *TaskObj) doworkfunc(ctx context.Context, reschan chan interface{}, v ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("任务异常panic", err)
			taskobj.setStatus(Task_Status_Finished_Panic)
		}
	}()
	if taskobj.GetStatus() != Task_Status_Running {
		return
	}

	resObj, err := taskobj.DoFunc(ctx, v...)
	log.Println("doworkfunc()", taskobj.GetStatus())

	if taskobj.GetStatus() == Task_Status_Running {
		if err != nil {
			resObj = err
		}
		close(reschan)
		select {
		case reschan <- resObj:
		default:
			log.Println("向结果队列添加数据失败", len(reschan), cap(reschan))
		}
		taskobj.setStatus(Task_Status_Finished_Normal)
	}
}

func (taskobj *TaskObj) err(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...) + fmt.Sprintf("%d_%s_%d", taskobj.Id, taskobj.Tag, taskobj.status))
}

//类似java线程池，控制最大并发数
type TaskPoolObj struct {
	numChan chan struct{}
	MaxCore int //最大并发数
}

func (taskpool *TaskPoolObj) SetData(maxCore int) {
	taskpool.MaxCore = maxCore
	taskpool.numChan = make(chan struct{}, maxCore)
}

func (taskpool *TaskPoolObj) Add(num int) error {
	var errmsg string
	for index := 0; index < num; index++ {
		select {
		case taskpool.numChan <- struct{}{}:
		default:
			errmsg = errmsg + "任务池chan异常，添加失败"
		}
	}
	if errmsg != "" {
		return errors.New(errmsg)
	}
	return nil
}

func (taskpool *TaskPoolObj) AddAsync(num int) {
	go taskpool.Add(num)
}

func (taskpool *TaskPoolObj) Done(num int) {
	if taskpool.numChan == nil {
		return
	}

	for index := 0; index < num; index++ {
		select {
		case <-taskpool.numChan:
		default:

		}
	}
}

func (taskpool *TaskPoolObj) Destroy() {
	close(taskpool.numChan)
	taskpool.numChan = nil
}
