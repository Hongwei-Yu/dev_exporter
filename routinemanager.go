package main

import (
	"context"
	"log"
	"sync"
)

// GoroutineManager 协程管理器
type GoroutineManager struct {
	wg   sync.WaitGroup
	ctx  context.Context
	stop context.CancelFunc
}

// NewGoroutineManager 创建管理器
func NewGoroutineManager(parentCtx context.Context) *GoroutineManager {
	ctx, stop := context.WithCancel(parentCtx)
	return &GoroutineManager{
		ctx:  ctx,
		stop: stop,
	}
}

// Start 启动协程（支持任意函数签名）
func (gm *GoroutineManager) Start(fn func(context.Context)) {
	gm.wg.Add(1)
	go func() {
		defer gm.wg.Done()
		fn(gm.ctx)
	}()
}

// Stop 触发停止所有协程
func (gm *GoroutineManager) Stop() {
	log.Println("协程停止")
	gm.stop()    // 发送停止信号
	gm.wg.Wait() // 等待所有协程退出
	log.Println("全部协程退出")
}
