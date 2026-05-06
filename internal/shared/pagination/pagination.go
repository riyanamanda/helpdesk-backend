package pagination

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type Params struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (p *Params) Normalize() (page, limit, offset int) {
	page = p.Page
	if page <= 0 {
		page = DefaultPage
	}

	limit = p.Limit
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset = (page - 1) * limit
	return
}
