package repository

import "errors"

var (
	ErrInsufficientPermissions = errors.New("权限不足")
	ErrDepartmentNotMatch      = errors.New("部门不匹配")
	ErrNotFound                = errors.New("记录不存在")
	ErrDuplicateEntry          = errors.New("记录已存在")
	ErrSubmitLate              = errors.New("提交过晚")
)
