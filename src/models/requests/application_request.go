package requests

type CreateApplicationRequest struct {
	JobId int `json:"job_id" validate:"required"`
}

type ProceedApplicatioRequest struct {
	StatusId int `json:"status_id" validate:"required,numeric,oneof=1 2 3 4"`
}

type ApplicationListsFilter struct {
	Page                int    `json:"page" form:"page"`
	Size                int    `json:"size" form:"size"`
	Title               string `json:"title" form:"title"`
	ApplicationStatusId int    `json:"application_status_id" form:"application_status_id"`
}
