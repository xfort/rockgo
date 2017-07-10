package rockgo

import (
	"testing"
	"context"
	"errors"
	"log"
	"time"
)

func TestTaskObj(t *testing.T) {
	taskobj := TaskObj{}
	taskobj.Id = 1
	taskobj.Tag = "test"
	taskobj.DoFunc = taskobjDoFunc

	ctx, cancelFunc := context.WithCancel(context.Background())
	taskobj.CtxCancelFunc = cancelFunc

	err := taskobj.StartContext(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	resobj, err := taskobj.GetResult(7 * time.Second)
	if err != nil {
		log.Println("task_res_err", err)
	}

	log.Println("task_res", resobj, taskobj.GetStatus())

	taskobj.Cancel()
	resobj, err = taskobj.GetResult(0 * time.Second)
	if err != nil {
		log.Println(err)
	}
	log.Println("task_over_2", resobj, taskobj.GetStatus())
}

func taskobjDoFunc(ctx context.Context, v ...interface{}) (interface{}, error) {
	time.Sleep(5 * time.Second)
	return "result ok", errors.New("错误，test， dofunc")
}
