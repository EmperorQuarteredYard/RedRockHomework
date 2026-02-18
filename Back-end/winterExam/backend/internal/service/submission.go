package service

import (
	"homeworkSystem/backend/internal/models"
	WEjwt "homeworkSystem/backend/pkg/middleware/jwt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Service) SubmitSubmission() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var (
			message struct {
				HomeworkID uint64 `json:"homework_id"`
				Content    string `json:"content"`
				FileURL    string `json:"file_url"`
			}
			err        error
			assignment *models.Assignment
			authUser   WEjwt.AuthUser
			timeout    bool
		)

		defer func() {
			if err != nil {
				delegation.handleErrorAbort(err)
			}
		}()
		delegation.handleDataBind(&message)
		delegation.handleUserVerify(&authUser)
		if !delegation.Success() {
			return
		}

		if assignment, err = s.assignmentDao.FindByID(message.HomeworkID); err != nil {
			return
		}
		delegation.handelDepartment(assignment.Department, authUser.Department)
		timeout = time.Now().After(assignment.Deadline)
		if timeout && !assignment.AllowLate {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    StatusSubmitLate,
				"message": "fail to submit submission,due to timeout",
				"data":    nil,
			})
			return
		}
	}
}

//func (s *Service) CreateSubmission(submission *models.Submission) error {
//	return s.submissionDao.Create(submission)
//}
//
//func (s *Service) FindSubmissionBySubmitter(submitterID uint64) (sub *[]models.Submission, err error) {
//	return s.submissionDao.FindBySubmitterID(submitterID)
//}
//
//func (s *Service) FindSubmissionByDepartment(department string) (sub *[]models.Submission, err error) {
//	return s.submissionDao.FindByDepartment(department)
//}
//
//func (s *Service) SetCommentAndScore(comment string, score int, submissionID uint64, user *models.User) error {
//	return s.submissionDao.SetCommentAndScore(comment, score, submissionID, user)
//}
//
//func (s *Service) SetExcellentSubmission(submissionID uint64, user *models.User) error {
//	return s.submissionDao.SetExcellent(submissionID, user)
//}
//
//func (s *Service) GetExcellentSubmission(department string) (sub *[]models.Submission, err error) {
//	return s.submissionDao.GetExcellent(department)
//}
