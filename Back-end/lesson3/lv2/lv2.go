package lv2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GradeRequest struct {
	Name  string    `json:"name"`
	Score []float64 `json:"score"`
}

type GradeResponse struct {
	Average float64 `json:"average"`
}

func calculateGrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println("传入", r.Body)
	w.Header().Set("Content-Type", "application/json")
	var request GradeRequest
	decoder := json.NewDecoder(r.Body)
	fmt.Println((decoder))
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, `{"error":"无效的JSON数据"}`, http.StatusBadRequest)
		return
	}

	if len(request.Score) == 0 {
		http.Error(w, `{"error":"分数数组不能为空"}`, http.StatusBadRequest)
		return
	}
	var result GradeResponse
	result.Average = 0
	for i := 0; i < len(request.Score); i++ {
		result.Average += request.Score[i]
	}
	result.Average /= float64(len(request.Score))
	json.NewEncoder(w).Encode(result)
	fmt.Printf("收到请求 - 学生: %s, 分数: %v, 平均分: %.2f\n",
		request.Name, request.Score, result)
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `
<!DOCTYPE html>
<html>
<head>
	<title>学生成绩计算器</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 40px; }
		.container { max-width: 600px; margin: 0 auto; }
		.code { background: #f5f5f5; padding: 10px; border-radius: 5px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>学生成绩计算器</h1>
		
		<h2>使用方法</h2>
		<p>发送POST请求到 <code>/calculate-grade</code> 端点</p>
		
		<h3>请求示例:</h3>
		<div class="code">
			<pre>{
  "name": "bob",
  "score": [68, 97.4, 94.2, 75.4]
}</pre>
		</div>
		
		<h3>响应示例:</h3>
		<div class="code">
			<pre>{
  "average": 83.75
}</pre>
		</div>
		
		<h3>使用curl测试:</h3>
		<div class="code">
			<pre>curl -X POST http://localhost:8080/calculate-grade \
  -H "Content-Type: application/json" \
  -d '{"name":"bob","score":[68,97.4,94.2,75.4]}'</pre>
		</div>
	</div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func main() {
	// 注册路由
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calculate-grade", calculateGrade)

	port := ":8080"
	fmt.Printf("服务器启动在 http://localhost%s\n", port)
	fmt.Println("可用端点:")
	fmt.Println("  GET  http://localhost:8080/ - 显示使用说明")
	fmt.Println("  POST http://localhost:8080/calculate-grade - 计算平均成绩")

	log.Fatal(http.ListenAndServe(port, nil))
}
func Main() {
	main()
}

//样例输入：powershell版本
//Invoke-RestMethod -Uri "http://localhost:8080/calculate-grade" -Method Post -ContentType "application/json" -Body '{"name":"bob","score":[68,97.4,94.2,75.4]}'
