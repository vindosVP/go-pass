// Package db consists the database utils
package db

import "fmt"

// PostgresDSN creates the postgres database dsn
func PostgresDSN(host string, port int, user string, password string, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
}
