package errors

type DesiredNumberOfSeatsMustBeGreaterThanZeroError struct{}

func (error DesiredNumberOfSeatsMustBeGreaterThanZeroError) Error() string {
	return "Desired  number of seats must be greater than zero."
}
