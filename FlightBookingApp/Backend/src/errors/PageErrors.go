package errors

type PageNumberMustBeGreaterThanZeroError struct{}

func (error PageNumberMustBeGreaterThanZeroError) Error() string {
	return "Page number must be greater than zero."
}

type ResultsPerPageMustBeGreaterThanZeroError struct{}

func (error ResultsPerPageMustBeGreaterThanZeroError) Error() string {
	return "Results per page must be greater than zero."
}

type SortDirectionError struct{}

func (error SortDirectionError) Error() string {
	return "Sort direction must be asc, dsc or no_sort"
}

type InvalidSortTypeError struct{}

func (error InvalidSortTypeError) Error() string {
	return "Sorting by this type is not valid."
}
