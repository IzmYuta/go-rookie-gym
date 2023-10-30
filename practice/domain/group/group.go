package group

type Group struct {
	ID     int
	UserID int
	Name   string
}

// Group構造体の初期化を簡単にするための関数
func NewGroup(user_id int, name string) *Group {
	return &Group{
		UserID: user_id,
		Name:   name,
	}
}
