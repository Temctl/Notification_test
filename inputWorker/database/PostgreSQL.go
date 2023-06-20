package database

import (
	"fmt"

	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
)

func ConnectPostgreSQL() {
	db, err := connections.ConnectPostgreSQL()
	if err != nil {
		elog.Error().Panic(err)
	}
	// Perform database operations...
	// Execute a query
	rows, err := db.Query("select email, username from users")
	if err != nil {
		elog.Error().Panic(err)
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var email string
		var username string
		err := rows.Scan(&email, &username)
		if err != nil {
			elog.Error().Panic(err)
		}
		fmt.Println("ID:", email, "Name:", username)
	}
	db.Close()
}
