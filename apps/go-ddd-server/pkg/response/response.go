package response

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Success: true, Data: data})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{Success: false, Error: message})
}

func BadRequest(c *gin.Context, message string) { Error(c, http.StatusBadRequest, message) }
func NotFound(c *gin.Context, message string)   { Error(c, http.StatusNotFound, message) }
func Conflict(c *gin.Context, message string)   { Error(c, http.StatusConflict, message) }
func Internal(c *gin.Context, message string)   { Error(c, http.StatusInternalServerError, message) }

func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make([]string, 0, len(ve))
			for _, fe := range ve {
				fields = append(fields, fe.Field()+" is "+fe.Tag())
			}
			BadRequest(c, strings.Join(fields, "; "))
		} else {
			BadRequest(c, "invalid request body")
		}
		return false
	}
	return true
}
