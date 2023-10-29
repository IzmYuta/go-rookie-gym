package handler

import(
	"net/http"
	"context"
	"io"
	"log"
	"encoding/json"
	"database/sql"
)

type UserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID      int `json:"id"`
	GroupID int `json:"group_id"`
}

func UserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POSTのみ許可
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// JSONリクエストボディから構造体を作成する
		var input UserRequest
		var output UserResponse
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

type GroupsResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GroupsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// GETのみ許可
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// こう定義しないと中身が空の時のレスポンスがnullになる
		groups := make([]GroupsResponse, 0)
		// クエリパラメータの取得
		user_id := r.URL.Query().Get("user_id")
		if user_id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// DBから取得
		// 取得できなかったときは空配列を返す
		rows, err := db.QueryContext(context.Background(), "SELECT id, name FROM `groups` WHERE user_id = ?", user_id)
		if err != nil {
			j, err := json.Marshal(&groups)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(j)
			return
		}
		defer rows.Close()
		// rows.Next()はbooleanを返す。次の要素があるときはTrueになる。
		// つまり、rowsの要素が全て取り出されるまで無限ループする
		for rows.Next() {
			var group GroupsResponse
			if err := rows.Scan(&group.ID, &group.Name); err != nil {
				return
			}
			groups = append(groups, group)
		}
		j, err := json.Marshal(&groups)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

type GroupRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type GroupResponse struct {
	ID int `json:"id"`
}

func GroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POSTのみ許可
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// JSONから構造体を作成
		var input GroupRequest
		var output GroupResponse
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Groupの作成
		id, err := db.ExecContext(context.Background(), "INSERT INTO `groups` (user_id, name) VALUES (?,?);", input.UserID, input.Name)
		if err != nil {
			log.Printf("failed to exec query err = %s", err.Error())
			return
		}
		last_id, err := id.LastInsertId()
		if err != nil {
			log.Printf("failed to get a last insert id err = %s", err.Error())
			return
		}
		// JSONレスポンスに変換
		output.ID = int(last_id)
		j, err := json.Marshal(&output)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}
