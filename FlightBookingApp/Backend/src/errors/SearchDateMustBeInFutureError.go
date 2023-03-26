package errors

type SearchDateMustBeInFutureError struct{}

func (error SearchDateMustBeInFutureError) Error() string {
	return "Search date must be in future."
}
