package errors

type SameStartPointAndDestinationError struct{}

func (error SameStartPointAndDestinationError) Error() string {
	return "Same start point and destination"
}

type FlightPassedError struct{}

func (error FlightPassedError) Error() string {
	return "Flight already passed"
}

type DesiredNumberOfSeatsMustBeGreaterThanZeroError struct{}

func (error DesiredNumberOfSeatsMustBeGreaterThanZeroError) Error() string {
	return "Desired  number of seats must be greater than zero."
}
