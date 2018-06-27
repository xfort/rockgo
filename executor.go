package rockgo

import (
	"errors"
	"context"
	"sync"
	"strconv"
)


type DoWorkFunc func(ctx context.Context, v ...interface{}) (interface{}, error)

//类似java线程池，控制最大并发数
type TaskPoolObj struct {
	numChan chan struct{}
	MaxCore int //最大并发数
}

func (taskpool *TaskPoolObj) SetData(maxCore int) {
	taskpool.MaxCore = maxCore
	taskpool.numChan = make(chan struct{}, maxCore)
}

//增加任务数，若任务已满，则阻塞等待
func (taskpool *TaskPoolObj) Add(num int) {
	for index := 0; index < num; index++ {
		taskpool.numChan <- struct{}{}
	}

}

//在任务完成时，必须执行此方法
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

type TaskSyncObj struct {
	Id         int
	Tag        string
	DoWorkFunc DoWorkFunc

	CtxCancelFunc context.CancelFunc
	CtxDone       <-chan struct{}

	statusCode int
	statusSync *sync.RWMutex
}

func (task *TaskSyncObj) GetInfo() (int, string) {
	return task.Id, task.Tag
}

func (task *TaskSyncObj) GetStatus() int {
	task.statusSync.RLock()
	var status = task.statusCode
	task.statusSync.RUnlock()
	return status
}

func (task *TaskSyncObj) Cancel() error {
	if task.CtxCancelFunc != nil {
		task.CtxCancelFunc()
	}
	if task.GetStatus() < Task_Status_Finished_Normal {
		task.setStatus(Task_Status_Finished_Canceled)
	}
	return nil
}

func (task *TaskSyncObj) StartContext(ctx context.Context, v ...interface{}) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			if task.GetStatus() != Task_Status_Finished_Panic {
				task.setStatus(Task_Status_Finished_Panic)
			}
		}
	}()

	if task.statusSync == nil {
		task.statusSync = &sync.RWMutex{}
	}

	if code := task.GetStatus(); code != Task_Status_Pending {
		task.setStatus(Task_Status_Finished_Panic)
		return nil, errors.New("任务启动失败，初始状态异常_" + strconv.Itoa(code))
	}
	if task.DoWorkFunc == nil {
		task.setStatus(Task_Status_Finished_Panic)
		return nil, errors.New("任务启动失败，DoWorkFunc为空")
	}
	task.CtxDone = ctx.Done()

	task.setStatus(Task_Status_Running)
	resData, err := task.DoWorkFunc(ctx, v...)

	select {
	case <-ctx.Done():
		if task.GetStatus() < Task_Status_Finished_Normal {
			task.setStatus(Task_Status_Finished_Canceled)
		}
		if err == nil {
			err = ctx.Err()
		} else {
			err = errors.New(ctx.Err().Error() + "__" + err.Error())
		}
	default:
	}

	if task.GetStatus() < Task_Status_Finished_Normal {
		task.setStatus(Task_Status_Finished_Normal)
	}
	return resData, err
}

func (task *TaskSyncObj) IsDone() (bool, error) {
	if task.CtxDone != nil {
		select {
		case <-task.CtxDone:
			if task.GetStatus() < Task_Status_Finished_Normal {
				task.setStatus(Task_Status_Finished_Canceled)
			}
		default:
		}
	}

	if task.GetStatus() >= Task_Status_Finished_Normal {
		return true, nil
	}
	return false, nil
}

func (task *TaskSyncObj) setStatus(code int) {
	task.statusSync.Lock()
	task.statusCode = code
	task.statusSync.Unlock()
}
