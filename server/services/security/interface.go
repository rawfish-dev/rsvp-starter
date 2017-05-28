package security

type SecurityServiceProvider interface {
	ValidateCredentials(username, password string) (valid bool)
	VerifyReCAPTCHA(token string) (valid bool)
}
