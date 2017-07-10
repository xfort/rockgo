package rockgo

import (
	"testing"
	"context"
	"github.com/pingcap/tidb/_vendor/src/github.com/juju/errors"
	"log"
	"time"
)

func TestTaskObj(t *testing.T) {
	taskobj := TaskObj{}
	taskobj.Id = 1
	taskobj.Tag = "test"
	taskobj.DoFunc = taskobjDoFunc

	ctx, _ := context.WithCancel(context.Background())

	err := taskobj.StartContext(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	resobj, err := taskobj.GetResult(10 * time.Second)
	if err != nil {
		log.Println("task_res_err", err)
	}
	log.Println("task_res", resobj)
}

func taskobjDoFunc(ctx context.Context, v ...interface{}) (interface{}, error) {
	time.Sleep(5 * time.Second)
	return nil, errors.New("错误，test， dofunc")
}
