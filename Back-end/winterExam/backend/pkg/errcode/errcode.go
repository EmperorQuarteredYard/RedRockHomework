package errcode

import (
	"errors"
	"homeworkSystem/backend/internal/repository"
	"homeworkSystem/backend/pkg/jwt"
)

const (
	Success                    = 0
	ErrInvalidParams           = 10001
	ErrUnauthorized            = 10002
	ErrInsufficientPermissions = 10003
	ErrDepartmentNotMatch      = 10004
	ErrSubmitLate              = 10005
	ErrNotFound                = 10006
	ErrDuplicateEntry          = 10007
	ErrServer                  = 30001
)

var codeMsgMap = map[int]string{
	Success:                           "success",
	ErrInvalidParams:                  "参数错误",
	ErrUnauthorized:                   "未认证",
	ErrInsufficientPermissions:        "权限不足",
	ErrDepartmentNotMatch:             "部门不匹配",
	ErrSubmitLate:                     "已过截止时间且不允许补交",
	ErrNotFound:                       "记录不存在",
	ErrDuplicateEntry:                 "记录已存在",
	ErrServer:                         "服务器内部错误",
	jwt.StatusInvalidToken:            "Token不合法",
	jwt.StatusInsufficientPermissions: "Token权限不足",
	jwt.StatusMissedToken:             "缺少Token",
	jwt.StatusInvalidUserData:         "用户数据格式不正确",
	jwt.StatusUserNotAuthenticated:    "用户未被授权",
}

func CodeMsg(code int) string {
	msg, ok := codeMsgMap[code]
	if ok {
		return msg
	}
	return "nil"
}

func FromError(err error) (int, string) {
	if err == nil {
		return Success, codeMsgMap[Success]
	}
	// 匹配自定义错误
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return ErrNotFound, codeMsgMap[ErrNotFound]
	case errors.Is(err, repository.ErrDepartmentNotMatch):
		return ErrDepartmentNotMatch, codeMsgMap[ErrDepartmentNotMatch]
	case errors.Is(err, repository.ErrInsufficientPermissions):
		return ErrInsufficientPermissions, codeMsgMap[ErrInsufficientPermissions]
	case errors.Is(err, repository.ErrDuplicateEntry):
		return ErrDuplicateEntry, codeMsgMap[ErrDuplicateEntry]
	// 其他错误...
	default:
		return ErrServer, codeMsgMap[ErrServer]
	}
}
