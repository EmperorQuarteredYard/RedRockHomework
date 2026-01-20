package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// 请求数据结构
type GradeRequest struct {
	Name  string    `json:"name"`
	Score []float64 `json:"score"`
}

// 响应数据结构
type GradeResponse struct {
	Average float64 `json:"average"`
}

func main() {
	router := gin.Default()

	// 首页 - 显示输入表单
	router.GET("/", func(c *gin.Context) {

		htmlContent, err := os.ReadFile("lv2\\index.html")

		if err != nil {
			// 如果文件不存在，返回一个简单的错误页面
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
				<!DOCTYPE html>
				<html>
				<head><title>错误</title></head>
				<body>
					<h1>找不到index.html文件</h1>
					<p>请确保index.html文件与程序在同一目录下</p>
				</body>
				</html>
			`))
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})

	router.Static("/static", "./lv2/static")

	router.POST("/calculate", func(c *gin.Context) {
		fmt.Println(0)
		// 解析请求体中的 JSON数据
		var req GradeRequest

		// 读取原始请求体
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
			return
		}
		fmt.Println(string(body))
		// 解析 JSON
		if err := json.Unmarshal(body, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON解析错误: " + err.Error()})
			return
		}

		// 验证数据
		fmt.Println(1)
		if req.Name == "" || len(req.Score) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "姓名和成绩数组不能为空"})
			return
		}

		// 计算平均分
		fmt.Println(2)
		var sum float64
		for _, score := range req.Score {
			sum += score
		}
		average := sum / float64(len(req.Score))

		// 返回JSON响应
		fmt.Println(average)
		resp := GradeResponse{
			Average: average,
		}

		c.JSON(http.StatusOK, resp)
	})
	// 启动服务器
	router.Run(":8080")
}
