package models

type GenericQueryOptions struct {
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}

type PaginationResponse struct {
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
}
