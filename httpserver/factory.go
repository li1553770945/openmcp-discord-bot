package httpserver

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/li1553770945/openmcp-discord-bot/httpserver/handler"
	"net/http"
	"sync"
	"time"
)

func StartHttpServer(ctx context.Context, wg *sync.WaitGroup) {
	// 1. 显式创建 http.Server 实例（才能控制 Shutdown）
	server := &http.Server{
		Addr:    ":80",
		Handler: nil,
	}
	// 2. 注册路由
	http.HandleFunc("/api/github-webhook/release", handler.GithubRelease)
	http.HandleFunc("/api/message", handler.SendMessageHandler)
	// 3. 监听 ctx 取消信号并触发关闭
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 启动一个子协程运行服务（因为 ListenAndServe 会阻塞）
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("HTTP 服务启动失败: %v", err))
			}
		}()
		// 阻塞等待退出信号
		<-ctx.Done()

		// 优雅关闭服务器（设置 5 秒超时）
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("HTTP 服务关闭错误: %v", err)
		}
	}()
}
