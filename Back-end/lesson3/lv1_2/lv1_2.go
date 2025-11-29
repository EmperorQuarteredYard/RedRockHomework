package lv2

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Data string `json:"data"`
}

func talkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("收到请求: %s %s\n", r.Method, r.URL.String())

	w.Header().Set("Content-Type", "application/json")
	msg := r.URL.Query().Get("msg")

	var response Response

	switch msg {
	case "ping":
		response.Data = "pong"
	case "helloserver":
		response.Data = "helloclient"
	case "甘雨是谁？":
		response.Data = "是你老婆啦"
	default:
		response.Data = "呜呜呜，Golang酱听不懂呢" + msg
	}

	fmt.Printf("返回响应: %s\n", response.Data)
	json.NewEncoder(w).Encode(response)
}

// 图片文件处理函数
func imageHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中提取文件名
	filename := "cat.jpg"

	fmt.Printf("请求图片: %s\n", filename)

	// 检查文件扩展名，确保是图片
	//ext := strings.ToLower(filepath.Ext(filename))
	//allowedExts := map[string]bool{
	//	".jpg":  true,
	//	".jpeg": true,
	//	".png":  true,
	//	".gif":  true,
	//	".bmp":  true,
	//	".webp": true,
	//}
	//
	//if !allowedExts[ext] {
	//	http.Error(w, "不支持的图片格式", http.StatusBadRequest)
	//	return
	//}

	// 尝试打开文件
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "图片不存在", http.StatusNotFound)
			return
		}
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 设置正确的Content-Type
	//switch ext {
	//case ".jpg", ".jpeg":
	//	w.Header().Set("Content-Type", "image/jpeg")
	//case ".png":
	//	w.Header().Set("Content-Type", "image/png")
	//case ".gif":
	//	w.Header().Set("Content-Type", "image/gif")
	//case ".bmp":
	//	w.Header().Set("Content-Type", "image/bmp")
	//case ".webp":
	//	w.Header().Set("Content-Type", "image/webp")
	//}

	// 将文件内容复制到响应中
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "无法发送图片", http.StatusInternalServerError)
		return
	}

	fmt.Printf("成功返回图片: %s\n", filename)
}

func Main() {
	// 注册/talk路由
	http.HandleFunc("/talk", talkHandler)

	// 注册图片文件路由 - 处理所有根路径下的文件请求
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 如果请求的是根路径，显示帮助信息
		if r.URL.Path == "/" {
			html := `
<!DOCTYPE html>
<html>
<head>
	<title>Go 服务器</title>
</head>
<body>
	<h1>服务器运行中</h1>
	<p>可用端点:</p>
	<ul>
		<li><a href="/talk?msg=ping">/talk?msg=ping</a> - 返回 {"data":"pong"}</li>
		<li><a href="/talk?msg=helloserver">/talk?msg=helloserver</a> - 返回 {"data":"helloclient"}</li>
		<li><a href="/cat.jpg">/cat.jpg</a> - 显示图片 (需要先将cat.jpg放在程序同目录下)</li>
	</ul>
</body>
</html>`
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, html)
			return
		}

		// 否则处理图片请求
		imageHandler(w, r)
	})

	port := ":8080"
	fmt.Printf("服务器启动在 http://localhost%s\n", port)
	fmt.Println("可用端点:")
	fmt.Println("  http://localhost:8080/talk?msg=ping")
	fmt.Println("  http://localhost:8080/talk?msg=helloserver")
	fmt.Println("  http://localhost:8080/cat.jpg")
	fmt.Println("\n注意: 请确保cat.jpg图片文件与程序在同一目录下")
	fmt.Println("按 Ctrl+C 停止服务器")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
