package controller

import (
	"homeworkSystem/backend/internal/models"
	"homeworkSystem/backend/internal/service"
	"homeworkSystem/backend/pkg/errcode"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubmissionController struct {
	BaseController
	svc *service.Service
}

func NewSubmissionController(svc *service.Service) *SubmissionController {
	return &SubmissionController{svc: svc}
}

// Submit 提交作业（小登）
func (ctl *SubmissionController) Submit(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.Role != "student" {
		ctl.HandleCode(c, errcode.ErrInsufficientPermissions)
		return
	}

	var req struct {
		HomeworkID uint64 `json:"homework_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
		FileURL    string `json:"file_url"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	input := service.SubmitInput{
		HomeworkID: req.HomeworkID,
		StudentID:  authUser.UserID,
		Content:    req.Content,
		FileURL:    req.FileURL,
	}
	submission, err := ctl.svc.Submit(input)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"id":           submission.ID,
		"homework_id":  submission.HomeworkID,
		"is_late":      submission.IsLate,
		"submitted_at": submission.SubmittedAt.Format("2006-01-02 15:04:05"),
	})
}

// MySubmissions 我的提交列表（小登）
func (ctl *SubmissionController) MySubmissions(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.Role != "student" {
		ctl.HandleCode(c, errcode.ErrInsufficientPermissions)
		return
	}

	list, err := ctl.svc.GetMySubmissions(authUser.UserID)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	var respList []gin.H
	for _, s := range list {
		homework, _ := ctl.svc.GetAssignmentDetail(s.HomeworkID) // 忽略错误
		respList = append(respList, gin.H{
			"id": s.ID,
			"homework": gin.H{
				"id":               homework.ID,
				"title":            homework.Title,
				"department":       homework.Department,
				"department_label": ctl.DepartmentLabel(homework.Department),
			},
			"score":        s.Score,
			"comment":      s.Comment,
			"is_excellent": s.IsExcellent,
			"submitted_at": s.SubmittedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctl.Success(c, gin.H{
		"list":  respList,
		"total": len(respList),
	})
}

// HomeworkSubmissions 获取作业的所有提交（老登）
func (ctl *SubmissionController) HomeworkSubmissions(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.Role != "admin" {
		ctl.HandleCode(c, errcode.ErrInsufficientPermissions)
		return
	}

	homeworkID, err := strconv.ParseUint(c.Param("homework_id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.ErrInvalidParams)
		return
	}

	// 检查作业部门是否与老登一致
	homework, err := ctl.svc.GetAssignmentDetail(homeworkID)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}
	if homework.Department != authUser.Department {
		ctl.HandleCode(c, errcode.ErrDepartmentNotMatch)
		return
	}

	list, err := ctl.svc.GetHomeworkSubmissions(homeworkID)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	var respList []gin.H
	for _, s := range list {
		student, _ := ctl.svc.GetUserByID(s.StudentID)
		respList = append(respList, gin.H{
			"id": s.ID,
			"student": gin.H{
				"id":               student.ID,
				"nickname":         student.Nickname,
				"department":       student.Department,
				"department_label": ctl.DepartmentLabel(student.Department),
			},
			"content":      s.Content,
			"file_url":     s.FileURL,
			"is_late":      s.IsLate,
			"score":        s.Score,
			"comment":      s.Comment,
			"submitted_at": s.SubmittedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctl.Success(c, gin.H{
		"list":  respList,
		"total": len(respList),
	})
}

// Review 批改作业（老登）
func (ctl *SubmissionController) Review(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.Role != "admin" {
		ctl.HandleCode(c, errcode.ErrInsufficientPermissions)
		return
	}

	submissionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.ErrInvalidParams)
		return
	}

	var req struct {
		Score       int    `json:"score"`
		Comment     string `json:"comment" binding:"required"`
		IsExcellent bool   `json:"is_excellent"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	reviewer, err := ctl.svc.GetUserByID(authUser.UserID)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	submission, err := ctl.svc.ReviewSubmission(reviewer, submissionID, req.Score, req.Comment, req.IsExcellent)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"id":           submission.ID,
		"score":        submission.Score,
		"comment":      submission.Comment,
		"is_excellent": submission.IsExcellent,
		"reviewed_at":  submission.ReviewedAt.Format("2006-01-02 15:04:05"),
	})
}

// MarkExcellent 标记/取消优秀作业（老登）
func (ctl *SubmissionController) MarkExcellent(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.Role != "admin" {
		ctl.HandleCode(c, errcode.ErrInsufficientPermissions)
		return
	}

	submissionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.ErrInvalidParams)
		return
	}

	var req struct {
		IsExcellent bool `json:"is_excellent" binding:"required"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	reviewer, err := ctl.svc.GetUserByID(authUser.UserID)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	submission, err := ctl.svc.MarkExcellent(reviewer, submissionID, req.IsExcellent)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"id":           submission.ID,
		"is_excellent": submission.IsExcellent,
	})
}

// ExcellentList 获取优秀作业列表
func (ctl *SubmissionController) ExcellentList(c *gin.Context) {
	department := c.Query("department")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	list, err := ctl.svc.GetExcellentSubmissions(department)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	// 简单分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(list) {
		list = []models.Submission{}
	} else if end > len(list) {
		list = list[start:]
	} else {
		list = list[start:end]
	}

	var respList []gin.H
	for _, s := range list {
		homework, _ := ctl.svc.GetAssignmentDetail(s.HomeworkID)
		student, _ := ctl.svc.GetUserByID(s.StudentID)
		respList = append(respList, gin.H{
			"id": s.ID,
			"homework": gin.H{
				"id":               homework.ID,
				"title":            homework.Title,
				"department":       homework.Department,
				"department_label": ctl.DepartmentLabel(homework.Department),
			},
			"student": gin.H{
				"id":       student.ID,
				"nickname": student.Nickname,
			},
			"score":   s.Score,
			"comment": s.Comment,
		})
	}

	ctl.Success(c, gin.H{
		"list":      respList,
		"total":     len(list), // 实际应返回总数，此处简化
		"page":      page,
		"page_size": pageSize,
	})
}
