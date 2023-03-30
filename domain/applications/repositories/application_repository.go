package repositories

import (
	"errors"
	"job-portal-lite/domain/applications/entities"
	UserInterface "job-portal-lite/domain/user/interfaces"
	UserRepository "job-portal-lite/domain/user/repositories"
	"job-portal-lite/models/requests"
	"job-portal-lite/shared/databases"

	"gorm.io/gorm"
)

type applicationRepository struct {
	db             *gorm.DB
	userRepository UserInterface.UserRepository
}

func NewApplicationRepository() *applicationRepository {
	return &applicationRepository{
		db:             databases.Connect(),
		userRepository: UserRepository.NewUserRepository(),
	}
}

func (a applicationRepository) CreateApplication(applications *entities.Applications) (*entities.Applications, error) {
	err := a.db.Save(applications).Error
	if err != nil {
		return nil, err
	}

	return applications, err
}

func (a applicationRepository) ReadAllApplication(userId int, jobId int, request *requests.ApplicationListsFilter) (interface{}, error) {
	applicationsInterface := []map[string]interface{}{}
	user, err := a.userRepository.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	// Validasi kalo gk ada status, makan tampilkan semua
	// jangan lupa cek juga buat yang job detail karena manggil read all application juga
	var status string = "Applied"
	switch request.ApplicationStatusId {
	case 1:
		status = "Applied"
	case 2:
		status = "HR Interview"
	case 3:
		status = "Client Interview"
	case 4:
		status = "Passed"
	case 5:
		status = "Rejected"
	case 6:
		status = "Cancelled"
	}

	if user.Role == "Candidate" {
		err := a.db.Table("applications").
			Select("applications.id, jobs.title as job_title, application_histories.last_process_date, application_histories.status as last_status, users.name as company").
			Joins("JOIN application_histories ON application_histories.application_id = applications.id").
			Joins("JOIN jobs ON applications.job_id = jobs.id").
			Joins("JOIN users ON jobs.company_id = users.id").
			Where(a.db.Where("applications.candidate_id = ?", userId).Where("application_histories.status = ?", status).Where("0 = ?", jobId)).Where(
			a.db.Where("applications.job_id = ?", jobId).
				Where("application_histories.status = ?", "Applied").
				Where("application_histories.status = ?", status).Where("0 != ?", jobId),
		).
			Limit(request.Size).
			Offset((request.Page - 1) * request.Size).
			Find(&applicationsInterface).Error
		if err != nil {
			return nil, err
		}
	} else {
		query := a.db.Table("applications").
			Select("applications.id, jobs.title as job_title, application_histories.last_process_date, application_histories.status as last_status, users.name as candidate").
			Joins("JOIN application_histories ON application_histories.application_id = applications.id").
			Joins("JOIN jobs ON applications.job_id = jobs.id").
			Joins("JOIN users ON applications.candidate_id = users.id").
			Limit(request.Size).
			Offset((request.Page-1)*request.Size).
			Where("applications.company_id = ?", userId)
		if jobId == 0 {
			err = query.Find(&applicationsInterface).Error
		} else {
			err = query.Where("applications.job_id = ?", jobId).
				Where("application_histories.status = ?", "Applied").
				Find(&applicationsInterface).Error
		}
	}

	if err != nil {
		return nil, err
	}

	return applicationsInterface, err
}

func (a applicationRepository) ApplicationDetail(role string, id int) (*entities.Applications, error) {
	preloadRole := "Candidate"
	if role == "Candidate" {
		preloadRole = "Company"
	}

	var aplication *entities.Applications
	err := a.db.Preload("Job").Preload(preloadRole).Where("id = ?", id).First(&aplication).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Application")
		}
		return nil, err
	}

	return aplication, err
}

func (a applicationRepository) ProceedApplication(applicationHistory *entities.ApplicationHistories) (*entities.ApplicationHistories, error) {
	err := a.db.Save(applicationHistory).Error
	if err != nil {
		return nil, err
	}

	return applicationHistory, err
}

func (a applicationRepository) CreateApplicationHistory(applicationHistory *entities.ApplicationHistories) (*entities.ApplicationHistories, error) {
	err := a.db.Save(applicationHistory).Error
	if err != nil {
		return nil, err
	}

	return applicationHistory, err
}

func (a applicationRepository) ShowJobApplicationHistory(id int) (interface{}, error) {
	histories := []map[string]interface{}{}
	err := a.db.Table("application_histories").
		Select("application_histories.status, application_histories.last_process_date as date").
		Joins("JOIN applications ON application_histories.application_id = applications.id").
		Where("application_histories.application_id = ?", id).
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}
