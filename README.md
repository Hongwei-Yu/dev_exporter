# 运维小工具v1.0
## 功能简介
### 监控
#### pc监控
```yaml
MonitorPC:
  disk:
    enabled: false
    mountpoints:
      - "E:"
    interval: 10s

  cpu:
    enabled: false        # 是否启用CPU监控
    per_core: true        # 是否监控每个核心
    interval: 5s         # 采集间隔

  memory:
    enabled: false
    interval: 5s

  network:
    enabled: false
    interfaces: [ "enp4s0" ]  # 指定监控的网卡
    interval: 10s
```
目前实现了cpu使用率监控、内存使用率监控、磁盘使用率监控、网卡监控

#### 进程监控
```yaml
MonitorProc:
  enabled: false
  process:
    - name: xxx
      interval: 10s
#      help: "LZStation online (1 online 0 offline)"
      keepalive: 1 # offline后是否拉起
      keepaliveShell: "cd /home/bin/ && ./startApp.sh" # 拉起命令
      keepaliveWait: 60s # 拉起命令执行后等待时间
```
进程监控主要监控进程是否在线，在线状态 1 掉线状态 0
同时添加了进程掉线拉起功能，此部分需要将keepalive设置为1，keepaliveShell中填写shell拉起指令
，为防止两个相邻监控周期重复拉起的情况需要设置每次shell执行后等待时间

#### 定时任务
```yaml
# 定时任务 shell 类型
ScheduleTask:
  enabled: true
#  skipIfStillRunning: true # 如果上次任务还正在运行，那么跳过本次任务的运行并记录日记
#  delayIfStillRunning: true # 如果上次任务还正在运行，那么延迟执行本次任务的运行并记录日记
#  recover: true # 如果 Job Panic，记录日记
  job:
    - name: "清理任务"
      cron: "00 23 * * *"
      type: "shell"
      parameters:
        shellcommand: "cd /home/ftp&&find ./ -mtime +1 -exec rm -rf {} \\;"
```
此功能实现shell定时任务，可以依据标准cron设置执行时间，shellcommand中填写执行的任务，复杂的shell任务可以单独防止脚本中，再通过该配置进行调用

#### 保活功能
为防止系统内存占用高出现的oom_kill情况，设置了程序保活机制将oom_score_adj设置为-17
```yaml
# 监控程序是否保活方式被系统oom killer  true表示保活 false表示不保活
Self:
  MonitorKeepalive: false
  Version: "v1.0"
  Port: ":8080"
```

# 运维小工具v1.1
## 准备实现功能
1. 配置文件热更新


