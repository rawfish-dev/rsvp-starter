package postgres

import (
	_ "github.com/lib/pq"

	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"

	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"
)

var _ interfaces.CategoryStorage = new(service)
var _ interfaces.InvitationStorage = new(service)
var _ interfaces.RSVPStorage = new(service)

type service struct {
	ctx    context.Context
	gorpDB *gorp.DbMap
}

var singletonService *service
var once sync.Once

func NewService(ctx context.Context, postgresConfig config.PostgresConfig) *service {
	once.Do(func() {
		ctxLogger := ctx.Value("logger").(interfaces.Logger)

		dbConnection, err := sql.Open("postgres", postgresConfig.URL)
		if err != nil {
			ctxLogger.Fatalf("postgres service - unable to open connection to postgres due to %v", err.Error())
		}

		dbConnection.SetMaxIdleConns(postgresConfig.MaxIdle)
		dbConnection.SetMaxOpenConns(postgresConfig.MaxConnections)

		gorpDB := &gorp.DbMap{Db: dbConnection, Dialect: gorp.PostgresDialect{}}
		gorpDB.AddTableWithName(category{}, "categories").SetKeys(true, "ID")
		gorpDB.AddTableWithName(invitation{}, "invitations").SetKeys(true, "ID")
		gorpDB.AddTableWithName(rsvp{}, "rsvps").SetKeys(true, "ID")

		gorpDB.TypeConverter = dbTypeConverter{}

		singletonService = &service{ctx, gorpDB}
	})

	return singletonService
}

func (s *service) Close() error {
	return s.gorpDB.Db.Close()
}

type dbTypeConverter struct{}

func (c dbTypeConverter) ToDb(val interface{}) (interface{}, error) {

	switch t := val.(type) {

	case []string, []int64:
		b, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	return val, nil
}

func (c dbTypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {

	switch target.(type) {

	case *[]string, *[]int64:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return fmt.Errorf("unable to convert %v from db", reflect.TypeOf(target))
			}
			b := []byte(*s)
			return json.Unmarshal(b, target)
		}
		return gorp.CustomScanner{new(string), target, binder}, true
	}

	return gorp.CustomScanner{}, false
}
