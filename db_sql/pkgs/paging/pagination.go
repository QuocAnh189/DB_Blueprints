package paging

import "math"

const DefaultPageSize int64 = 10

type Pagination struct {
	Page        int64 `json:"page"`
	Size        int64 `json:"size"`
	TakeAll     bool  `json:"take_all"`
	TotalCount  int64 `json:"total_count"`
	TotalPages  int64 `json:"total_pages"`
	HasPrevious bool  `json:"has_previous"`
	HasNext     bool  `json:"has_next"`
}

func NewPagination(page, size, total int64) *Pagination {
	if size <= 0 {
		size = DefaultPageSize
	}
	totalPages := int64(math.Ceil(float64(total) / float64(size)))
	if page < 1 {
		page = 1
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}
	if totalPages == 0 {
		page = 1
	}

	p := &Pagination{
		Page: page, Size: size, TotalCount: total, TotalPages: totalPages,
		HasPrevious: page > 1, HasNext: page < totalPages,
	}
	return p
}
