package rockgo

type TaskPool struct {
	maxNum  int           //最大并发数
	numChan chan struct{} //用于控制并发数量
}

func NewTaskPool(maxNum int) *TaskPool {
	taskpool := &TaskPool{maxNum: maxNum, numChan: make(chan struct{}, maxNum)}
	return taskpool
}

/**
设置最大并发数
 */
func (pool *TaskPool) SetMaxNum(maxNum int) {
	if pool.numChan == nil {
		pool.maxNum = maxNum
		pool.numChan = make(chan struct{}, maxNum)
	} else {
		pool.maxNum = maxNum
		oldChan := pool.numChan
		pool.numChan = make(chan struct{}, maxNum)
		if len(oldChan) > 0 {
			for {
				select {
				case _, ok := <-oldChan:
					if ok {
						select {
						case pool.numChan <- struct{}{}:
						default:
						}
					} else {
						return
					}
				default:
					return
				}
			}
		}
	}
}

/**
添加了新任务
 */
func (pool *TaskPool) Add(num int) {
	for index := 0; index < num; index++ {
		pool.numChan <- struct{}{}
	}
}

/**
任务执行完成
 */
func (pool *TaskPool) Done(num int) {
	for index := 0; index < num; index++ {
		select {
		case <-pool.numChan:
		default:
		}
	}
}

func (pool *TaskPool) Destroy() {
	if pool.numChan != nil {
		close(pool.numChan)
	}
}
