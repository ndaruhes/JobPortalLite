package middlewares

import (
	"job-portal-lite/models/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompanyRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		decoded := ctx.MustGet("Authenticated").(*responses.TokenDecoded)
		if decoded.Role != "Company" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden",
				"status":  http.StatusForbidden,
			})

			return
		}

		ctx.Next()
	}
}
