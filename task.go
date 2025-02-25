package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"os/exec"
)

type ShellJob struct {
	Name    string
	Command string
	Shut    chan int
}

func CreateShellJob(job JobConfig) *ShellJob {
	return &ShellJob{
		Name:    job.Name,
		Command: job.Parameters.Parameter["shellcommand"].(string),
	}
}

func (j *ShellJob) Run() {

	log.Println("开始执行" + j.Name)
	err := exec.Command("bash", "-c", j.Command).Run()
	if err != nil {
		log.Println(j.Name + " 执行失败")
	}
}

func StartJob(spec string, job ShellJob) {

	c := cron.New()

	c.AddJob(spec, &job)

	// 启动执行任务
	c.Start()
	// 退出时关闭计划任务
	defer c.Stop()

	// 如果使用 select{} 那么就一直会循环
	select {
	case <-job.Shut:
		return
	}
}
func StopJob(shut chan int) {
	shut <- 0
}
