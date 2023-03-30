package errors

type JwtError struct{}

func (error JwtError) Error() string {
	return "Error with jwt"
}
