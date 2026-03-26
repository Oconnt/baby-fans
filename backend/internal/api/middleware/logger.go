package middleware

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// requestLog 用于 JSON 格式的请求日志
type requestLog struct {
	Timestamp  string `json:"timestamp"`
	Level      string `json:"level"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	Latency    string `json:"latency"`
	ClientIP   string `json:"client_ip"`
	Error      string `json:"error,omitempty"`
	StackTrace string `json:"stack_trace,omitempty"`
}

// jsonLogger 输出 JSON 格式日志到 stdout（Docker 可收集）
func JSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 请求结束后记录日志
		latency := time.Since(start)

		entry := requestLog{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "INFO",
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			Latency:   latency.String(),
			ClientIP:  c.ClientIP(),
		}

		// 如果有错误信息，从 context 中获取
		if len(c.Errors) > 0 {
			entry.Error = c.Errors.String()
			entry.Level = "ERROR"
		}

		// 如果是 500 错误，降低级别
		if c.Writer.Status() >= 500 {
			entry.Level = "ERROR"
		} else if c.Writer.Status() >= 400 {
			entry.Level = "WARN"
		}

		logBytes, _ := json.Marshal(entry)
		fmt.Println(string(logBytes))
	}
}

// RecoveryWithLogger 带有日志输出的 Recovery 中间件
func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())

				entry := requestLog{
					Timestamp:  time.Now().Format(time.RFC3339),
					Level:      "PANIC",
					Method:     c.Request.Method,
					Path:       c.Request.URL.Path,
					Status:     500,
					ClientIP:   c.ClientIP(),
					Error:      fmt.Sprintf("%v", err),
					StackTrace: stack,
				}
				logBytes, _ := json.Marshal(entry)
				fmt.Println(string(logBytes))

				c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			}
		}()
		c.Next()
	}
}
