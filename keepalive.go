package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
)

func HandUpProc(cfg *ProcConfig) {
	err := exec.Command("bash", "-c", cfg.KeepaliveShell).Run()
	if err != nil {
		log.Println(cfg.Name + " 拉起失败")
	} else {
		online := CheckProcess(cfg.Name)
		if online {
			log.Println(cfg.Name + " 拉起成功")
		}
	}
}

func KeepGoProcessalive() {
	pid := strconv.Itoa(os.Getpid())
	cmd := "echo -1000 > /proc/" + pid + "/oom_score_adj"
	err := exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		log.Println("监控程序保活失败")
	} else {
		log.Println("监控程序保活成功")
	}
}
