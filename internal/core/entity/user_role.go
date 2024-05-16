package entity

type UserRole struct {
	UserId int64  `db:"user_id"`
	Role   string `db:"role"`
}
