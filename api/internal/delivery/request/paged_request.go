package request

type PagedRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func (r *PagedRequest) Normalise() {
	if r.Page <= 0 {
		r.Page = DefaultPage
	}
	if r.PageSize <= 0 || r.PageSize > MaxPageSize {
		r.PageSize = DefaultPageSize
	}
}

func (r *PagedRequest) Offset() int {
	return (r.Page - 1) * r.PageSize
}
