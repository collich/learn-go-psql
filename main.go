package main

import (
	"database/sql"
	// "log"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Product struct{
	Name string
	Price float64
	Available bool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// pg_user := os.Getenv("PG_USER")
	// pg_pass := os.Getenv("PG_PASS")
	// pg_port := os.Getenv("PG_PORT")
	// pg_db_name := os.Getenv("PG_DB_NAME")
	pg_user, pg_pass, pg_port, pg_db_name := os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_PORT"), os.Getenv("PG_DB_NAME")

	variable := fmt.Sprintf("postgres://%v:%v@localhost:%v/%v?sslmode=disable", pg_user, pg_pass, pg_port, pg_db_name)
	connStr := variable
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	createProductTable(db)

	// product := Product{"LootBox - 1000x", 999.99, false}
	// prod_id := insertProduct(db, product)

	// fmt.Printf("ID = %d\n", prod_id)

	// name, price, available := getProduct(db, 111)
	// fmt.Printf("Name: %v\n", name)
	// fmt.Printf("Price: %v\n", price)
	// fmt.Printf("Availability: %v\n", available)

	data := getProducts(db)

	for i, s := range data {
		fmt.Printf("ID: %v\nName: %v\nPrice: %v\nAvailable: %v\n", i, s.Name, s.Price, s.Available)
	}
}

func createProductTable(db *sql.DB)  {
	query := `CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(6,2) NOT NULL,
		available BOOLEAN,
		created timestamp DEFAULT NOW()
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insertProduct(db *sql.DB, product Product) int {
	var product_id int

	query := `INSERT INTO product (name, price, available)
	VALUES ($1, $2, $3) RETURNING id`

	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&product_id)
	if err != nil {
		log.Fatal(err)
	}
	return product_id
}

func getProduct(db *sql.DB, id int) (string, float64, bool) {
	var name string
	var price float64
	var available bool
	
	query := `SELECT name, price, available FROM product WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&name, &price, &available)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found with ID %d", id)
		}
		log.Fatal(err)
	}
	return name, price, available
}

func getProducts(db *sql.DB) []Product {
	var data []Product
	var name string
	var price float64
	var available bool

	rows, err := db.Query("SELECT name, price, available FROM product")
	if err != nil {
		log.Fatal(err)
	}
	
	for rows.Next() {
		err := rows.Scan(&name, &price, &available)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Product{name, price, available})
	}

	defer rows.Close()
	return data
}