package errors

type SearchDateMustBeTodayOrInFutureError struct{}

func (error SearchDateMustBeTodayOrInFutureError) Error() string {
	return "Search date must be today or in future."
}

type DesiredNumberOfSeatsMustBeGreaterThanZeroError struct{}

func (error DesiredNumberOfSeatsMustBeGreaterThanZeroError) Error() string {
	return "Desired number of seats must be greater than zero."
}
