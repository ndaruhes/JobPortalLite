package entities

import (
	Jobs "job-portal-lite/domain/jobs/entities"
	User "job-portal-lite/domain/user/entities"
)

type Applications struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Candidate   User.User `json:"candidate" gorm:"foreignKey:CandidateId"`
	CandidateId int       `json:"candidate_id" gorm:"type:int;not null"`
	Job         Jobs.Jobs `json:"job" gorm:"foreignKey:JobId"`
	JobId       int       `json:"job_id" gorm:"type:int;not null"`
	Company     User.User `json:"company" gorm:"foreignKey:CompanyId"`
	CompanyId   int       `json:"company_id" gorm:"type:int;not null"`
	IsDeleted   string    `json:"is_deleted" gorm:"type:enum('True', 'False');default:'False';not null"`
}
