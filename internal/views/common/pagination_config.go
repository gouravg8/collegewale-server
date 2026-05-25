package view

// PaginationConfig defines pagination parameters
type PaginationConfig struct {
	PageFilter
	TotalRecords int64 `json:"total_records"`
	HasNext      bool  `json:"has_next,omitempty"`
	HasPrev      bool  `json:"has_prev,omitempty"`
}
