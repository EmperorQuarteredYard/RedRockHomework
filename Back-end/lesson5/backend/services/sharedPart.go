package services

import (
	"fmt"
	"os"
)

var (
	htmlNotFound []byte = []byte("<!DOCTYPE html>\n<html>\n<head><title>错误</title></head>\n<body>\n\t<h1>找不到html文件</h1>\n</body>\n</html>")
)

func loadHtml(filename string) []byte {
	htmlContent, err := os.ReadFile("frontend/" + filename + ".html")
	if err != nil {
		fmt.Println("Failed to load HTML:", err)
		return htmlNotFound
	}
	return htmlContent
}
