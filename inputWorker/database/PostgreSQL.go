package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectPostgreSQL() {
	fmt.Println("Dsfdsf")
	// Define the database connection string
	db, err := sql.Open("postgres", "host=172.72.0.11 user=postgres password=changeme dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Println("SDfsdf")
		log.Fatal(err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		fmt.Println("123")
		log.Fatal(err)
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
}
