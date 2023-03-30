package repositories

import (
	"errors"
	"job-portal-lite/domain/jobs/entities"
	"job-portal-lite/models/requests"
	"job-portal-lite/shared/databases"

	"gorm.io/gorm"
)

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository() *jobRepository {
	return &jobRepository{
		db: databases.Connect(),
	}
}

func (j jobRepository) CreateJob(job *entities.Jobs) (*entities.Jobs, error) {
	err := j.db.Save(job).Error
	if err != nil {
		return nil, err
	}

	return job, err
}

func (j jobRepository) UpdateJob(job *entities.Jobs) (*entities.Jobs, error) {
	err := j.db.Model(&entities.Jobs{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
		"title":       job.Title,
		"description": job.Description,
		"open_date":   job.OpenDate,
		"close_date":  job.CloseDate,
		"company_id":  job.CompanyId,
	}).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (j jobRepository) DeleteJob(id int) error {
	err := j.db.Where("id = ?", id).Delete(&entities.Jobs{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (j jobRepository) GetJobLists(request *requests.JobListsFilter) (interface{}, error) {
	jobsInterface := []map[string]interface{}{}

	query := j.db.Table("jobs").
		Select("jobs.id, title, count(applications.id) as application_count, description, open_date, close_date, users.name as company").
		Joins("JOIN users ON jobs.company_id = users.id").
		Joins("LEFT JOIN applications ON applications.job_id = jobs.id").
		Group("jobs.id").
		Limit(request.Size).Offset((request.Page - 1) * request.Size)

	var err error
	if request.Title == "" {
		err = query.Find(&jobsInterface).Error
	} else {
		err = query.Where("title LIKE ?", "%"+request.Title+"%").Find(&jobsInterface).Error
	}

	if err != nil {
		return nil, err
	}

	return jobsInterface, nil
}

func (j jobRepository) GetCompanyJobLists(companyId int, request *requests.JobListsFilter) (interface{}, error) {
	jobsInterface := []map[string]interface{}{}

	query := j.db.Table("jobs").
		Select("jobs.id, title, count(applications.id) as application_count, description, open_date, close_date, users.name as company").
		Joins("JOIN users ON jobs.company_id = users.id").
		Joins("LEFT JOIN applications ON applications.job_id = jobs.id").
		Where("jobs.company_id = ?", companyId).
		Group("jobs.id").
		Limit(request.Size).Offset((request.Page - 1) * request.Size)

	var err error
	if request.Title == "" {
		err = query.Find(&jobsInterface).Error
	} else {
		err = query.Where("title LIKE ?", "%"+request.Title+"%").Find(&jobsInterface).Error
	}

	if err != nil {
		return nil, err
	}

	return jobsInterface, nil
}

func (j jobRepository) GetJobDetail(id int) (*entities.Jobs, error) {
	var jobResponse *entities.Jobs
	err := j.db.Preload("Company").Where("id = ?", id).First(&jobResponse).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Job not found")
		}
		return nil, err
	}

	return jobResponse, err
}
