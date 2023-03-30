package responses

type ApplicationResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// 1. READ APPLICATIONS
type ReadApplicationResponse struct {
	Data interface{} `json:"data"`
	ApplicationResponse
	TotalData int `json:"total_data"`
}

// 2. APPLICATION DETAIL
type ApplicationDetail struct {
	Id        int         `json:"id"`
	JobTitle  string      `json:"job_title"`
	Processes interface{} `json:"processes"`
}

// Candidate
type ApplicationDetailCandidateResponse struct {
	ApplicationDetail
	Company string `json:"company"`
	Status  string `json:"status"`
}

// Company
type ApplicationDetailCompanyResponse struct {
	ApplicationDetail
	Candidate string `json:"candidate"`
	Status    string `json:"status"`
}
