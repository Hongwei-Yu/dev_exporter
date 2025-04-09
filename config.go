package main

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
	"time"
)

type ResourceConfig struct {
	Enabled  bool           `mapstructure:"enabled"`
	Interval time.Duration  `mapstructure:"interval"`
	Params   map[string]any `mapstructure:",remain"` // 存储类型特定参数
}

type MonitorPcConfig struct {
	CPU     ResourceConfig `mapstructure:"cpu"`
	Memory  ResourceConfig `mapstructure:"memory"`
	Disk    ResourceConfig `mapstructure:"disk"`
	Network ResourceConfig `mapstructure:"network"`
}

type ProcConfig struct {
	Name           string        `yaml:"name"`
	Help           string        `yaml:"help"`
	KeepAlive      bool          `yaml:"keepAlive"`
	NotifyID       string        `yaml:"notifyID"`
	KeepaliveShell string        `yaml:"keepaliveShell"`
	KeepaliveWait  time.Duration `yaml:"keepaliveWait"`
	Interval       time.Duration `yaml:"interval"`
}

type MonitorProcConfig struct {
	Enabled bool `yaml:"enabled"`

	Process []ProcConfig `yaml:"process"`
}

type SelfConfig struct {
	MonitorKeepalive bool           `yaml:"monitorKeepalive"`
	Version          string         `yaml:"version"`
	Auth             string         `yaml:"auth"`
	Port             string         `yaml:"port"`
	Params           map[string]any `mapstructure:",remain"`
}

type ParameterConf struct {
	Parameter map[string]any `mapstructure:",remain"`
}

type JobConfig struct {
	Type       string        `mapstructure:"type"`
	Name       string        `mapstructure:"name"`
	Cron       string        `mapstructure:"cron"`
	Parameters ParameterConf `mapstructure:"parameters"`
}

type ScheduleTaskConfig struct {
	Enabled             bool        `yaml:"enabled" mapstructure:"enabled"`
	SkipIfStillRunning  bool        `yaml:"skipIfStillRunning" mapstructure:"skipIfStillRunning"`
	DelayIfStillRunning bool        `yaml:"delayIfStillRunning" mapstructure:"delayIfStillRunning"`
	Recover             bool        `yaml:"recover" mapstructure:"recover"`
	Job                 []JobConfig `yaml:"job" mapstructure:"job"`
}

func loadMonitorExporterConfig() (*MonitorPcConfig, *MonitorProcConfig, *SelfConfig) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var cfg MonitorPcConfig
	var ProCcfg MonitorProcConfig
	var SelfCfg SelfConfig
	err := viper.UnmarshalKey("MonitorPC", &cfg, func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.TextUnmarshallerHookFunc(),
		)
	})
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
	if err := viper.UnmarshalKey("MonitorProc", &ProCcfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
	if err := viper.UnmarshalKey("Self", &SelfCfg); err != nil {
		log.Fatalf("Self unmarshaling config: %v", err)
	}
	return &cfg, &ProCcfg, &SelfCfg
}

func loadScheduleTaskConfig() *ScheduleTaskConfig {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	var cfg ScheduleTaskConfig
	err := viper.UnmarshalKey("ScheduleTask", &cfg, func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
		)
	})
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
	return &cfg
}
