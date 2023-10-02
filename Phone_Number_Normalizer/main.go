package main

import (
	"database/sql"
	"fmt"
	"log"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

func normalize(phone string) string {
	var runes []rune
	for _, c := range phone {
		if unicode.IsDigit(c) {
			runes = append(runes, c)
		}
	}
	return string(runes)
}

func createTable(db *sql.DB) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS phone (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			phone TEXT NOT NULL
		);
	`
	_, err := db.Exec(createTableSQL)
	return err
}
func insertPhone(db *sql.DB, p string) error {
	insert := `
		INSERT INTO phone (phone) VALUES (?)
	`
	_, err := db.Exec(insert, normalize(p))
	return err
}

func getPhoneById(db *sql.DB, id int) (string, error) {
	var phone string
	err := db.QueryRow(`SELECT phone FROM phone WHERE id = ?`, id).Scan(&phone)
	if err != nil {
		return "", err
	}
	return phone, nil
}

func deletePhoneById(db *sql.DB, id int) error {
	deleteSQL := "DELETE FROM phone WHERE id = ?"
	_, err := db.Exec(deleteSQL, id)
	return err
}

func main() {
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}

	err = insertPhone(db, "(123) 456-7893")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone number inserted successfully.")

	phone, err := getPhoneById(db, 1)
	fmt.Println(phone)

	err = deletePhoneById(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone number deleted successfully.")

	return
}
