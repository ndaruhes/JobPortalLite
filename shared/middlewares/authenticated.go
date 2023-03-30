package middlewares

import (
	"job-portal-lite/shared/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticated() gin.HandlerFunc {
	return func(context *gin.Context) {
		verifyToken, err := helpers.VerifyToken(context)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthenticated",
				"status":  http.StatusUnauthorized,
			})

			return
		}

		context.Set("Authenticated", verifyToken)
		context.Next()
	}
}
