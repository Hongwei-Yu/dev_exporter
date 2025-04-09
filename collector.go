package main

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"sync/atomic"
	"time"
)

// CPU监控
func monitorCPU(ctx context.Context, cfg ResourceConfig) {
	// 初始化原子变量（标记运行状态）
	var isRunning int32 = 1
	defer atomic.StoreInt32(&isRunning, 0) // 退出时标记为停止

	gid := GetGID()
	log.Println(gid, "CPU监控开始")
	defer log.Println(gid, "CPU监控停止")

	perCore, _ := cfg.Params["per_core"].(bool)
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		// 优先检查退出信号和运行状态
		if atomic.LoadInt32(&isRunning) == 0 {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 执行采集前再次检查运行状态
			if atomic.LoadInt32(&isRunning) == 0 {
				return
			}

			// 核心采集逻辑
			if perCore {
				percents, _ := cpu.Percent(time.Second, true)
				for core, pct := range percents {
					cpuUsage.WithLabelValues(fmt.Sprintf("%d", core)).Set(pct)
				}
			} else {
				pct, _ := cpu.Percent(time.Second, false)
				cpuUsage.WithLabelValues("all").Set(pct[0])
			}
		}
	}
}

// 内存监控
func monitorMemory(ctx context.Context, cfg ResourceConfig) {
	var isRunning int32 = 1
	defer atomic.StoreInt32(&isRunning, 0) // 退出时标记为停止
	gid := GetGID()
	log.Println(gid, "内存监控开始")
	defer log.Println(gid, "内存监控停止")
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	for {
		// 优先检查退出信号和运行状态
		if atomic.LoadInt32(&isRunning) == 0 {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 执行采集前再次检查运行状态
			if atomic.LoadInt32(&isRunning) == 0 {
				return
			}

			// 核心采集逻辑
			memory, _ := mem.VirtualMemory()
			//fmt.Println(mem)
			memUsage.Set(memory.UsedPercent)

		}
	}

}

// 磁盘监控
func monitorDisk(ctx context.Context, cfg ResourceConfig) {
	var isRunning int32 = 1
	defer atomic.StoreInt32(&isRunning, 0) // 退出时标记为停止
	gid := GetGID()
	log.Println(gid, "磁盘监控开始")
	defer log.Println(gid, "磁盘监控停止")
	mounts, _ := cfg.Params["mountpoints"].([]interface{})

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		// 优先检查退出信号和运行状态
		if atomic.LoadInt32(&isRunning) == 0 {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if atomic.LoadInt32(&isRunning) == 0 {
				return
			}
			// 执行采集前再次检查运行状态
			for _, m := range mounts {
				usage, err := disk.Usage(m.(string))
				if err != nil {
					log.Fatal(err)
				}
				diskUsage.WithLabelValues(m.(string)).Set(usage.UsedPercent)
				//log.Println(m.(string))
			}

		}
	}

}

// 网络监控
func monitorNetwork(ctx context.Context, cfg ResourceConfig) {
	var isRunning int32 = 1
	defer atomic.StoreInt32(&isRunning, 0) // 退出时标记为停止
	gid := GetGID()
	log.Println(gid, "网卡监控开始")
	defer log.Println(gid, "网卡监控停止")
	ifaces, _ := cfg.Params["interfaces"].([]interface{})

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	prevStats := make(map[string]net.IOCountersStat)

	for {
		// 优先检查退出信号和运行状态
		if atomic.LoadInt32(&isRunning) == 0 {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if atomic.LoadInt32(&isRunning) == 0 {
				return
			}
			// 执行采集前再次检查运行状态
			stats, _ := net.IOCounters(true)
			for _, s := range stats {
				for _, target := range ifaces {
					if s.Name == target.(string) {
						// 计算差值速率
						if prev, ok := prevStats[s.Name]; ok {
							delta := s.BytesSent - prev.BytesSent
							networkTx.WithLabelValues(s.Name).Set(float64(delta))
						}
						prevStats[s.Name] = s
					}
				}
			}

		}
	}

}

func monitorProc(ctx context.Context, cfg ProcConfig) {
	var isRunning int32 = 1
	defer atomic.StoreInt32(&isRunning, 0) // 退出时标记为停止
	gid := GetGID()
	log.Println(gid, "进程"+cfg.Name+"开始监控")
	defer log.Println(gid, "进程"+cfg.Name+"监控停止")
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	for {
		// 优先检查退出信号和运行状态
		if atomic.LoadInt32(&isRunning) == 0 {
			return
		}

		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if atomic.LoadInt32(&isRunning) == 0 {
				return
			}
			online := CheckProcess(cfg.Name)
			if online {
				procMonitor.WithLabelValues(cfg.Name).Set(1)
			} else {
				procMonitor.WithLabelValues(cfg.Name).Set(0)
				log.Println(gid, cfg.Name+" down")
				if cfg.KeepAlive {
					log.Println(gid, "开始拉起: "+cfg.Name)
					HandUpProc(&cfg)

					log.Println(gid, "本协程程开始沉睡 "+cfg.KeepaliveWait.String())
					time.Sleep(cfg.KeepaliveWait)
				}
				// 执行采集前再次检查运行状态

			}
		}

	}

}

// v1.2版本优化进程监控

func monitorProc_CPU(ctx context.Context, cfg ProcConfig) {
	var isRunning int32 = 1
	defer atomic.StoreInt32(&isRunning, 0) // 退出时标记为停止
	gid := GetGID()
	log.Println(gid, "进程"+cfg.Name+"开始CPU使用率监控")
	defer log.Println(gid, "进程"+cfg.Name+"CPU使用率监控停止")
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	for {
		// 优先检查退出信号和运行状态
		if atomic.LoadInt32(&isRunning) == 0 {
			return
		}

		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if atomic.LoadInt32(&isRunning) == 0 {
				return
			}
			pid, err := GetProcessPidByName(cfg.Name)

			if err != nil {
				log.Println(gid, "进程"+cfg.Name+"cpu使用率监控错误: ", err)

			} else {
				percent := ProcCpuMonitor(*pid)
				if percent != nil {
					procMonitorCpu.WithLabelValues(cfg.Name).Set(*percent)
				} else {
					log.Println(gid, "进程"+cfg.Name+"cpu使用率监控错误: ", err)
				}

			}
		}

	}

}

func monitorProc_MEM(ctx context.Context, cfg ProcConfig) {
	var isRunning int32 = 1
	defer atomic.StoreInt32(&isRunning, 0) // 退出时标记为停止
	gid := GetGID()
	log.Println(gid, "进程"+cfg.Name+"开始CPU使用率监控")
	defer log.Println(gid, "进程"+cfg.Name+"CPU使用率监控停止")
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	for {
		// 优先检查退出信号和运行状态
		if atomic.LoadInt32(&isRunning) == 0 {
			return
		}

		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if atomic.LoadInt32(&isRunning) == 0 {
				return
			}
			pid, err := GetProcessPidByName(cfg.Name)

			if err != nil {
				log.Println(gid, "进程"+cfg.Name+"cpu使用率监控错误: ", err)

			} else {
				percent := ProcMemMonitor(*pid)
				if percent != nil {
					procMonitorMem.WithLabelValues(cfg.Name).Set(float64(*percent))
				} else {
					log.Println(gid, "进程"+cfg.Name+"cpu使用率监控错误: ", err)
				}

			}
		}

	}

}
