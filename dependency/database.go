package dependency

import (
	"errors"
	"fmt"
	"github.com/SananGuliyev/goddd/config"
	"github.com/go-pg/pg/v9"
)

func NewPostgresConnection() (*pg.DB, error) {
	host, port, user, password, database := config.GetDatabaseConfig()

	connection := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		User:     user,
		Password: password,
		Database: database,
	})

	_, err := connection.Exec("SELECT 1")
	if err != nil {
		return nil, errors.New("service is not available")
	}

	return connection, nil
}

func Close(db interface{}) {
	switch connection := db.(type) {
	case *pg.DB:
		connection.Close()
	}
}
