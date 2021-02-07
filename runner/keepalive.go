package runner

import "time"

// CheckaliveForever 用于周期性对所有 runner keepalive
func CheckaliveForever(interval time.Duration, handleTimeout func(name string, runner *Runner, curTime time.Time)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		RunningRunners.Range(
			func(key, value interface{}) bool {
				name := key.(string)
				runner := value.(*Runner)

				curTime := time.Now()
				if runner.IsTimeout(curTime) {
					handleTimeout(name, runner, curTime)
				}

				return true
			},
		)
	}

}
