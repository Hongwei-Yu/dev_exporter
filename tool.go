package main

import (
	"bytes"
	"github.com/shirou/gopsutil/process"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// 进程在线检测
func CheckProcess(processName string) bool {
	// 执行任务以获取当前所有进程
	cmd := exec.Command("ps", "-e")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
		return false
	}

	// 将输出转换为字符串，并查找目标进程
	outputString := string(output)
	return strings.Contains(outputString, processName)
}

// v1.2 版本细化进程监控，添加进程cpu、内存使用率

// cpu使用率
func ProcCpuMonitor(pid int32) *float64 {

	proc, err := GetProcessesByPid(pid)
	if err != nil {
		log.Println("Get proc err:", err)
		return nil
	}
	percent, err := proc.CPUPercent()
	if err != nil {
		log.Println("proc cpu monitor err:", err)
	}
	return &percent
}

// 内存使用率
func ProcMemMonitor(pid int32) *float64 {
	proc, err := GetProcessesByPid(pid)
	if err != nil {
		log.Println("Get proc err:", err)
		return nil
	}
	percent, err := proc.MemoryPercent()
	if err != nil {
		log.Println("proc cpu monitor err:", err)
	}
	percent64 := float64(percent)
	return &percent64
}

// 打开文件描述符
func ProcFDMonitor(pid int32) *float64 {
	proc, err := GetProcessesByPid(pid)
	if err != nil {
		log.Println("Get proc err:", err)
		return nil
	}
	files, err := proc.OpenFiles()
	if err != nil {
		log.Println("proc cpu monitor err:", err)
	}
	//fmt.Println("percent: ", files)
	num := float64(len(files))
	return &num
}

// 磁盘io
//func ProcDiskMonitor(pid int32) *float64 {
//	proc, err := GetProcessesByPid(pid)
//	if err != nil {
//		log.Println("Get proc err:", err)
//		return nil
//	}
//	io, err := proc.IOCounters()
//	if err != nil {
//		log.Println("proc cpu monitor err:", err)
//	}
//	log.Println("percent: ", io.ReadBytes)
//	num := float64(len(percent))
//	return &num
//}

// 根据进程名称获取进程对象
func GetProcessesByName(name string) ([]*process.Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var result []*process.Process
	for _, p := range processes {
		pName, err := p.Name()
		if err != nil {
			continue
		}
		if pName == name {
			result = append(result, p)
		}
	}
	return result, nil
}

// 根据进程pid获取进程对象
func GetProcessesByPid(pid int32) (*process.Process, error) {
	processes, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	return processes, nil

}

// 进程在线检测
func GetProcessPidByName(processName string) (*int32, error) {

	// 执行任务以获取当前所有进程
	cmd := exec.Command("pgrep", processName)
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
		log.Println("out:", output)
		return nil, err
	}
	log.Println("out:", string(output))

	// 将输出转换为字符串，并查找目标进程
	pid, err := strconv.ParseInt(strings.Replace(string(output), "\n", "", -1), 10, 32)
	if err != nil {
		log.Println("转化错误:", err)
		return nil, err
	}

	pid32 := int32(pid)
	return &pid32, nil
}
