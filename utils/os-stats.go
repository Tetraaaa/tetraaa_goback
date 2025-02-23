package utils

import (
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/prometheus/procfs"
)

func GetCPUTemp() float64 {
	fallbackValue := 44.3
	if runtime.GOOS == "windows" {
		return fallbackValue
	}
	cmd := exec.Command("vcgencmd", "measure_temp")
	stdout, err := cmd.Output()
	if err != nil {
		log.Print(err)
		return fallbackValue
	}

	tempString := string(stdout)
	tempString = strings.ReplaceAll(tempString, "temp=", "")
	tempString = strings.ReplaceAll(tempString, "'C", "")
	tempString = strings.ReplaceAll(tempString, "\n", "")
	temp, err := strconv.ParseFloat(tempString, 64)
	if err != nil {
		log.Print(err)
		return fallbackValue
	}
	return temp
}

func GetCPUSAverages() []int64 {
	fallbackValue := []int64{0, 0, 0, 0}
	if runtime.GOOS == "windows" {
		return fallbackValue
	}

	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return fallbackValue
	}
	stats, err := fs.Stat()
	if err != nil {
		return fallbackValue
	}
	resp := []int64{}
	for cpu := range stats.CPU {
		resp = append(resp, cpu)
	}
	return resp

}

func GetMemoryTotalAndFree() (uint64, uint64) {
	var fallbackValue1, fallbackValue2 uint64 = 4294967296, 4294967296

	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return fallbackValue1, fallbackValue2
	}
	stats, err := fs.Meminfo()
	if err != nil {
		return fallbackValue1, fallbackValue2
	}
	return *stats.MemTotal, *stats.MemAvailable

}

type PeribotResponse struct {
	Status string `json:"status"`
	Uptime int64  `json:"uptime"`
}
