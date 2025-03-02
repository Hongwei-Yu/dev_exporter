package main

import "github.com/prometheus/client_golang/prometheus"

var (
	// CPU指标
	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_cpu_usage_percent",
			Help: "CPU usage percentage",
		},
		[]string{"core"}, // 标签：核心编号
	)

	// 内存指标
	memUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "node_memory_usage_percent",
			Help: "Memory usage percentage",
		},
	)

	// 磁盘指标
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_disk_usage_percent",
			Help: "Disk usage percentage",
		},
		[]string{"mountpoint"}, // 标签：挂载点
	)

	// 网络指标
	networkTx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_network_transmit_bytes",
			Help: "Network transmit bytes",
		},
		[]string{"interface"},
	)

	procMonitor = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_monitor_online",
			Help: "Monitor process online (1 online 0 offline)",
		},
		[]string{"Name"},
	)
)

func initPrometheus(PCcfg *MonitorPcConfig, ProCcfg *MonitorProcConfig) {
	if PCcfg.CPU.Enabled {
		prometheus.MustRegister(cpuUsage)
	}
	if PCcfg.Memory.Enabled {
		prometheus.MustRegister(memUsage)
	}
	if PCcfg.Disk.Enabled {
		prometheus.MustRegister(diskUsage)
	}
	if PCcfg.Network.Enabled {
		prometheus.MustRegister(networkTx)
	}
	if ProCcfg.Enabled {
		prometheus.MustRegister(procMonitor)
	}
}

func UnRegister() {
	if PCcfg.CPU.Enabled {
		prometheus.Unregister(cpuUsage)
	}
	if PCcfg.Memory.Enabled {
		prometheus.Unregister(memUsage)
	}
	if PCcfg.Disk.Enabled {
		prometheus.Unregister(diskUsage)
	}
	if PCcfg.Network.Enabled {
		prometheus.Unregister(networkTx)
	}
	if ProCcfg.Enabled {
		prometheus.Unregister(procMonitor)
	}
}
