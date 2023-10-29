//go:generate mockgen -source=$GOFILE -package=mock -destination=../mock/$GOPACKAGE/$GOFILE
package group

import(
	"database/sql"
)

type Repository interface{
	Group(db *sql.DB, group *Group) (int , error)
	Groups(db *sql.DB, id int) ([]Group, error)
}