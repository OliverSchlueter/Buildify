package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	// foregrounds
	ColorReset  = "\033[0m"
	ColorWhite  = "\033[37m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	// backgrounds
	ColorRedBg    = "\033[41m"
	ColorGreenBg  = "\033[42m"
	ColorYellowBg = "\033[43m"
)

var (
	startupTime int64 = 0
)

func SetStartupTime(t int64) {
	if startupTime == 0 {
		startupTime = t
	}
}

func GetUptime() int64 {
	return time.Now().UnixMilli() - startupTime
}

func FastCopyFile(source, destination *os.File) error {
	buf := make([]byte, 4096)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func GetArtifactFileName(filePath string) string {
	return path.Base(filePath)
}

func GetMemoryStats() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem
}

func PrintMemUsage() {
	mem := GetMemoryStats()

	log.Printf(
		"Memory usage: Alloc=%v | TotalAlloc=%v | Sys=%v | NumGC=%v",
		FormatMemory(mem.Alloc),
		FormatMemory(mem.TotalAlloc),
		FormatMemory(mem.Sys),
		mem.NumGC,
	)
}

func FormatMemory(bytes uint64) string {
	const (
		kB = 1024
		mB = kB * 1024
	)

	switch {
	case bytes >= mB:
		return fmt.Sprintf("%.2f MiB", float64(bytes)/float64(mB))
	case bytes >= kB:
		return fmt.Sprintf("%.2f KiB", float64(bytes)/float64(kB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
