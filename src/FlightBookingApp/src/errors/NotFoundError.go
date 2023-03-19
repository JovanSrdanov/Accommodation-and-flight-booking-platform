package errors

type NotFoundError struct{}

func (error *NotFoundError) Error() string {
	return "Entity not found"
}
