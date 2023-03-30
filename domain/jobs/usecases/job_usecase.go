package usecases

import (
	"errors"
	ApplicationInterface "job-portal-lite/domain/applications/interfaces"
	ApplicationRepository "job-portal-lite/domain/applications/repositories"
	"job-portal-lite/domain/jobs/entities"
	JobInterface "job-portal-lite/domain/jobs/interfaces"
	JobRepository "job-portal-lite/domain/jobs/repositories"
	"job-portal-lite/models/requests"
	"job-portal-lite/models/responses"

	"github.com/gin-gonic/gin"
)

type jobUseCase struct {
	jobRepository         JobInterface.JobRepository
	applicationRepository ApplicationInterface.ApplicationRepository
}

func NewJobUseCase() *jobUseCase {
	return &jobUseCase{
		jobRepository:         JobRepository.NewJobRepository(),
		applicationRepository: ApplicationRepository.NewApplicationRepository(),
	}
}

func (uc *jobUseCase) CreateJob(ctx *gin.Context, request *requests.UpsertJob) (*responses.JobResponse, error) {
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	_, err := uc.jobRepository.CreateJob(&entities.Jobs{
		Title:       request.Title,
		Description: request.Description,
		OpenDate:    request.OpenDate,
		CloseDate:   request.CloseDate,
		CompanyId:   user.ID,
	})

	if err != nil {
		response := &responses.JobResponse{
			Message: err.Error(),
			Status:  "Bad Request",
		}

		return response, err
	}

	response := &responses.JobResponse{
		Message: "Job created succesfully",
		Status:  "OK",
	}

	return response, nil
}

func (uc *jobUseCase) UpdateJob(ctx *gin.Context, id int, request *requests.UpsertJob) (*responses.JobResponse, error) {
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	jobDetail, err := uc.jobRepository.GetJobDetail(id)
	if err != nil {
		return nil, errors.New("Job not found")
	}

	if jobDetail.CompanyId != user.ID {
		return nil, errors.New("Company not authorized to update this job")
	}

	_, err = uc.jobRepository.UpdateJob(&entities.Jobs{
		ID:          jobDetail.ID,
		Title:       request.Title,
		Description: request.Description,
		OpenDate:    request.OpenDate,
		CloseDate:   request.CloseDate,
		CompanyId:   user.ID,
	})

	if err != nil {
		response := &responses.JobResponse{
			Message: err.Error(),
			Status:  "Bad Request",
		}

		return response, err
	}

	response := &responses.JobResponse{
		Message: "Job updated succesfully",
		Status:  "OK",
	}

	return response, nil
}

func (uc *jobUseCase) DeleteJob(ctx *gin.Context, id int) (*responses.JobResponse, error) {
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	jobDetail, err := uc.jobRepository.GetJobDetail(id)
	if err != nil {
		return nil, errors.New("Job not found")
	}

	if jobDetail.CompanyId != user.ID {
		return nil, errors.New("Company not authorized to delete this job")
	}

	uc.jobRepository.DeleteJob(id)

	response := &responses.JobResponse{
		Message: "Job deleted succesfully",
		Status:  "OK",
	}

	return response, nil
}

func (uc *jobUseCase) GetJobLists(ctx *gin.Context, request *requests.JobListsFilter) (interface{}, error) {
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	var jobs interface{}
	var err error
	if user.Role == "Company" {
		jobs, err = uc.jobRepository.GetCompanyJobLists(user.ID, request)
	} else {
		if request.Size == 0 {
			request.Size = 6
		} else if request.Page == 0 {
			request.Page = 1
		}
		jobs, err = uc.jobRepository.GetJobLists(request)
	}

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (uc *jobUseCase) GetJobDetailCompany(ctx *gin.Context, id int) (interface{}, error) {
	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	job, err := uc.jobRepository.GetJobDetail(id)
	if (err != nil) || (user.ID != job.CompanyId) {
		return nil, errors.New("Job not found")
	}

	applications, applicationsErr := uc.applicationRepository.ReadAllApplication(user.ID, job.ID, &requests.ApplicationListsFilter{})
	if applicationsErr != nil {
		return nil, nil
	}

	return &responses.JobDetailCompanyResponse{
		Data: responses.JobDetailCompany{
			ID:             job.ID,
			Title:          job.Title,
			Description:    job.Description,
			ApplicantCount: len(applications.([]map[string]interface{})),
			OpenDate:       job.OpenDate,
			CloseDate:      job.CloseDate,
			Company:        job.Company.Name,
			Applications:   applications,
		},
		JobResponse: responses.JobResponse{
			Message: "Job Detail Showed",
			Status:  "OK",
		},
	}, nil
}

func (uc *jobUseCase) GetJobDetailCandidate(ctx *gin.Context, id int) (interface{}, error) {
	job, err := uc.jobRepository.GetJobDetail(id)
	if err != nil {
		return nil, errors.New("Job not found")
	}

	user := ctx.Value("Authenticated").(*responses.TokenDecoded)
	applications, applicationsErr := uc.applicationRepository.ReadAllApplication(user.ID, job.ID, &requests.ApplicationListsFilter{})
	if applicationsErr != nil {
		return nil, nil
	}

	return &responses.JobDetailCandidateResponse{
		Data: responses.JobDetailCandidate{
			ID:             job.ID,
			Title:          job.Title,
			Description:    job.Description,
			ApplicantCount: len(applications.([]map[string]interface{})),
			OpenDate:       job.OpenDate,
			CloseDate:      job.CloseDate,
			Company:        job.Company.Name,
		},
		JobResponse: responses.JobResponse{
			Message: "Job Detail Showed",
			Status:  "OK",
		},
	}, nil
}
