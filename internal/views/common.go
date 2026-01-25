package views

type PageFilter struct {
	AllPages bool `json:"all_pages"`
	PageNum  int  `json:"page_num"`
	PageSize int  `json:"page_size"`
}

type SortFilter struct {
	SortField string `json:"sort_field"`
	SortOrder string `json:"sort_order"`
}

type StatusFilter struct {
	Status string `json:"status"`
}

type DateFilter struct {
	Date string `json:"date"`
}

type DateRangeFilter struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Response struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ListResponse struct {
	TotalRecords int `json:"total_records"`
	Data         any `json:"data,omitempty"`

	FailedCount  int64 `json:"failed_count,omitempty"`
	MessageCount int64 `json:"message_count,omitempty"`
}

type CountResult struct {
	FailedCount  int64
	MessageCount int64
}
