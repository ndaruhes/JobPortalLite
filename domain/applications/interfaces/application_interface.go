package interfaces

import (
	"job-portal-lite/domain/applications/entities"
	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"

	"github.com/gin-gonic/gin"
)

type ApplicationUseCase interface {
	CreateApplication(ctx *gin.Context, request *requests.CreateApplicationRequest) (*responses.ApplicationResponse, error)
	ReadAllApplication(id int, request *requests.ApplicationListsFilter) (interface{}, error)
	ApplicationDetailCandidate(ctx *gin.Context, role string, id int) (*responses.ApplicationDetailCandidateResponse, error)
	ApplicationDetailCompany(ctx *gin.Context, role string, id int) (*responses.ApplicationDetailCompanyResponse, error)
	ProceedApplication(id int, ctx *gin.Context, request *requests.ProceedApplicatioRequest) (*responses.ApplicationResponse, error)
	CancelApplication(ctx *gin.Context, id int) (*responses.ApplicationResponse, error)
}

type ApplicationRepository interface {
	CreateApplication(*entities.Applications) (*entities.Applications, error)
	ReadAllApplication(int, int, *requests.ApplicationListsFilter) (interface{}, error)
	ApplicationDetail(role string, id int) (*entities.Applications, error)
	ProceedApplication(*entities.ApplicationHistories) (*entities.ApplicationHistories, error)
	CreateApplicationHistory(*entities.ApplicationHistories) (*entities.ApplicationHistories, error)
	ShowJobApplicationHistory(id int) (interface{}, error)
}
