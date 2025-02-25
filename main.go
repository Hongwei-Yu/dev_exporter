package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func main() {
	PCcfg, ProCcfg, Selfcfg := loadMonitorExporterConfig()
	log.Println(Selfcfg.Version, "write by YHW")
	initPrometheus(PCcfg, ProCcfg)
	//log.Println(PCcfg, ProCcfg, Selfcfg)
	log.Println("监控程序开始")
	if Selfcfg.MonitorKeepalive {
		log.Println("开启本程序保活")
		KeepGoProcessalive()
		time.Sleep(time.Second)
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(Selfcfg.Port, nil)
		if err != nil {
			log.Panicln(err)
		}
	}()
	// 启动PC监控协程
	if PCcfg.CPU.Enabled {
		go monitorCPU(PCcfg.CPU)
	}
	if PCcfg.Memory.Enabled {
		go monitorMemory(PCcfg.Memory)
	}
	if PCcfg.Disk.Enabled {
		go monitorDisk(PCcfg.Disk)
	}
	if PCcfg.Network.Enabled {
		go monitorNetwork(PCcfg.Network)
	}
	// 启动进程监控协程
	if ProCcfg.Enabled {
		for _, proc := range ProCcfg.Process {
			go monitorProc(proc)
		}
	}

	//启动定时任务
	ScheduleCfg := loadScheduleTaskConfig()
	log.Println("存在以下定时任务")
	for _, job := range ScheduleCfg.Job {
		log.Println(job.Name + " : " + job.Cron)
		log.Println(job.Parameters)
	}
	log.Println("等待3m请确认定时任务是否正确")
	time.Sleep(180 * time.Second)
	log.Println("定时任务开始")
	if ScheduleCfg.Enabled {
		for _, job := range ScheduleCfg.Job {
			switch job.Type {
			case "shell":
				go StartJob(job.Cron, *CreateShellJob(job))

			}
		}
	}

	select {} // 阻塞主线程

}
