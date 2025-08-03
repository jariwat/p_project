package models

type Paginator struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	TotalRows  int `json:"total_rows"`
}

func NewPaginator(page, perPage int) *Paginator {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	return &Paginator{
		Page:    page,
		PerPage: perPage,
	}
}

func (p *Paginator) SetTotal(totalRows int) {
	if totalRows < 1 {
		p.TotalPages = 0
	} else {
		p.TotalRows = totalRows
		p.TotalPages = (totalRows + p.PerPage - 1) / p.PerPage // Ceiling division
	}
}