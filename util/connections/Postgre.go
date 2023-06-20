package connections

import (
	"database/sql"

	"github.com/Temctl/E-Notification/util/elog"
	_ "github.com/lib/pq"
)

func ConnectPostgreSQL() (*sql.DB, error) {
	// Define the database connection string
	db, err := sql.Open("postgres", "host=172.72.0.11 user=postgres password=changeme dbname=postgres sslmode=disable")
	if err != nil {
		elog.Error().Println(err)
		return nil, err
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		elog.Error().Println(err)
		return nil, err
	}
	return db, nil
}
