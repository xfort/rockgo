package rockgo

import (
	"testing"
	"log"
	"time"
	"strconv"
)

func TestAsyncTaskObj(t *testing.T) {
	asynctask := AsyncTaskObj{Id: 0, Tag: "test"}

	asynctask.DoInBackgroundFunc = dobackground

	err := asynctask.Start()

	if err != nil {
		log.Fatalln("启动失败", err)
	}

	//log.Println(asynctask.GetStatus())

	resObj, err := asynctask.GetResult(0)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resObj)
}

func dobackground(v ...interface{}) (interface{}, error) {
	return v[0], nil
}

func TestTaskExecutor(t *testing.T) {

	taskExecutor := &TaskExecutor{}
	taskExecutor.Start(3)

	for index := 0; index < 10; index++ {
		taskobj := &TaskObj{}
		taskobj.Id = time.Now().UTC().UnixNano()
		taskobj.Tag = "test" + strconv.Itoa(index)

		//taskExecutor.AddTask(taskobj)
	}

	time.Sleep(10 * time.Minute)
}
