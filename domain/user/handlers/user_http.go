package handlers

import (
	"job-portal-lite/domain/user/entities"
	"job-portal-lite/domain/user/interfaces"
	"job-portal-lite/domain/user/usecases"
	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"
	errorResponse "job-portal-lite/models/responses/errors"
	"job-portal-lite/shared/middlewares"
	validators "job-portal-lite/shared/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	uc interfaces.AuthUseCase
}

func NewUserHandler(router *gin.Engine) {
	handler := userHandler{
		uc: usecases.NewUserUseCase(),
	}
	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.RegisterUser)
		auth.POST("/login", handler.LoginUser)
		auth.GET("/user", middlewares.Authenticated(), handler.GetCurrentUser)
	}
}

func (handler userHandler) RegisterUser(ctx *gin.Context) {
	validate := validator.New()

	request := &requests.RegisterAccountRequest{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validate.Struct(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseData, err := handler.uc.RegisterUser(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func (handler userHandler) LoginUser(ctx *gin.Context) {
	request := &requests.LoginAccountRequest{}

	jsonErrors := ctx.ShouldBindJSON(&request)
	if jsonErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   jsonErrors.Error(),
			"Message": "error",
		})
		return
	}

	structErrors := validators.ValidateStruct(request)
	if structErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &errorResponse.FormErrors{
			Errors: structErrors,
			BasicResponse: errorResponse.BasicResponse{
				Message:    "Harap isi form dengan benar",
				StatusCode: 400,
			},
		})
		return
	}

	responseData, err := handler.uc.LoginUser(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, responseData)
}

func (handler userHandler) GetCurrentUser(ctx *gin.Context) {
	decoded := ctx.MustGet("Authenticated").(*responses.TokenDecoded)
	user := &entities.User{
		Email: decoded.Email,
	}
	data, err := handler.uc.GetCurrentUser(user)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
