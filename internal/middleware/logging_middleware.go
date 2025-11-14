package middleware

import (
	"context"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CreateLoggingMiddleware MCP 中间件，用于记录方法调用日志。
func CreateLoggingMiddleware() mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			start := time.Now()
			sessionID := req.GetSession().ID()
			log.Printf("[请求] Session: %s | Method: %s", sessionID, method)

			result, err := next(ctx, method, req)

			duration := time.Since(start)

			if err != nil {
				log.Printf("[响应] Session: %s | Method: %s | Status: 失败 | Duration: %v | Error: %v", sessionID, method, duration, err)
			} else {
				log.Printf("[响应] Session: %s | Method: %s | Status: 成功 | Duration: %v", sessionID, method, duration)
			}

			return result, err
		}
	}
}
