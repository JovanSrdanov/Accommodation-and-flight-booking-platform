package errors

type InvalidSortTypeError struct{}

func (error InvalidSortTypeError) Error() string {
	return "Sorting by this type is not valid."
}
