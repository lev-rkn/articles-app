package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, code int, err error) {
	var errMessage string
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		errMessage = err.Error()
	} else {
		errMessage = unwrapped.Error()
	}

	// TODO: unwrap rpc error

	c.JSON(code, gin.H{
		"error": errMessage,
	})
}
