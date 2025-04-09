package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupSignalHandler(mgr *GoroutineManager) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			select {
			case sig := <-sigChan:
				switch sig {
				case syscall.SIGHUP:
					log.Println("Received SIGHUP, reloading config...")
					Restart(mgr)

				case syscall.SIGTERM, syscall.SIGINT:
					log.Println("Received termination signal, shutting down...")
					Stop(mgr)
				}

			}
		}
	}()
}

func Start(mgr *GoroutineManager) {
	PCcfg, ProCcfg, Selfcfg = loadMonitorExporterConfig()
	initPrometheus(PCcfg, ProCcfg)
	//log.Println(PCcfg, ProCcfg, Selfcfg)
	log.Println("监控程序开始")
	if Selfcfg.MonitorKeepalive {
		log.Println("开启本程序保活")
		KeepGoProcessalive()
		time.Sleep(time.Second)
	}
	RestartServer(*Selfcfg)
	// 启动PC监控协程
	if PCcfg.CPU.Enabled {
		mgr.Start(func(ctx context.Context) {
			monitorCPU(ctx, PCcfg.CPU)
		})
	}
	if PCcfg.Memory.Enabled {
		mgr.Start(func(ctx context.Context) {
			monitorMemory(ctx, PCcfg.Memory)
		})
	}
	if PCcfg.Disk.Enabled {
		mgr.Start(func(ctx context.Context) {
			monitorDisk(ctx, PCcfg.Disk)
		})
	}
	if PCcfg.Network.Enabled {
		mgr.Start(func(ctx context.Context) {
			monitorNetwork(ctx, PCcfg.Network)
		})
	}

	// 启动进程监控协程
	if ProCcfg.Enabled {
		for _, proc := range ProCcfg.Process {
			mgr.Start(func(ctx context.Context) {
				monitorProc(ctx, proc)
			})
		}
	}

	// v1.2 进程监控优化

	//启动定时任务
	ScheduleCfg = loadScheduleTaskConfig()
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
				mgr.Start(func(ctx context.Context) {
					StartJob(ctx, job.Cron, *CreateShellJob(job))
				})

			}
		}
	}
}

func Restart(mgr *GoroutineManager) {
	mgr.Stop()
	mgr.Reset()
	UnRegister()
	log.Println("重新开始程序")
	Start(mgr)
}

func Stop(mgr *GoroutineManager) {
	mgr.Stop()
	StopServer()
	os.Exit(0)
}
