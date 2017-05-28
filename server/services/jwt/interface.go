package jwt

type JWTServiceProvider interface {
	GenerateAuthToken(additionalClaims map[string]string, expiryInSeconds int) (authToken string, err error)
	ParseToken(token string) (claims map[string]interface{}, err error)
	IsAuthTokenValid(authToken string) (valid bool)
}
