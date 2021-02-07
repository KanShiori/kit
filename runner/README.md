# runner

runner 用于构建长久运行的 goroutine 对象, 并提供 Start, Stop, Keepalive 等功能.


## Usage
```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/KanShiori/kit/runner"
)

type Demo struct {
	*runner.Runner
}

func NewDemo() *Demo {
	d := &Demo{}

	d.Runner = runner.NewRunner(d, "Demo", time.Second)
	return d
}

func (d *Demo) Handle() {
	fmt.Println("is in handle")
}

func (d *Demo) OnStart() error {
	return nil
}

func (d *Demo) OnExit() {
}

func main() {
	d := NewDemo()

	// - start runner
	err := d.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	// - handling ...
	time.Sleep(10 * time.Second)

	// - stop runner
	d.Stop()
}
```