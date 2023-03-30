package handlers

import (
	"errors"
	JobInterface "job-portal-lite/domain/jobs/interfaces"
	JobUseCase "job-portal-lite/domain/jobs/usecases"
	UserEntities "job-portal-lite/domain/user/entities"
	UserInterface "job-portal-lite/domain/user/interfaces"
	UserUseCase "job-portal-lite/domain/user/usecases"
	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"
	"job-portal-lite/shared/middlewares"
	"job-portal-lite/shared/utils"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type jobHandler struct {
	jobUseCase  JobInterface.JobUseCase
	userUseCase UserInterface.AuthUseCase
}

func NewJobHandler(router *gin.Engine) {
	handler := jobHandler{
		jobUseCase:  JobUseCase.NewJobUseCase(),
		userUseCase: UserUseCase.NewUserUseCase(),
	}
	jobs := router.Group("/jobs")
	{
		jobs.GET("", middlewares.Authenticated(), handler.GetJobLists)
		jobs.POST("", middlewares.Authenticated(), middlewares.CompanyRole(), handler.CreateJob)
		jobs.GET("/:id", middlewares.Authenticated(), handler.GetJobDetail)
		jobs.PUT("/:id", middlewares.Authenticated(), middlewares.CompanyRole(), handler.UpdateJob)
		jobs.DELETE("/:id", middlewares.Authenticated(), middlewares.CompanyRole(), handler.DeleteJob)
	}
}

// VALIDATION
func validateOpenDate(fl validator.FieldLevel) bool {
	openDate, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	closeDate := fl.Parent().Elem().FieldByName("CloseDate").Interface().(time.Time)
	if openDate.After(closeDate) {
		return false
	}
	return true
}

func noSpaces(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^\s*$`).MatchString(fl.Field().String())
}

func validateUpsert(request *requests.UpsertJob) error {
	validate := validator.New()
	var err error

	validate.RegisterValidation("noSpaces", noSpaces)
	err = validate.Struct(request)
	if err != nil {
		return errors.New("No space input")
	}

	validate.RegisterValidation("validateOpenDate", validateOpenDate)
	err = validate.Struct(request)
	if err != nil {
		return errors.New("The open date cannot be later than the close date")
	}

	return nil
}

// END VALIDATION

func (handler jobHandler) CreateJob(ctx *gin.Context) {
	request := &requests.UpsertJob{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validateUpsert(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseData, err := handler.jobUseCase.CreateJob(ctx, request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func (handler jobHandler) UpdateJob(ctx *gin.Context) {
	request := &requests.UpsertJob{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validateUpsert(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseData, err := handler.jobUseCase.UpdateJob(ctx, utils.ParseInt(ctx.Param("id")), request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, responseData)
}

func (handler jobHandler) DeleteJob(ctx *gin.Context) {
	responseData, err := handler.jobUseCase.DeleteJob(ctx, utils.ParseInt(ctx.Param("id")))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, responseData)
}

func (handler jobHandler) GetJobLists(ctx *gin.Context) {
	validate := validator.New()
	filterRequest := &requests.JobListsFilter{}
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

	var jobsInterface interface{}
	jobsInterface, err = handler.jobUseCase.GetJobLists(ctx, filterRequest)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobListsResponse := &responses.JobListsResponse{}
	jobListsResponse.Data = jobsInterface
	jobListsResponse.Status = "OK"
	jobListsResponse.TotalData = len(jobListsResponse.Data.([]map[string]interface{}))

	ctx.JSON(http.StatusOK, jobListsResponse)
}

func (handler jobHandler) GetJobDetail(ctx *gin.Context) {
	decoded := ctx.MustGet("Authenticated").(*responses.TokenDecoded)
	user := &UserEntities.User{
		Email: decoded.Email,
	}
	userData, err := handler.userUseCase.GetCurrentUser(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var job interface{}
	var jobDetailError error

	if userData.Data.Role == "Company" {
		jobInterface, err := handler.jobUseCase.GetJobDetailCompany(ctx, utils.ParseInt(ctx.Param("id")))
		if err != nil {
			jobDetailError = err
		} else {
			job = jobInterface.(*responses.JobDetailCompanyResponse)
		}
	} else {
		jobInterface, err := handler.jobUseCase.GetJobDetailCandidate(ctx, utils.ParseInt(ctx.Param("id")))
		if err != nil {
			jobDetailError = err
		} else {
			job = jobInterface.(*responses.JobDetailCandidateResponse)
		}
	}

	if jobDetailError != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": jobDetailError.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, job)
}
