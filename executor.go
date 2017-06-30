package rockgo

import (
	"time"
	"sync"
)

const (
	Task_Pending  = 0
	Task_Running  = 2
	Task_Finished = 4
)

type TaskIn interface {
	Start() error
	Stop()
	GetStatus() (status int)
	GetRes(time.Duration) (interface{}, error)
}

type FutureTaskObj struct {
	Id  int64
	Tag string

	Status     int
	statusSync *sync.RWMutex
	resChan    chan interface{}
}

func (ft *FutureTaskObj) Start() error {
	ft.statusSync = new(sync.RWMutex)

	ft.statusSync.Lock()
	ft.Status = Task_Running
	ft.statusSync.Unlock()
	return nil
}

func (ft *FutureTaskObj) Stop() {
	ft.statusSync.Lock()
	ft.Status = Task_Finished
	ft.statusSync.Unlock()
}

func (ft *FutureTaskObj) GetStatus() (int) {
	ft.statusSync.RLock()
	status := ft.Status
	ft.statusSync.RUnlock()
	return status
}

func (ft *FutureTaskObj) GetRes(duration time.Duration) (interface{}, error) {
	if duration <= 0 {
	}
	return nil, nil
}

//结束
func (ft *FutureTaskObj) onFinished(resobj interface{}, err error) {

}
