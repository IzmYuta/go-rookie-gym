//go:generate mockgen -source=$GOFILE -package=mock -destination=../mock/$GOPACKAGE/$GOFILE
package user

import (
	"database/sql"
)

type Repository interface {
	User(db *sql.DB, user *User) (int, error)
}
