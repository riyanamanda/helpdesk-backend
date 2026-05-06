package category

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrInvalidCategory       = errors.New("invalid category")
)
