package database

import (
	"fmt"
	"log"

	"github.com/Temctl/E-Notification/util/connections"
	_ "github.com/lib/pq"
)

func ConnectPostgreSQL() {
	db, err := connections.ConnectPostgreSQL()
	if err != nil {
		fmt.Println("Dasd")
	}
	// Perform database operations...
	// Execute a query
	rows, err := db.Query("select email, username from users")
	if err != nil {
		fmt.Println("zxc")
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var email string
		var username string
		err := rows.Scan(&email, &username)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		fmt.Println("ID:", email, "Name:", username)
	}
	db.Close()
}
