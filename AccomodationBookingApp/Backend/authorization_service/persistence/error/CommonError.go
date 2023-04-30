package error

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "Entity not found"
}
