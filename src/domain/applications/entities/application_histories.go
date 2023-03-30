package entities

import "time"

type ApplicationHistories struct {
	ID              int          `json:"id" gorm:"primaryKey"`
	Status          string       `json:"status" gorm:"type:enum('Applied', 'HR Interview', 'Client Interview', 'Passed', 'Rejected', 'Cancelled');default:'Applied';not null"`
	LastProcessDate time.Time    `json:"last_process_date" gorm:"type:datetime;not null"`
	Application     Applications `json:"applications" gorm:"foreignKey:ApplicationId"`
	ApplicationId   int          `json:"application_id" gorm:"type:int;not null"`
}
