package errors

type BasicResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type FormErrors struct {
	Errors interface{} `json:"errors"`
	BasicResponse
}
