package errors

type NotEnoughVacantSeats struct{}

func (error NotEnoughVacantSeats) Error() string {
	return "Not enough vacant seats."
}
