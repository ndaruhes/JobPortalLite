package handlers

import (
	ApplicationInterface "job-portal-lite/domain/applications/interfaces"
	ApplicationUseCase "job-portal-lite/domain/applications/usecases"
	UserEntities "job-portal-lite/domain/user/entities"
	UserInterface "job-portal-lite/domain/user/interfaces"
	UserUseCase "job-portal-lite/domain/user/usecases"
	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"
	"job-portal-lite/shared/middlewares"
	"job-portal-lite/shared/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type applicationHandler struct {
	applicationUseCase ApplicationInterface.ApplicationUseCase
	userUseCase        UserInterface.AuthUseCase
}

func NewApplicationHandler(router *gin.Engine) {
	handler := applicationHandler{
		applicationUseCase: ApplicationUseCase.NewApplicationUseCase(),
		userUseCase:        UserUseCase.NewUserUseCase(),
	}
	applications := router.Group("/applications")
	{
		applications.POST("", middlewares.Authenticated(), handler.CreateApplication)
		applications.GET("", middlewares.Authenticated(), handler.ReadApplication)
		applications.GET("/:id", middlewares.Authenticated(), handler.ApplicationDetail)
		applications.POST("/:id/proceed", middlewares.Authenticated(), handler.ProceedApplication)
		applications.POST("/:id/cancel", middlewares.Authenticated(), handler.CancelApplication)
	}
}

func (handler applicationHandler) CreateApplication(ctx *gin.Context) {
	validate := validator.New()
	request := &requests.CreateApplicationRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validate.Struct(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseData, err := handler.applicationUseCase.CreateApplication(ctx, request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func (handler applicationHandler) ReadApplication(ctx *gin.Context) {
	decoded := ctx.MustGet("Authenticated").(*responses.TokenDecoded)
	user := &UserEntities.User{
		ID:    decoded.ID,
		Email: decoded.Email,
	}

	validate := validator.New()
	filterRequest := &requests.ApplicationListsFilter{}
	err := ctx.ShouldBindQuery(filterRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validate.Struct(filterRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var applicationsInterface interface{}
	applicationsInterface, err = handler.applicationUseCase.ReadAllApplication(user.ID, filterRequest)
	// applicationsInterface, _ = handler.applicationUseCase.ReadAllApplication("apapaun@gmail.com' OR 1 = 1 -- ")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	applicationResponse := &responses.ReadApplicationResponse{}
	applicationResponse.Data = applicationsInterface
	applicationResponse.Message = "Applications List Showed"
	applicationResponse.Status = "OK"
	applicationResponse.TotalData = len(applicationResponse.Data.([]map[string]interface{}))

	ctx.JSON(http.StatusOK, applicationResponse)
}

func (handler applicationHandler) ApplicationDetail(ctx *gin.Context) {
	decoded := ctx.MustGet("Authenticated").(*responses.TokenDecoded)

	var err error
	var applicationDetailCandidate *responses.ApplicationDetailCandidateResponse
	var applicationDetailCompany *responses.ApplicationDetailCompanyResponse
	if decoded.Role == "Company" {
		applicationDetailCompany, err = handler.applicationUseCase.ApplicationDetailCompany(ctx, decoded.Role, utils.ParseInt(ctx.Param("id")))
	} else {
		applicationDetailCandidate, err = handler.applicationUseCase.ApplicationDetailCandidate(ctx, decoded.Role, utils.ParseInt(ctx.Param("id")))
	}

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if decoded.Role == "Company" {
		ctx.JSON(http.StatusOK, applicationDetailCompany)
	} else {
		ctx.JSON(http.StatusOK, applicationDetailCandidate)
	}
}

func (handler applicationHandler) ProceedApplication(ctx *gin.Context) {
	validate := validator.New()
	request := &requests.ProceedApplicatioRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validate.Struct(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseData, err := handler.applicationUseCase.ProceedApplication(utils.ParseInt(ctx.Param("id")), ctx, request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func (handler applicationHandler) CancelApplication(ctx *gin.Context) {
	responseData, err := handler.applicationUseCase.CancelApplication(ctx, utils.ParseInt(ctx.Param("id")))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, responseData)
}
