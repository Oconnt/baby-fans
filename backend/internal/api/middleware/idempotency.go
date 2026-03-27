package middleware

import (
	"bytes"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IdempotencyMiddleware prevents duplicate request processing using X-Request-ID header
// Thread-safe in-memory storage with automatic expiration

type responseCapture struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseCapture) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

type requestRecord struct {
	processedAt time.Time
	statusCode  int
	body        []byte
}

var (
	idempotencyStore = sync.Map{}
	cleanupInterval  = 10 * time.Minute
	requestTTL       = 24 * time.Hour
)

func init() {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			cleanupExpiredRecords()
		}
	}()
}

func cleanupExpiredRecords() {
	now := time.Now()
	idempotencyStore.Range(func(key, value interface{}) bool {
		record := value.(*requestRecord)
		if now.Sub(record.processedAt) > requestTTL {
			idempotencyStore.Delete(key)
		}
		return true
	})
}

func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Header("X-Request-ID", requestID)
		}

		// Check if already processed
		if existing, ok := idempotencyStore.Load(requestID); ok {
			record := existing.(*requestRecord)
			c.Header("X-Idempotency-Replayed", "true")
			c.Data(record.statusCode, "application/json", record.body)
			c.Abort()
			return
		}

		// Capture response
		capture := &responseCapture{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = capture

		c.Next()

		// Only cache successful responses
		if capture.statusCode >= 200 && capture.statusCode < 300 {
			idempotencyStore.Store(requestID, &requestRecord{
				processedAt: time.Now(),
				statusCode:  capture.statusCode,
				body:        capture.body.Bytes(),
			})
		}
	}
}

func (r *responseCapture) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
