package entity

type Project struct {
	ID               int64   `db:"id"`
	Name             string  `db:"name"`
	Description      string  `db:"description"`
	TargetAmount     float64 `db:"target_amount"`
	CollectionAmount float64 `db:"collection_amount"`
	Currency         string  `db:"currency"`
	Status           string  `db:"status"`
	RequesterID      int     `db:"requester_id"`
	DueAt            int64   `db:"due_at"`
	CreatedAt        int64   `db:"created_at"`
	UpdatedAt        int64   `db:"updated_at"`
}
