package session

type SessionServiceProvider interface {
	CreateWithExpiry(username string) (authToken string, err error)
	IsSessionValid(authToken string) (valid bool, err error)
	Destroy(authToken string) (err error)
}
