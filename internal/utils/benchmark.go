package utils

import (
	"fmt"
	"time"
)

func Timer(out *string) func() {
	start := time.Now()
	return func() {
		perf := time.Since(start)
		if perf.Milliseconds() > 0 {
			*out = fmt.Sprintf("%d ms", perf.Milliseconds())
		} else {
			*out = fmt.Sprintf("%d us", perf.Microseconds())
		}
	}
}
