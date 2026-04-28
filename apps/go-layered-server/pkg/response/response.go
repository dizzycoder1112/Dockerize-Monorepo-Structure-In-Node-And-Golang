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

func BadRequest(c *gin.Context, message string)          { Error(c, http.StatusBadRequest, message) }
func Unauthorized(c *gin.Context, message string)        { Error(c, http.StatusUnauthorized, message) }
func NotFound(c *gin.Context, message string)            { Error(c, http.StatusNotFound, message) }
func Conflict(c *gin.Context, message string)            { Error(c, http.StatusConflict, message) }
func UnprocessableEntity(c *gin.Context, message string) { Error(c, http.StatusUnprocessableEntity, message) }

// InternalServerError responds with a generic 500 to the client and stashes the
// real message + caller location on the gin context for upstream logging
// middleware to read. The client never sees the raw error.
func InternalServerError(c *gin.Context, message string) {
	c.Set("error_message", message)
	Error(c, http.StatusInternalServerError, "Internal Server Error")
}

// BindJSON binds the JSON body and validates struct tags. On failure it writes
// a 400 with a snake_cased field listing and returns false.
func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make([]string, 0, len(ve))
			for _, fe := range ve {
				fields = append(fields, formatFieldError(fe))
			}
			BadRequest(c, strings.Join(fields, "; "))
		} else {
			BadRequest(c, "invalid request body")
		}
		return false
	}
	return true
}

func formatFieldError(fe validator.FieldError) string {
	field := toSnakeCase(fe.Field())
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "gt":
		return field + " must be greater than " + fe.Param()
	case "gte":
		return field + " must be >= " + fe.Param()
	case "min":
		return field + " must have at least " + fe.Param() + " items"
	default:
		return field + " is invalid"
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				prev := runes[i-1]
				if prev >= 'a' && prev <= 'z' || (prev >= 'A' && prev <= 'Z' && i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z') {
					result.WriteByte('_')
				}
			}
			result.WriteRune(r + 32)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
