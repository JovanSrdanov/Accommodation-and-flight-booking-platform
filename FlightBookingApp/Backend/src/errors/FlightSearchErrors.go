package errors

type SearchDateMustBeInFutureError struct{}

func (error SearchDateMustBeInFutureError) Error() string {
	return "Search date must be in future."
}

type DesiredNumberOfSeatsMustBeGreaterThanZeroError struct{}

func (error DesiredNumberOfSeatsMustBeGreaterThanZeroError) Error() string {
	return "Desired number of seats must be greater than zero."
}
