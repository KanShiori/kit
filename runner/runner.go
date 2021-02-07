package runner

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Version = "0.1.1"
)

var (
	ErrIsRunning = errors.New("runner is running")

	// RunningRunners 记录所有运行中的 Runner, name : *Runner
	RunningRunners sync.Map = sync.Map{}
)

// IRunner 包含 Runner 提供给外部的接口, 用于接口的继承
type IRunner interface {
	// Start 启动 Runner
	Start() error

	// Stop 停止 Runner
	Stop()

	// KeepAlive 显式进行一次 keepalive
	KeepAlive()

	// IsTimeout 用于检查是否超时
	IsTimeout(curTime time.Time) bool

	// Name 返回 Runner 命名
	Name() string
}

// Runner 是封装了永久循环的 goroutine 对象
type Runner struct {
	handler Handler
	name    string

	Timeout  time.Duration
	Interval time.Duration
	Logger   io.Writer

	lastHandleTime *atomic.Value // time.Time

	// 流程控制相关
	mutex   sync.Mutex
	running bool
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

func NewRunner(handler Handler, name string, interval time.Duration) *Runner {
	r := &Runner{
		handler:        handler,
		name:           name,
		mutex:          sync.Mutex{},
		lastHandleTime: &atomic.Value{},
		Interval:       interval,
		Timeout:        time.Hour,

		running: false,
		stopCh:  make(chan struct{}),
		wg:      sync.WaitGroup{},
	}
	r.lastHandleTime.Store(time.Now())

	return r
}

// Start 开始进行永久循环执行 handle 方法
func (r *Runner) Start() error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.running {
		return ErrIsRunning
	}

	err := r.handler.OnStart()
	if err != nil {
		return err
	}

	r.running = true
	go r.run()

	// 加入记录
	RunningRunners.Store(r.name, r)

	return nil
}

// Stop 调用 OnExit 回调, 停止并等待 runner 的退出
func (r *Runner) Stop() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.running {
		return
	}

	r.handler.OnExit()

	close(r.stopCh)

	// 等待 run 退出
	r.wg.Wait()
	r.running = false

	// 清除记录
	RunningRunners.Delete(r.name)
}

// KeepAlive 刷新 LastHandleTime. 默认会在每次 handle 执行后执行
func (r *Runner) KeepAlive() {
	r.lastHandleTime.Store(time.Now())
}

// LastHandleTime 返回上一次 handle 执行结束的时间
func (r *Runner) LastHandleTime() time.Time {
	return r.lastHandleTime.Load().(time.Time)
}

// IsTimeout 用于检查 Runner 是否阻塞
func (r *Runner) IsTimeout(curTime time.Time) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	last := r.lastHandleTime.Load().(time.Time)

	return curTime.Sub(last) > r.Timeout
}

func (r *Runner) Name() string {
	return r.name
}

func (r *Runner) run() {
	defer func() {
		if x := recover(); x != nil {
			stackBuf := make([]byte, 1024*10)
			size := runtime.Stack(stackBuf, true)
			msg := fmt.Sprintf("catch panic in Runner.run {panic=%s, name=%s}: %s", x, r.name, stackBuf[0:size])

			if r.Logger != nil {
				r.Logger.Write([]byte(msg))
			}
			panic(x)
		}
	}()

	r.wg.Add(1)
	defer r.wg.Done()

	// do while
	select {
	case <-r.stopCh:
		return
	default:
		r.handler.Handle()
		r.KeepAlive()
	}

	for {
		select {
		case <-r.stopCh:
			return
		case <-time.After(r.Interval):
			r.handler.Handle()
			r.KeepAlive()
		}
	}
}
