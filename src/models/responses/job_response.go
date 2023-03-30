package responses

import "time"

type JobResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// 1. JOB LIST
type JobListsResponse struct {
	Data      interface{} `json:"data"`
	Status    string      `json:"status"`
	TotalData int         `json:"total_data"`
}

// 2. CANDIDATE JOB DETAIL
type JobDetailCandidateResponse struct {
	Data JobDetailCandidate `json:"data"`
	JobResponse
}
type JobDetailCandidate struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ApplicantCount int       `json:"applicant_count"`
	OpenDate       time.Time `json:"open_date"`
	CloseDate      time.Time `json:"close_date"`
	Company        string    `json:"company"`
}

// 3. COMPANY JOB DETAIL
type JobDetailCompanyResponse struct {
	Data JobDetailCompany `json:"data"`
	JobResponse
}

type JobDetailCompany struct {
	ID             int         `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	ApplicantCount int         `json:"applicant_count"`
	OpenDate       time.Time   `json:"open_date"`
	CloseDate      time.Time   `json:"close_date"`
	Company        string      `json:"company"`
	Applications   interface{} `json:"applications"`
}
