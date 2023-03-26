package errors

type ResultsPerPageMustBeGreaterThanZeroError struct{}

func (error ResultsPerPageMustBeGreaterThanZeroError) Error() string {
	return "Results per page must be greater than zero."
}
