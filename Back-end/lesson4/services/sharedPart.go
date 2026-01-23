package services

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	htmlNotFound []byte = []byte("<!DOCTYPE html>\n\t\t\t\t<html>\n\t\t\t\t<head><title>错误</title></head>\n\t\t\t\t<body>\n\t\t\t\t\t<h1>找不到index.html文件</h1>\n\t\t\t\t\t<p>请确保index.html文件与程序在同一目录下</p>\n\t\t\t\t</body>\n\t\t\t\t</html>")
)

func loadHtml(filename string) []byte {
	htmlContent, err := os.ReadFile("static/" + filename + ".html")
	if err != nil {
		fmt.Println(err)
		return htmlNotFound
	}
	return htmlContent
}

func registerCRUDRoutes(url string, get func(c *gin.Context), post func(c *gin.Context), create func(c *gin.Context) error, read func(c *gin.Context) error, update func(c *gin.Context), delete func(c *gin.Context)) {

}
