package interfaces

type CacheServiceProvider interface {
	Get(key string) (value string, err error)
	SetWithExpiry(key string, value string, expiryInSeconds int) (err error)
	Delete(key string) (err error)
	Exists(key string) (exists bool, err error)
	Flush() (err error)
}
