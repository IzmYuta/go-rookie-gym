package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"practice/handler"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database?charset=utf8&parseTime=true")
	if err != nil {
		log.Printf("failed to open a db err = %s", err.Error())
		return
	}
	defer db.Close()
	if err := db.PingContext(context.Background()); err != nil {
		log.Printf("failed to ping err = %s", err.Error())
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/user", handler.UserHandler(db))
	http.HandleFunc("/groups", handler.GroupsHandler(db))
	http.HandleFunc("/group", handler.GroupHandler(db))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

