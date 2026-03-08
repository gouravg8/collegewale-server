package view

type DataList struct {
	TotalRecords int64 `json:"total_records"`
	Data         []any `json:"data"`
}
