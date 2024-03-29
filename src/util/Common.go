package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
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

func Log(level string, message string) {
	time := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [%s]: %s", time, level, message)
}

func LogInfo(message string) {
	Log("INFO", message)
}

func LogWarning(message string) {
	Log("WARNING", message)
}

func LogError(message string) {
	Log("ERROR", message)
}

func LogDebug(message string) {
	Log("DEBUG", message)
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

func PrintUptime() {
	uptime := GetUptime()
	log.Printf("Uptime: %v", FormatDuration(uptime))
}

func PrintAmountRequests() {
	amount := GetAmountRequests()
	log.Printf("Amount Requests: %d", amount)
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

func FormatDuration(milliseconds int64) string {
	duration := time.Duration(milliseconds) * time.Millisecond

	days := int(duration / (24 * time.Hour))
	duration = duration % (24 * time.Hour)

	hours := int(duration / time.Hour)
	duration = duration % time.Hour

	minutes := int(duration / time.Minute)
	duration = duration % time.Minute

	seconds := int(duration / time.Second)

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dmin ", minutes)
	}
	result += fmt.Sprintf("%ds", seconds)

	return result
}
