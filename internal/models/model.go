package models

type Response[T any] struct {
	Message    string              `json:"message,omitempty"`
	Data       T                   `json:"data,omitempty"`
	Pagination *PaginationMetaData `json:"pagination,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

type PaginationMetaData struct {
	CurrentPage int    `json:"current_page"`
	Previous    string `json:"previous"`
	Next        string `json:"next"`
	TotalPage   int    `json:"total_page"`
	TotalItem   int    `json:"total_item"`
}
