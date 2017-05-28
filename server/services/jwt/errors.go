package jwt

type JWTInvalidError struct {
}

func NewJWTInvalidError() error {
	return JWTInvalidError{}
}

func (p JWTInvalidError) Error() string {
	return ""
}
