package requests

import "time"

type UpsertJob struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	OpenDate    time.Time `json:"open_date" validate:"required"`
	CloseDate   time.Time `json:"close_date" validate:"required"`
}

type JobListsFilter struct {
	Page  int    `json:"page" form:"page"`
	Size  int    `json:"size" form:"size"`
	Title string `json:"title" form:"title"`
}
