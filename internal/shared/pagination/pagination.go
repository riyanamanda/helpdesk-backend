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

func (p *Params) Normalize() {
	if p.Page <= 0 {
		p.Page = DefaultPage
	}
	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
}
