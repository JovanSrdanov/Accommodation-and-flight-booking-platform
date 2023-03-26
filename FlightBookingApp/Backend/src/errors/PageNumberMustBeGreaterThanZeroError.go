package errors

type PageNumberMustBeGreaterThanZeroError struct{}

func (error PageNumberMustBeGreaterThanZeroError) Error() string {
	return "Page number must be greater than zero."
}
