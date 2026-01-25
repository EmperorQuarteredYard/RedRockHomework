package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// bodyLogWriter 用于记录响应体的 writer
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 日志中间件
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 记录请求体（对于非二进制请求）
		var requestBody []byte
		if c.Request.Body != nil && c.ContentType() != "multipart/form-data" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录响应体
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 状态码
		statusCode := c.Writer.Status()

		// 客户端IP
		clientIP := c.ClientIP()

		// 请求方法
		reqMethod := c.Request.Method

		// 请求路径
		reqPath := c.Request.URL.Path

		// 查询参数
		query := c.Request.URL.RawQuery
		if query != "" {
			reqPath = reqPath + "?" + query
		}

		// 响应体
		responseBody := blw.body.String()

		// 根据状态码选择日志级别
		logFields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", reqMethod),
			zap.String("path", reqPath),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// 记录错误信息
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error("请求错误", zap.String("error", e))
			}
		} else if statusCode >= 400 {
			// 只记录简化的请求/响应体信息
			if len(requestBody) > 0 && len(requestBody) < 1024 {
				logFields = append(logFields, zap.String("request_body", string(requestBody)))
			}
			if len(responseBody) > 0 && len(responseBody) < 1024 {
				logFields = append(logFields, zap.String("response_body", responseBody))
			}

			if statusCode >= 500 {
				logger.Error("服务器错误", logFields...)
			} else {
				logger.Warn("客户端错误", logFields...)
			}
		} else {
			logger.Info("请求处理完成", logFields...)
		}
	}
}
