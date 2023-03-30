package entities

import (
	User "job-portal-lite/domain/user/entities"
	"time"
)

type Jobs struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	OpenDate    time.Time `json:"open_date" gorm:"type:datetime;not null"`
	CloseDate   time.Time `json:"close_date" gorm:"type:datetime;not null"`
	Company     User.User `json:"company" gorm:"foreignKey:CompanyId"`
	CompanyId   int       `json:"company_id" gorm:"type:int;not null"`
}
