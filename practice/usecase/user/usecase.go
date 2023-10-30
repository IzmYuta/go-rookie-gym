package user

import (
	"context"
	"database/sql"
	"log"
	// groupd "practice/domain/group"
	userd "practice/domain/user"
)

type Usecase interface {
	User(db *sql.DB, input UserInput) (*UserOutput, error)
}


type UserInput struct {
	Name string
}

// DBから取得した数値がint64型のためこの形式にしている
type UserOutput struct {
	UserID  int64
	GroupID int64
}

// (u *usecase)は、メソッドレシーバによってusecaseにUser関数を結びつける
// これなんの意味があるの？
func User(db *sql.DB, input *UserInput) (*UserOutput, error) {
	user := userd.NewUser(input.Name)
	var output UserOutput
	// DB処理
	id, err := db.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?);", user.Name)
	if err != nil {
		log.Printf("failed to exec query err = %s", err.Error())
		return &output, err
	}
	last_user_id, err := id.LastInsertId()
	if err != nil {
		log.Printf("failed to get a last insert id err = %s", err.Error())
		return &output, err
	}
	group_id, err := db.ExecContext(context.Background(), "INSERT INTO `groups` (user_id, name) VALUES (?, ?);", last_user_id, "default group")
	if err != nil {
		log.Printf("failed to exec query err = %s", err.Error())
		return &output, err
	}
	last_group_id, err := group_id.LastInsertId()
	if err != nil {
		log.Printf("failed to get a last insert id err = %s", err.Error())
		return &output, err
	}
	// 値をreturn
	output.UserID = last_user_id
	output.GroupID = last_group_id
	return &output, nil
}
