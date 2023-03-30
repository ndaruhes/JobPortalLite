package usecases

import (
	"errors"
	"job-portal-lite/domain/applications/entities"
	ApplicationInterface "job-portal-lite/domain/applications/interfaces"
	ApplicationRepository "job-portal-lite/domain/applications/repositories"
	JobInterface "job-portal-lite/domain/jobs/interfaces"
	JobRepository "job-portal-lite/domain/jobs/repositories"
	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"
	"time"

	"github.com/gin-gonic/gin"
)

type applicationUseCase struct {
	applicationRepo ApplicationInterface.ApplicationRepository
	jobRepo         JobInterface.JobRepository
}

func NewApplicationUseCase() *applicationUseCase {
	return &applicationUseCase{
		applicationRepo: ApplicationRepository.NewApplicationRepository(),
		jobRepo:         JobRepository.NewJobRepository(),
	}
}

func (a applicationUseCase) CreateApplication(ctx *gin.Context, request *requests.CreateApplicationRequest) (*responses.ApplicationResponse, error) {
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	if user.Role == "Company" {
		return nil, errors.New("Company couldn't apply")
	}

	job, err := a.jobRepo.GetJobDetail(request.JobId)
	if err != nil {
		return nil, errors.New("Job not found")
	}

	// historiesCheck, err := a.applicationRepo.ReadAllApplication()
	if err != nil {
		return nil, errors.New("Application not found")
	}

	// isApplied := false
	// for _, historyCheck := range historiesCheck.([]map[string]interface{}) {
	// 	if historyCheck["status"] == "Applied" {
	// 		isApplied = true
	// 	}
	// }

	// if isApplied {
	// 	return nil, errors.New("You already apply this job")
	// }

	application, err := a.applicationRepo.CreateApplication(&entities.Applications{
		CandidateId: user.ID,
		JobId:       request.JobId,
		CompanyId:   job.CompanyId,
	})

	if err != nil {
		response := &responses.ApplicationResponse{
			Message: err.Error(),
			Status:  "Bad Request",
		}

		return response, err
	}

	_, err = a.applicationRepo.CreateApplicationHistory(&entities.ApplicationHistories{
		LastProcessDate: time.Now(),
		ApplicationId:   application.ID,
	})

	if err != nil {
		response := &responses.ApplicationResponse{
			Message: "Can't Create Application History",
			Status:  "Bad Request",
		}

		return response, err
	}

	response := &responses.ApplicationResponse{
		Message: "Application Applied Succesfully",
		Status:  "OK",
	}

	return response, nil
}

func (a applicationUseCase) ReadAllApplication(userId int, request *requests.ApplicationListsFilter) (interface{}, error) {
	if request.Size == 0 {
		request.Size = 5
	} else if request.Page == 0 {
		request.Page = 1
	}

	applications, err := a.applicationRepo.ReadAllApplication(userId, 0, request)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (a applicationUseCase) ApplicationDetailCandidate(ctx *gin.Context, role string, id int) (*responses.ApplicationDetailCandidateResponse, error) {
	application, err := a.applicationRepo.ApplicationDetail(role, id)
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	if application.CandidateId != user.ID {
		return nil, errors.New("Candidate not authorized to see this application")
	}
	if err != nil {
		return nil, errors.New("Application not found")
	}

	var historiesInterface interface{}
	historiesInterface, err = a.applicationRepo.ShowJobApplicationHistory(application.ID)
	if err != nil {
		return nil, errors.New("Application Histories Error")
	}

	return &responses.ApplicationDetailCandidateResponse{
		ApplicationDetail: responses.ApplicationDetail{
			Id:        application.ID,
			JobTitle:  application.Job.Title,
			Processes: historiesInterface,
		},
		Company: application.Company.Name,
		Status:  "OK",
	}, nil
}

func (a applicationUseCase) ApplicationDetailCompany(ctx *gin.Context, role string, id int) (*responses.ApplicationDetailCompanyResponse, error) {
	application, err := a.applicationRepo.ApplicationDetail(role, id)
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	if application.CompanyId != user.ID {
		return nil, errors.New("Company not authorized to see this application")
	}
	if err != nil {
		return nil, errors.New("Application not found")
	}

	var historiesInterface interface{}
	historiesInterface, err = a.applicationRepo.ShowJobApplicationHistory(application.ID)
	if err != nil {
		return nil, errors.New("Application Histories Error")
	}

	return &responses.ApplicationDetailCompanyResponse{
		ApplicationDetail: responses.ApplicationDetail{
			Id:        application.ID,
			JobTitle:  application.Job.Title,
			Processes: historiesInterface,
		},
		Candidate: application.Candidate.Name,
		Status:    "OK",
	}, nil
}

func (a applicationUseCase) ProceedApplication(id int, ctx *gin.Context, request *requests.ProceedApplicatioRequest) (*responses.ApplicationResponse, error) {
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	application, err := a.applicationRepo.ApplicationDetail(user.Role, id)
	if err != nil {
		return nil, errors.New("Application not found")
	}

	if application.CompanyId != user.ID {
		return nil, errors.New("Company not authorized proceed this application")
	}

	historiesCheck, err := a.applicationRepo.ShowJobApplicationHistory(id)
	if err != nil {
		return nil, errors.New("Application not found")
	}

	isRejected := false
	var currStatus string
	for _, historyCheck := range historiesCheck.([]map[string]interface{}) {
		if historyCheck["status"] == "Rejected" || historyCheck["status"] == "Cancelled" {
			isRejected = true
			currStatus = historyCheck["status"].(string)
		}
	}

	if isRejected {
		return nil, errors.New("Cannot proceed application again because it was " + currStatus)
	}

	var status string
	switch request.StatusId {
	case 1:
		status = "HR Interview"
		break
	case 2:
		status = "Client Interview"
		break
	case 3:
		status = "Passed"
		break
	case 4:
		status = "Rejected"
		break
	}
	_, err = a.applicationRepo.ProceedApplication(&entities.ApplicationHistories{
		Status:          status,
		LastProcessDate: time.Now(),
		ApplicationId:   application.ID,
	})

	if err != nil {
		response := &responses.ApplicationResponse{
			Message: "Can't Proceed Application",
			Status:  "Bad Request",
		}

		return response, err
	}

	response := &responses.ApplicationResponse{
		Message: "Proceed Application Succesfully",
		Status:  "OK",
	}

	return response, nil
}

func (a applicationUseCase) CancelApplication(ctx *gin.Context, id int) (*responses.ApplicationResponse, error) {
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	application, err := a.applicationRepo.ApplicationDetail(user.Role, id)
	if err != nil {
		return nil, errors.New("Application not found")
	}

	if user.Role == "Company" {
		return nil, errors.New("Company not authorized to cancel this application")
	}

	if application.CandidateId != user.ID {
		return nil, errors.New("Candidate not authorized to cancel this application")
	}

	_, err = a.applicationRepo.ProceedApplication(&entities.ApplicationHistories{
		Status:          "Cancelled",
		LastProcessDate: time.Now(),
		ApplicationId:   application.ID,
	})

	if err != nil {
		response := &responses.ApplicationResponse{
			Message: "Can't Cancel Your Application",
			Status:  "Bad Request",
		}

		return response, err
	}

	response := &responses.ApplicationResponse{
		Message: "Application Cancelled Succesfully",
		Status:  "OK",
	}

	return response, nil
}
