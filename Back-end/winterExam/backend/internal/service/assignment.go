package service

import (
	"homeworkSystem/backend/internal/models"
	WEjwt "homeworkSystem/backend/pkg/middleware/jwt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Service) FindAssignmentByDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var (
			message = struct {
				Department string `json:"department"`
				Page       int    `json:"page"`
				PageSize   int    `json:"page_size"`
			}{
				Page:       1,
				PageSize:   10,
				Department: "",
			}
			department string
			authUser   WEjwt.AuthUser
		)

		delegation.handleDataBind(&message)
		delegation.handleDataForm(message.PageSize < 10 || message.PageSize > 100)
		delegation.handleUserVerify(&authUser)
		if !delegation.Success() {
			return
		}

		if message.Department != "" {
			department = message.Department
		} else {
			department = authUser.Department
		}
		list, total, err := s.assignmentDao.FindByDepartment(message.Page, message.PageSize, department)
		if err != nil {
			delegation.handleErrorAbort(err)
			return
		}
		delegation.handleSuccessResponse(gin.H{
			"list":  list,
			"total": total,
			"page":  message.Page,
			"size":  message.PageSize,
		})
		return
	}
}

func (s *Service) PublishAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var (
			authUser WEjwt.AuthUser
			err      error
			ass      *models.Assignment
			deadline time.Time
			label    string
			message  = struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Department  string `json:"department"`
				Deadline    string `json:"deadline"`
				AllowLate   bool   `json:"allow_late"`
			}{
				AllowLate: false,
			}
		)

		delegation.handleUserVerify(&authUser)
		delegation.handleDataBind(&message)
		delegation.handleDataForm(message.Title == "" || message.Description == "" || message.Department == "" || message.Deadline == "")
		delegation.handleLabelParse(message.Department, &label)

		deadline, err = time.Parse("2006-01-02 15:04:05", message.Deadline)
		delegation.handleDataForm(err != nil)
		if !delegation.Success() {
			return
		}

		ass = &models.Assignment{
			Department:  message.Department,
			Title:       message.Title,
			Description: message.Description,
			AllowLate:   message.AllowLate,
			Deadline:    deadline,
		}
		if err = s.assignmentDao.Publish(ass); err != nil {
			delegation.handleErrorAbort(err)
			return
		}
		delegation.handleSuccessResponse(gin.H{
			"id":               ass.ID,
			"title":            ass.Title,
			"department":       message.Department,
			"department_label": label,
			"deadline":         message.Deadline,
			"allowLate":        message.AllowLate,
		})
	}
}

func (s *Service) FindAssignmentByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var (
			authUser         WEjwt.AuthUser
			user             *models.User
			ass              *models.Assignment
			err              error
			id               uint64
			submission_count int
			submissions      *[]models.Submission
			submission       *models.Submission
			label            string
			my_submission    map[string]interface{}
		)
		defer func() {
			if err != nil {
				delegation.handleErrorAbort(err)
			}
		}()

		delegation.handleUserVerify(&authUser)
		if id, err = strconv.ParseUint(c.Param("id"), 10, 64); err != nil {
			return
		}
		if user, err = s.userDao.FindByID(authUser.UserID); err != nil {
			return
		}
		if ass, err = s.assignmentDao.FindByID(id); err != nil {
			return
		}
		delegation.handleLabelParse(ass.Department, &label)
		submission_count = len(*submissions)
		if authUser.Role == models.RoleNewLight {
			if submissions, err = s.submissionDao.FindBySubmitterID(authUser.UserID); err != nil {
				return
			}
			if submissions != nil && len(*submissions) != 0 {
				var maxTime time.Time = (*submissions)[0].CreatedAt
				submission = &(*submissions)[0]
				for _, sub := range *submissions {
					if maxTime.Before(sub.CreatedAt) {
						maxTime = sub.CreatedAt
						submission = &sub
					}
				}
				my_submission = map[string]interface{}{
					"id":           submission.ID,
					"score":        submission.Score,
					"is_excellent": submission.Excellent,
				}
			} else {
				submission = nil
			}
		}

		delegation.handleSuccessResponse(gin.H{
			"id":               ass.ID,
			"title":            ass.Title,
			"description":      ass.Description,
			"department":       ass.Department,
			"department_label": label,
			"creator": gin.H{
				"id":       user.ID,
				"nickname": user.Nickname,
			},
			"deadline":         ass.Deadline.Format("2006-01-02 15:04:05"),
			"allow_late":       ass.AllowLate,
			"submission_count": submission_count,
			"my_submission":    my_submission,
		})
		return
	}
}

func (s *Service) UpdateAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var (
			id       uint64
			err      error
			ass      *models.Assignment
			authUser WEjwt.AuthUser
			user     *models.User
		)

		defer func() {
			if err != nil {
				delegation.handleErrorAbort(err)
			}
		}()

		delegation.handleUserVerify(&authUser)
		delegation.handlePermissionRole(authUser.Role, models.RoleNewLight)
		if id, err = strconv.ParseUint(c.Param("id"), 10, 64); err != nil {
			return
		}
		if !delegation.Success() {
			return
		}
		if ass, err = s.assignmentDao.FindByID(id); err != nil {
			return
		}
		if user, err = s.userDao.FindByID(authUser.UserID); err != nil {
			return
		}
		delegation.handleDataBind(&ass)
		if err = s.assignmentDao.Update(user, ass); err != nil {
			return
		}
		delegation.handleSuccessResponse(gin.H{
			"id":       id,
			"title":    ass.Title,
			"deadline": ass.Deadline.Format("2006-01-02 15:04:05"),
		})
		return
	}
}

func (s *Service) DeleteAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var (
			id       uint64
			err      error
			authUser WEjwt.AuthUser
			user     *models.User
		)

		defer func() {
			if err != nil {
				delegation.handleErrorAbort(err)
			}
		}()

		id, err = strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return
		}
		delegation.handleUserVerify(&authUser)
		user, err = s.userDao.FindByID(authUser.UserID)
		if err != nil {
			return
		}

		err = s.assignmentDao.Delete(user, id)
		if err != nil {
			return
		}
		delegation.handleSuccessResponse(nil)
	}
}

//func (s *Service) FindAssignmentByID(id int64) (*models.Assignment, error) {
//	return s.assignmentDao.FindByID(id)
//}
//
//func (s *Service) UpdateAssignment(Updater *models.User, newOne models.Assignment) error {
//	return s.assignmentDao.Update(Updater, newOne)
//}
//
//func (s *Service) DeleteAssignment(deleter *models.User, id int64) error {
//	return s.assignmentDao.Delete(deleter, id)
//}
