package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"time"
)

// CPU监控
func monitorCPU(cfg ResourceConfig) {
	gid := GetGID()
	log.Println(gid, "CPU监控开始")
	perCore, _ := cfg.Params["per_core"].(bool)

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
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

// 内存监控
func monitorMemory(cfg ResourceConfig) {
	gid := GetGID()
	log.Println(gid, "内存监控开始")
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
		memory, _ := mem.VirtualMemory()
		//fmt.Println(mem)
		memUsage.Set(memory.UsedPercent)
		//memUsage.Set(1)
	}
}

// 磁盘监控
func monitorDisk(cfg ResourceConfig) {
	gid := GetGID()
	log.Println(gid, "磁盘监控开始")
	mounts, _ := cfg.Params["mountpoints"].([]interface{})

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
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

// 网络监控
func monitorNetwork(cfg ResourceConfig) {
	gid := GetGID()
	log.Println(gid, "网卡监控开始")
	ifaces, _ := cfg.Params["interfaces"].([]interface{})

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	prevStats := make(map[string]net.IOCountersStat)

	for range ticker.C {
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

func monitorProc(cfg ProcConfig) {
	gid := GetGID()
	log.Println(gid, "进程"+cfg.Name+"开始监控")
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	for range ticker.C {
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

		}
	}

}
