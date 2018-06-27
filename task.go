package rockgo

import "context"

type TaskStatus int

const (
	Pending           TaskStatus = 1 + iota //准备，等待
	Running                                 //正常运行中
	Finished_Normal                         //正常结束
	Finished_Error                          //出现异常，recover()返回错误
	Finished_Canceled                       //被取消结束
)

/**
任务类抽象
 */
type TaskFutureIn interface {
	Start(ctx context.Context) error //准备数据
	Do() error    //具体任务操作
	Cancel(mayInterruptIfRunning bool) error
	GetInfo() (id int64, status TaskStatus) //任务Id，当前状态
	//IsCancelled() bool
	//IsDone() bool
}
