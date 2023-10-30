package user

type User struct {
	ID   int
	Name string
}

// User構造体の初期化を簡単にするための関数
func NewUser(name string) *User {
	return &User{
		Name: name,
	}
}