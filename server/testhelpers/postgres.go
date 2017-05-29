package testhelpers

import (
	_ "github.com/lib/pq"

	"database/sql"
	"log"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"

	"gopkg.in/gorp.v1"
)

func NewTestPostgresService() postgres.PostgresServiceProvider {
	testPostgresConfig := config.LoadConfig().Postgres
	testBaseService := NewTestBaseService()

	return postgres.NewService(testBaseService, testPostgresConfig)
}

func TruncateTestPostgresDB() {
	testPostgresConfig := config.LoadConfig().Postgres

	dbConnection, err := sql.Open("postgres", testPostgresConfig.URL)
	if err != nil {
		log.Fatalf("postgres test service - unable to open connection to postgres due to %v", err.Error())
	}
	defer dbConnection.Close()

	dbConnection.SetMaxIdleConns(1)
	dbConnection.SetMaxOpenConns(1)

	gorpDB := &gorp.DbMap{Db: dbConnection, Dialect: gorp.PostgresDialect{}}

	_, err = gorpDB.Exec("TRUNCATE rsvps,invitations,categories")
	if err != nil {
		log.Fatalf("postgres test service - unable to truncate tables due to %v", err)
	}
}
