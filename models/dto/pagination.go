package dto

type Pagination struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

type PaginationParam struct {
	Current  int `query:"current"`
	PageSize int `query:"page_size" validate:"max=128"`
}

func (p *PaginationParam) GetCurrent() int {
	return p.Current
}

func (p *PaginationParam) GetPageSize() int {
	pageSize := p.PageSize
	if pageSize == 0 {
		pageSize = 15
	}

	return pageSize
}
