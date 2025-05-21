package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type SseHandler struct {
	clients map[chan string]bool
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSseHandler() *SseHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &SseHandler{
		clients: make(map[chan string]bool),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Serve 处理 SSE 连接
func (h *SseHandler) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientChan := make(chan string, 10) // 使用缓冲通道避免阻塞

	h.mu.Lock()
	h.clients[clientChan] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, clientChan)
		h.mu.Unlock()
		close(clientChan)
	}()

	// 使用Flush中间件确保正确刷新
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case msg := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		case <-h.ctx.Done(): // 服务器关闭时退出
			return
		}
	}
}

// SimulateEvents 模拟事件并安全发送
func (h *SseHandler) SimulateEvents() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			message := fmt.Sprintf("Server time: %s", time.Now().Format(time.RFC3339))

			h.mu.RLock()
			clients := make([]chan string, 0, len(h.clients))
			for ch := range h.clients {
				clients = append(clients, ch)
			}
			h.mu.RUnlock()

			for _, ch := range clients {
				// 安全发送，避免panic
				func() {
					defer func() {
						if r := recover(); r != nil {
							logx.Errorf("send error: %v", r)
						}
					}()
					select {
					case ch <- message:
					default:
						// 跳过阻塞的通道
					}
				}()
			}
		case <-h.ctx.Done():
			return
		}
	}
}

func main() {
	server := rest.MustNewServer(rest.RestConf{
		Host: "0.0.0.0",
		Port: 8080,
	}, rest.WithFileServer("/static", http.Dir("./go_demo/sse_demo/zero/static"))) // 确认静态文件路径

	sseHandler := NewSseHandler()
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/r3labs_sse",
		Handler: sseHandler.Serve,
	}, rest.WithTimeout(0))

	go sseHandler.SimulateEvents()

	// 优雅关闭处理
	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		// signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) // 按需引入os/signal
		<-c
		sseHandler.cancel()
		server.Stop()
		close(done)
	}()

	logx.Info("Server starting on :8080")
	server.Start()
	<-done
}
