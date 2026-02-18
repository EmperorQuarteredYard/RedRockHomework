package controller

import (
	"homeworkSystem/backend/internal/service"
	"homeworkSystem/backend/pkg/errcode"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type HomeworkController struct {
	BaseController
	svc *service.Service
}

func NewHomeworkController(svc *service.Service) *HomeworkController {
	return &HomeworkController{svc: svc}
}

// Publish 发布作业
func (ctl *HomeworkController) Publish(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.Role != "admin" {
		ctl.HandleCode(c, errcode.ErrInsufficientPermissions)
		return
	}

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Department  string `json:"department" binding:"required"`
		Deadline    string `json:"deadline" binding:"required"` // 格式: 2006-01-02 15:04:05
		AllowLate   bool   `json:"allow_late"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	deadline, err := time.Parse("2006-01-02 15:04:05", req.Deadline)
	if err != nil {
		ctl.HandleCode(c, errcode.ErrInvalidParams)
		return
	}

	input := service.PublishAssignmentInput{
		Title:       req.Title,
		Description: req.Description,
		Department:  req.Department,
		Deadline:    deadline,
		AllowLate:   req.AllowLate,
		CreatorID:   authUser.UserID,
	}

	assignment, err := ctl.svc.PublishAssignment(input)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"id":               assignment.ID,
		"title":            assignment.Title,
		"department":       assignment.Department,
		"department_label": ctl.DepartmentLabel(assignment.Department),
		"deadline":         assignment.DeadlineString(),
		"allow_late":       assignment.AllowLate,
	})
}

// List 获取作业列表
func (ctl *HomeworkController) List(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	department := c.Query("department")
	if department == "" {
		department = authUser.Department // 默认当前用户部门
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	list, total, err := ctl.svc.ListAssignments(department, page, pageSize)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}
	var count int64
	var nickname string
	var respList []gin.H
	for _, a := range list {
		count, err = ctl.svc.CountSubmissionsByHomework(a.ID)
		if err != nil {
			count = -1
		}
		nickname, err = ctl.svc.GetNicknameByID(a.ID)
		if err != nil {
			nickname = "用户不存在"
		}
		respList = append(respList, gin.H{
			"id":               a.ID,
			"title":            a.Title,
			"department":       a.Department,
			"department_label": ctl.DepartmentLabel(a.Department),
			"creator": gin.H{
				"id":       a.CreatorID,
				"nickname": nickname,
			},
			"deadline":         a.DeadlineString(),
			"allow_late":       a.AllowLate,
			"submission_count": count,
		})
	}

	ctl.Success(c, gin.H{
		"list":      respList,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetDetail 获取作业详情
func (ctl *HomeworkController) GetDetail(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.ErrInvalidParams)
		return
	}

	assignment, err := ctl.svc.GetAssignmentDetail(id)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	// 获取发布者信息
	creator, _ := ctl.svc.GetUserByID(assignment.CreatorID) // 忽略错误

	// 小登需要看到自己的提交
	var mySubmission interface{}
	if authUser.Role == "student" {
		submissions, _ := ctl.svc.GetMySubmissions(authUser.UserID)
		// 查找对应作业的最新提交
		for _, s := range submissions {
			if s.HomeworkID == id {
				mySubmission = gin.H{
					"id":           s.ID,
					"score":        s.Score,
					"is_excellent": s.IsExcellent,
				}
				break
			}
		}
	}

	var count int64
	count, err = ctl.svc.CountSubmissionsByHomework(assignment.ID)
	if err != nil {
		count = -1
	}

	ctl.Success(c, gin.H{
		"id":               assignment.ID,
		"title":            assignment.Title,
		"description":      assignment.Description,
		"department":       assignment.Department,
		"department_label": ctl.DepartmentLabel(assignment.Department),
		"creator": gin.H{
			"id":       creator.ID,
			"nickname": creator.Nickname,
		},
		"deadline":         assignment.DeadlineString(),
		"allow_late":       assignment.AllowLate,
		"submission_count": count,
		"my_submission":    mySubmission,
	})
}

// Update 修改作业
func (ctl *HomeworkController) Update(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.Role != "admin" {
		ctl.HandleCode(c, errcode.ErrInsufficientPermissions)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.ErrInvalidParams)
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Deadline    *string `json:"deadline"`
		AllowLate   *bool   `json:"allow_late"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Deadline != nil {
		deadline, err := time.Parse("2006-01-02 15:04:05", *req.Deadline)
		if err != nil {
			ctl.HandleCode(c, errcode.ErrInvalidParams)
			return
		}
		updates["deadline"] = deadline
	}
	if req.AllowLate != nil {
		updates["allow_late"] = *req.AllowLate
	}

	// 获取当前用户完整信息
	updater, err := ctl.svc.GetUserByID(authUser.UserID)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	updated, err := ctl.svc.UpdateAssignment(updater, id, updates)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"id":       updated.ID,
		"title":    updated.Title,
		"deadline": updated.DeadlineString(),
	})
}

// Delete 删除作业
func (ctl *HomeworkController) Delete(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.Role != "admin" {
		ctl.HandleCode(c, errcode.ErrInsufficientPermissions)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.ErrInvalidParams)
		return
	}

	deleter, err := ctl.svc.GetUserByID(authUser.UserID)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	if err := ctl.svc.DeleteAssignment(deleter, id); err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, nil)
}
