package runner

import "time"

type Handler interface {
	// Handle 为 Runner 中循环执行的函数
	Handle()

	// OnStart 在 Runner 启动前回调
	OnStart() error

	// OnExit 在 Runner 退出前回调
	OnExit()
}

type NoopHandler struct {
}

func (h *NoopHandler) Handle() {
	time.Sleep(30 * time.Minute)
}

func (h *NoopHandler) OnStart() error {
	return nil
}

func (h *NoopHandler) OnExit() {
}
