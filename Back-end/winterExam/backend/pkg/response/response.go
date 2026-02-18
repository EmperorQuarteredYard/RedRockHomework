package response

import (
	"homeworkSystem/backend/pkg/errcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {
	// 根据错误类型获取对应的错误码和信息
	code, msg := errcode.FromError(err)
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}
