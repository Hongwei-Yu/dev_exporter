package main

import (
	"context"
	"log"
)

var PCcfg *MonitorPcConfig

var ProCcfg *MonitorProcConfig

var Selfcfg *SelfConfig

var ScheduleCfg *ScheduleTaskConfig

func main() {
	rootCtx := context.Background()
	mgr := NewGoroutineManager(rootCtx)
	PCcfg, ProCcfg, Selfcfg = loadMonitorExporterConfig()
	log.Println(Selfcfg.Version, "write by YHW")
	Start(mgr)
	setupSignalHandler(mgr)

	select {} // 阻塞主线程

}
