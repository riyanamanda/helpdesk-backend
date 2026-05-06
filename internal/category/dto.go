package category

import "time"

type CategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListCategoriesParams struct {
	Limit  int
	Offset int
}

type ListCategoriesResult struct {
	Data  []Category
	Total int
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,max=50"`
}
