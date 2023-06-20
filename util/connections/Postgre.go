package connections

import (
	"database/sql"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/elog"
	_ "github.com/lib/pq"
)

func ConnectPostgreSQL() (*sql.DB, error) {
	// Define the database connection string
	url := "host=" + util.DB_HOST + " port=" + util.DB_HOST + " user=" + util.DB_USERNAME + " password=" + util.DB_PASSWORD + " dbname=" + util.DB_DBNAME + " sslmode=disable"
	db, err := sql.Open("postgres", url)
	if err != nil {
		elog.Error().Println(err)
		return nil, err
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		db.Close()
		elog.Error().Println(err)
		return nil, err
	}
	return db, nil
}
