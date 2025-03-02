package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var server *http.Server

func StartServer(config SelfConfig) {
	server = &http.Server{
		Addr: fmt.Sprintf(":%s", config.Port),
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
}

func RestartServer(config SelfConfig) {
	oldServer := server
	newServer := &http.Server{
		Addr: fmt.Sprintf(":%s", config.Port),
	}

	// 启动新服务器
	go func() {
		log.Printf("Starting new server on %s", newServer.Addr)
		if err := newServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 优雅关闭旧服务器
	if oldServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		if err := oldServer.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown old server: %v", err)
		}
	}

	server = newServer
}

func StopServer() {
	oldServer := server
	// 优雅关闭旧服务器
	if oldServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := oldServer.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown old server: %v", err)
		} else {
			log.Printf("Shutdown Server sucess")
		}
	}

}
