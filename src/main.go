package main

import (
	applications "job-portal-lite/domain/applications/handlers"
	jobs "job-portal-lite/domain/jobs/handlers"
	user "job-portal-lite/domain/user/handlers"
	"job-portal-lite/shared/databases"
	"job-portal-lite/shared/middlewares"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	databases.Migrate()
	router := gin.Default()

	router.Use(middlewares.NewCors(router))

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome to job portal lite",
		})
	})

	user.NewUserHandler(router)
	jobs.NewJobHandler(router)
	applications.NewApplicationHandler(router)
	router.Run(":" + os.Getenv("PORT"))
}
