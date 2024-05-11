package entity

type User struct {
	ID       uint32
	Username string
	Email    string
	Password string
	Role     []string
}
