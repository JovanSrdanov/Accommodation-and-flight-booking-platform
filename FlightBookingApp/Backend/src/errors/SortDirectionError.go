package errors

type SortDirectionError struct{}

func (error SortDirectionError) Error() string {
	return "Sort direction must be asc, dsc or no_sort"
}
