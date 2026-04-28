package middleware

import "github.com/gin-gonic/gin"

type Middlewares struct {
	Logger gin.HandlerFunc
}
