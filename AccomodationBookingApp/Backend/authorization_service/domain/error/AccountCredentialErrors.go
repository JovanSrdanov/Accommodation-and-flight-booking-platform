package error

type FooError struct{}

func (f FooError) Error() string {
	return "Foo error message"
}
