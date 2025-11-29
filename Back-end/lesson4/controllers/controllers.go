package controllers

import (
	"RedrockHomework/lesson4/models"
	"RedrockHomework/lesson4/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StudentController struct {
	studentService services.StudentService
}

// 这里传入的应该是一个Student类型的JSON
func (sc StudentController) CreateStudent(c *gin.Context) {
	var student models.STUDENT
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := sc.studentService.CreateStudent(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"student": student})
}

func (sc StudentController) SelectLesson(c *gin.Context) {
	studentID, _ := strconv.Atoi(c.Param("studentId"))
	if err := sc.studentService.SelectLesson(studentID, c.Param("lessonCode")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "选课成功"})
}

func (sc StudentController) GetStudentLessons(c *gin.Context) {
	studentID, _ := strconv.Atoi(c.Param("studentId"))
	lessons, err := sc.studentService.GetStudentLessons(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": lessons})
}
