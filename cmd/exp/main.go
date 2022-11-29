package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)
}

func main() {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "lenslocked",
		Password: "lenslocked",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to PSQL")

	// Create a table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR (255) UNIQUE NOT NULL,
		name TEXT
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		amount INT,
		description TEXT
	);
	`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tables Created")

	// name := "Danny Festor"
	// email := "denakino@gmail.com"
	// // _, err = db.Exec(`
	// // INSERT INTO users(name, email)
	// // VALUES ($1, $2);
	// // `, name, email)
	// row := db.QueryRow(`
	// INSERT INTO users(name, email)
	// VALUES ($1, $2)
	// RETURNING id;
	// `, name, email)
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("User created with ID", id)

	// query := `
	// SELECT name, email
	// FROM users
	// WHERE id=$1;
	// `
	// id := 1
	// row := db.QueryRow(query, id)
	// var name, email string
	// err = row.Scan(&name, &email)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		fmt.Println("Error, no rows")
	// 		return
	// 	}
	// 	panic(err)
	// }
	// fmt.Printf("User information: name=%s email=%s\n", name, email)

	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Fake order #%d", i)
	// 	_, err := db.Exec(`
	// 	INSERT INTO orders (user_id, amount, description)
	// 	VALUES ($1, $2, $3);
	// 	`, id, amount, desc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	userID := 1

	type Order struct {
		ID          int
		UserID      int
		Amount      int
		Description string
	}

	var orders []Order
	query := `
	SELECT id, amount, description
	FROM orders
	WHERE user_id = $1
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		order.UserID = userID
		err = rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	fmt.Println("Orders:", orders)
}
