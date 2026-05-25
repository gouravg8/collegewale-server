package view

type DataList struct {
	Pagination PaginationConfig `json:"pagination"`
	Data       any              `json:"data"`
	Meta       interface{}      `json:"meta,omitempty"` // For additional context like filters applied, sorting, etc.
}

func NewDataList(data []any, totalRecords int64, page int, pageSize int) DataList {
	pagination := PaginationConfig{
		PageFilter: PageFilter{
			AllPages: true,
			PageNum:  page,
			PageSize: pageSize,
		},
		TotalRecords: totalRecords,
		HasNext:      false,
		HasPrev:      false,
	}

	hasPrev := page > 1 && totalRecords > 0

	lastRecordOnPage := int64((page-1)*pageSize + len(data))
	hasNext := lastRecordOnPage < totalRecords

	pagination.HasNext = hasNext
	pagination.HasPrev = hasPrev
	pagination.PageFilter.AllPages = int64(len(data)) == totalRecords

	return DataList{
		Pagination: pagination,
		Data:       data,
	}
}

func NewAllDataList(data []any) DataList {
	tRecords := len(data)
	pagination := PaginationConfig{
		PageFilter: PageFilter{
			AllPages: true,
			PageNum:  1,
			PageSize: tRecords,
		},
		TotalRecords: int64(tRecords),
		HasNext:      false,
		HasPrev:      false,
	}
	return DataList{
		Pagination: pagination,
		Data:       data,
	}
}
