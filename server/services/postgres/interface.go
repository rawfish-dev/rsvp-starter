package postgres

import (
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/interfaces"
)

type PostgresServiceProvider interface {
	interfaces.Storage
}
