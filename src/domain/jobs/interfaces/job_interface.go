package interfaces

import (
	"job-portal-lite/domain/jobs/entities"
	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"

	gin "github.com/gin-gonic/gin"
)

type JobUseCase interface {
	CreateJob(ctx *gin.Context, requests *requests.UpsertJob) (*responses.JobResponse, error)
	UpdateJob(ctx *gin.Context, id int, requests *requests.UpsertJob) (*responses.JobResponse, error)
	DeleteJob(ctx *gin.Context, id int) (*responses.JobResponse, error)
	GetJobLists(ctx *gin.Context, request *requests.JobListsFilter) (interface{}, error)
	GetJobDetailCompany(ctx *gin.Context, id int) (interface{}, error)
	GetJobDetailCandidate(ctx *gin.Context, id int) (interface{}, error)
}

type JobRepository interface {
	CreateJob(*entities.Jobs) (*entities.Jobs, error)
	UpdateJob(*entities.Jobs) (*entities.Jobs, error)
	DeleteJob(id int) error
	GetJobLists(request *requests.JobListsFilter) (interface{}, error)
	GetCompanyJobLists(id int, request *requests.JobListsFilter) (interface{}, error)
	GetJobDetail(id int) (*entities.Jobs, error)
}
