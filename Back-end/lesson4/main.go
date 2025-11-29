package main

import (
	"RedrockHomework/lesson4/controllers"
	"RedrockHomework/lesson4/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	r := gin.Default()
	studentController := controllers.StudentController{}

	studentRoutes := r.Group("/students")
	{
		studentRoutes.POST("/", studentController.CreateStudent)
		studentRoutes.POST("/:studentId/lessons/:lessonId", studentController.SelectLesson)
		studentRoutes.GET("/:studentId/lessons", studentController.GetStudentLessons)
	}
	r.Run(":8080")
}
