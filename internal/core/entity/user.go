package entity

type User struct {
	ID        int64   `db:"id" json:"id"`
	Username  string  `db:"username" json:"username"`
	Email     string  `db:"email" json:"email"`
	Password  string  `db:"password" json:"password"`
	CreatedAt int64   `db:"created_at" json:"created_at"`
	Roles     []uint8 `db:"roles" json:"roles"`
}