package rockgo

import (
	"testing"
	"time"
)

func TestJobPool(t *testing.T) {

	jobpool := &JobPoolExecutor{}
	jobpool.Init(10, 1024)

	var idIndex int64 = 1
	for index := 0; index < 2048; index++ {
		task := &JobTask{}
		task.Id = idIndex
		task.Tag = index
		jobpool.Put(task)
		idIndex++
	}
	time.Sleep(7 * time.Second)

	idIndex = 1
	for idIndex = 1; idIndex < 10; idIndex++ {
		jobpool.StopJobId(100 * idIndex)
	}
	jobpool.StopAll()

	jobpool.Destroy()
	time.Sleep(30 * time.Second)

}
