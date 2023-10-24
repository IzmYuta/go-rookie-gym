package main

import (
	"fmt"
	"log"
	"io"
	"context"
	"encoding/json"
	"net/http"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct{
	ID int
	Name string
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database?charset=utf8&parseTime=true")
	if err != nil {
		log.Printf("failed to open a db err = %s", err.Error())
		return
	}

	if err := db.PingContext(context.Background()); err != nil {
		log.Printf("failed to ping err = %s", err.Error())
		return
	}
	id, err := db.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?);", "sample User")
	if err != nil {
		log.Printf("failed to exec query err = %s", err.Error())
		return
	}

	lastid, err := id.LastInsertId()
	if err != nil {
		log.Printf("failed to get a last insert id err = %s", err.Error())
		return
	}
	var s User

	if err := db.QueryRowContext(context.Background(), "SELECT id, name FROM users WHERE id = ?", lastid).Scan(&s.ID, &s.Name); err != nil {
		log.Printf("failed to scan err = %s", err.Error())
		return
	}

	log.Printf("User is %#v", s)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/user",userHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type UserInput struct {
	Name string `json:"name"`
}

type UserOutput struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

func userHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var input UserInput
	var output UserOutput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output.ID = 1
	output.Name = input.Name
	j, err := json.Marshal(&output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}