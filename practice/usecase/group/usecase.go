package group

import (
	"context"
	"database/sql"
	"log"
	domain "practice/domain/group"
)

// 構造体を扱うときは、ポインタで操作すると何かとメリットがある
type Usecase interface {
	Group(db *sql.DB, input *GroupInput) (int64, error)
	Groups(db *sql.DB, id int) ([]*domain.Group, error)
}
type GroupInput struct {
	UserID    int
	GroupName string
}

func Groups(db *sql.DB, id int) ([]*domain.Group, error) {
	outputs := make([]*domain.Group, 0)
	// DB処理
	rows, err := db.QueryContext(context.Background(), "SELECT id, name FROM `groups` WHERE user_id = ?", id)
	if err != nil {
		// 取得件数が0だった場合空配列を返す
		return outputs, err
	}
	defer rows.Close()
	for rows.Next() {
		var op domain.Group
		if err := rows.Scan(&op.ID, &op.Name); err != nil {
			return outputs, err
		}
		outputs = append(outputs, &op)
	}
	return outputs, nil
}

func Group(db *sql.DB, input *GroupInput) (int64, error) {
	id, err := db.ExecContext(context.Background(), "INSERT INTO `groups` (user_id, name) VALUES (?,?);", input.UserID, input.GroupName)
	if err != nil {
		log.Printf("failed to exec query err = %s", err.Error())
		return 0, err
	}
	last_id, err := id.LastInsertId()
	if err != nil {
		log.Printf("failed to get a last insert id err = %s", err.Error())
		return 0, err
	}
	return last_id, nil
}
