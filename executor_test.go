package rockgo

import (
	"testing"
	"log"
	"github.com/pingcap/tidb/_vendor/src/github.com/juju/errors"
)

func TestAsyncTaskObj(t *testing.T) {
	asynctask := AsyncTaskObj{Id: 0, Tag: "test"}

	asynctask.DoInBackgroundFunc = dobackground

	err := asynctask.Start("test_start", "test")
	if err != nil {
		log.Fatalln("启动失败", err)
	}
	log.Println(asynctask.GetStatus())
	resObj, err := asynctask.GetResult(0)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(resObj, asynctask.GetStatus())
}

func dobackground(v ...interface{}) (interface{}, error) {
	return v[0], errors.New("测试错误")
}
