package main

import (
	"bytes"
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
