package main

import (
	"fmt"
	"testing"
)

//func TestGetProcessesByName(t *testing.T) {
//
//	a, _ := GetProcessesByName("bs")
//	name, err := a[0].MemoryPercent()
//	if err != nil {
//		fmt.Println("err :", err)
//	}
//	fmt.Println("len ", len(a))
//	fmt.Println("a[0]: ", name)
//	assert.Equal(t, name, "init")
//}

//func TestGetProcessesByPid(t *testing.T) {
//	process, err := GetProcessesByPid(19)
//	if err != nil {
//		fmt.Println("err :", err)
//	}
//	assert.Equal(t, process.Pid, int32(19))
//}

func TestGetProcessPidByName(t *testing.T) {
	pid, _ := GetProcessPidByName("top")
	fmt.Println(*ProcCpuMonitor(*pid))
	fmt.Println(*ProcMemMonitor(*pid))
	fmt.Println(*ProcFDMonitor(*pid))
}
