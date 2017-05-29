package postgres

import (
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
)

type PostgresServiceProvider interface {
	interfaces.Storage
}
