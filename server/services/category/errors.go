package category

var _ error = new(CategoryNotFoundError)

type CategoryNotFoundError struct {
}

func NewCategoryNotFoundError() error {
	return CategoryNotFoundError{}
}

func (c CategoryNotFoundError) Error() string {
	return "category not found"
}
