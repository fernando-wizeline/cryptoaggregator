package types

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	failureHeader = "x-failure-reason"
)

type jsonResponseFunc func() (result any, err error)

type RouterError struct {
	StatusCode int                    `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty" swaggerignore:"true"`
	Err        error                  `json:"-"`
}

func (r RouterError) Error() string {
	return r.Message
}

func JSONResponse(c *gin.Context, f jsonResponseFunc) {
	if result, err := f(); err != nil {
		AbortIfErr(c, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func AbortIfErr(c *gin.Context, err error) bool {
	if err != nil {
		_ = c.Error(err) //nolint:errcheck
		var routerErr RouterError
		if errors.As(err, &routerErr) {
			c.Header(failureHeader, routerErr.Message)
			c.AbortWithStatusJSON(routerErr.StatusCode, &routerErr)
		} else {
			c.Header(failureHeader, err.Error())
			routerErr = RouterError{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
				Err:        err,
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &routerErr)
		}
	}

	return err == nil
}
