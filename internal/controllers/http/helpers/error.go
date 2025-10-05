package errors

import "github.com/gin-gonic/gin"

type ErrorPayload struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func RespondError(c *gin.Context, status int, code, msg string) {
	var p ErrorPayload

	p.Error.Code, p.Error.Message = code, msg
	c.IndentedJSON(status, p)
}
