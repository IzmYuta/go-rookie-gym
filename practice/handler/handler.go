package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	groupuc "practice/usecase/group"
	useruc "practice/usecase/user"
)

type UserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID      int64 `json:"id"`
	GroupID int64 `json:"group_id"`
}

func UserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POSTのみ許可
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// JSONリクエストボディから構造体を作成する
		var req UserRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// usecaseから取得
		output, err := useruc.User(db, &useruc.UserInput{
			Name: req.Name,
		})
		if err != nil {
			log.Printf("failed to put user err = %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// JSONレスポンスの作成
		var res UserResponse
		res.ID = output.UserID
		res.GroupID = output.GroupID
		j, err := json.Marshal(&res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

type Group struct {
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
		// クエリパラメータをint型に変換する
		// パラメータが空、型が違うときは BadRequestを返す
		user_id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// usecaseから取得
		output, err := groupuc.Groups(db, user_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// JSONレスポンスの作成
		res := make([]Group, len(output))
		for i, g := range output {
			res[i] = Group{
				ID:   g.ID,
				Name: g.Name,
			}
		}
		j, err := json.Marshal(&res)
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

// DBから取得した値はint64型
type GroupResponse struct {
	ID int64 `json:"id"`
}

func GroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POSTのみ許可
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// JSONから構造体を作成
		var req GroupRequest
		var res GroupResponse
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Groupの作成
		group_id, err := groupuc.Group(db, &groupuc.GroupInput{
			UserID:    req.UserID,
			GroupName: req.Name,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// JSONレスポンスに変換
		res.ID = group_id
		j, err := json.Marshal(&res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}
