package entity

type Donation struct {
	ID        int64   `db:"id"`
	ProjectID int64   `db:"project_id"`
	DonorID   int64   `db:"donor_id"`
	Message   string  `db:"message"`
	Amount    float64 `db:"amount"`
	Currency  string  `db:"currency"`
	CreatedAt int64   `db:"created_at"`
}
