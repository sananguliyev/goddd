package config

import "os"

func GetDatabaseConfig() (string, string, string, string, string) {
	var host, port, user, password, database string

	if host = os.Getenv("DATABASE_HOST"); len(host) == 0 {
		host = "localhost"
	}
	if port = os.Getenv("DATABASE_PORT"); len(port) == 0 {
		port = "5432"
	}
	if user = os.Getenv("DATABASE_USER"); len(user) == 0 {
		user = "goddd"
	}
	if password = os.Getenv("DATABASE_PASSWORD"); len(password) == 0 {
		password = "goddd"
	}
	if database = os.Getenv("DATABASE_DB"); len(database) == 0 {
		database = "goddd"
	}

	return host, port, user, password, database
}
