package service

import "homeworkSystem/backend/internal/models"

func (s *Service) CreateSubmission(submission *models.Submission) error {
	return s.submissionDao.Create(submission)
}

func (s *Service) FindSubmissionBySubmitter(submitterID uint64) (sub *[]models.Submission, err error) {
	return s.submissionDao.FindBySubmitterID(submitterID)
}

func (s *Service) FindSubmissionByDepartment(department string) (sub *[]models.Submission, err error) {
	return s.submissionDao.FindByDepartment(department)
}

func (s *Service) SetCommentAndScore(comment string, score int, submissionID uint64, user *models.User) error {
	return s.submissionDao.SetCommentAndScore(comment, score, submissionID, user)
}

func (s *Service) SetExcellentSubmission(submissionID uint64, user *models.User) error {
	return s.submissionDao.SetExcellent(submissionID, user)
}

func (s *Service) GetExcellentSubmission(department string) (sub *[]models.Submission, err error) {
	return s.submissionDao.GetExcellent(department)
}
