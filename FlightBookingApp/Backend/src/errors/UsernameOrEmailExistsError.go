package errors

type UsernameOrEmailExistsError struct{}

func (error *UsernameOrEmailExistsError) Error() string {
	return "An account with the same username and/or password already exists."
}
