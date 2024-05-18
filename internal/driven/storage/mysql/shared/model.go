package shared

// abaikan dulu file ini
import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRole struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

type Project struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	TargetAmount     float64   `json:"target_amount"`
	CollectionAmount float64   `json:"collection_amount"`
	Currency         string    `json:"currency"`
	Status           string    `json:"status"`
	RequesterID      int64     `json:"requester_id"`
	DueAt            time.Time `json:"due_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ProjectImage struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	URL       string `json:"url"`
}

type Donation struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	DonorID   int64     `json:"donor_id"`
	Message   string    `json:"message,omitempty"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}
