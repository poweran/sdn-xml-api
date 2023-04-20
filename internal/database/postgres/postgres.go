package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sdn-xml-api/internal/database/repository"
)

// DB Struct to represent a PostgreSQL database connection
type DB struct {
	*sql.DB
}

// NewDB Function to create a new PostgreSQL database connection
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}

	return &DB{db}, nil
}

// InsertIndividual Function to insert a new individual into the sdn_individuals table
func InsertIndividual(db *DB, uid int, firstName string, lastName string) error {
	_, err := db.Exec(`
		INSERT INTO people (uid, first_name, last_name, sdn_type)
		VALUES ($1, $2, $3, $4);
	`, uid, firstName, lastName, "")
	if err != nil {
		return fmt.Errorf("failed to insert individual into sdn_individuals table: %v", err)
	}

	return nil
}

// GetIndividuals Function to retrieve individuals from the sdn_individuals table by name
func GetIndividuals(db *DB, name string, matchType string) ([]repository.Person, error) {
	// Build the SQL query based on the match type
	var query string
	if matchType == "strong" {
		query = `
			SELECT uid, first_name, last_name
			FROM people
			WHERE CONCAT(first_name, ' ', last_name) = $1;
		`
	} else if matchType == "weak" {
		query = `
			SELECT uid, first_name, last_name
			FROM people
			WHERE first_name ILIKE $1 OR last_name ILIKE $1;
		`
	} else {
		return nil, fmt.Errorf("invalid match type: %s", matchType)
	}

	// Execute the query and build the result set
	rows, err := db.Query(query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve individuals from sdn_individuals table: %v", err)
	}
	defer rows.Close()

	var individuals []repository.Person
	for rows.Next() {
		Person := repository.Person{}
		err := rows.Scan(&Person.UID, &Person.FirstName, &Person.LastName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan Person from people table: %v", err)
		}
		individuals = append(individuals, Person)
	}

	return individuals, nil
}
