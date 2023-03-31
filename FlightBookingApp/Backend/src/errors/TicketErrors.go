package errors

type NotEnoughVacantSeatsError struct{}

func (error NotEnoughVacantSeatsError) Error() string {
	return "Not enough vacant seats."
}

type FlightIsCanceledError struct{}

func (error FlightIsCanceledError) Error() string {
	return "Can not but tickets for canceled flight."
}
