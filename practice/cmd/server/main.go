package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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
	http.HandleFunc("/user", userHandler(db))
	http.HandleFunc("/groups", groupsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type UserInput struct {
	Name string `json:"name"`
}

type UserOutput struct {
	ID      int `json:"id"`
	GroupID int `json:"group_id"`
}

func userHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POSTのみ許可
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// JSONリクエストボディから構造体を作成する
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
		// Userの作成
		id, err := db.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?);", input.Name)
		if err != nil {
			log.Printf("failed to exec query err = %s", err.Error())
			return
		}
		last_user_id, err := id.LastInsertId()
		if err != nil {
			log.Printf("failed to get a last insert id err = %s", err.Error())
			return
		}
		// groupの作成　groupsはMYSQLの予約語なのでバッククオートで囲む必要がある
		group_id, err := db.ExecContext(context.Background(), "INSERT INTO `groups` (user_id, name) VALUES (?, ?);", last_user_id, input.Name)
		if err != nil {
			log.Printf("failed to exec query err = %s", err.Error())
			return
		}
		last_group_id, err := group_id.LastInsertId()
		if err != nil {
			log.Printf("failed to get a last insert id err = %s", err.Error())
			return
		}
		// JSONレスポンスの作成
		output.ID = int(last_user_id)
		output.GroupID = int(last_group_id)
		j, err := json.Marshal(&output)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	// GETのみ許可
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// クエリパラメータの取得
	user_id:= r.URL.Query().Get("user_id")
	fmt.Println(user_id)
}
