package util

import (
	"runtime"
	"time"
)

var (
	startupTime    int64 = 0
	amountRequests int64 = 0
)

func SetStartupTime(t int64) {
	if startupTime == 0 {
		startupTime = t
	}
}

func GetUptime() int64 {
	return time.Now().UnixMilli() - startupTime
}

func GetAmountRequests() int64 {
	return amountRequests
}

func IncreamentAmountRequests() {
	amountRequests++
}

func GetMemoryStats() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem
}
