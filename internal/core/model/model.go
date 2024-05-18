package model

import "time"

type UserRegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserRegisterOutput struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginOutput struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type ListUserInput struct {
	Limit int64  `json:"limit"`
	Page  int64  `json:"page"`
	Role  string `json:"role"`
}

type ListUser struct {
	ID       int64          `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Roles    []ListUserRole `json:"roles"`
}

type ListUserRole struct {
	Role string `json:"role"`
}

type ListUserMeta struct {
	Page       int64 `json:"page"`
	TotalPages int64 `json:"total_pages"`
}

type ListUserOutput struct {
	Users      []ListUser   `json:"users"`
	Pagination ListUserMeta `json:"pagination"`
}

type RequestUploadUrlInput struct {
	MimeType string `json:"mime_type"`
	FileSize int64  `json:"file_size"`
}

type RequestUploadUrlOutput struct {
	MimeType  string `json:"mime_type"`
	FileSize  int64  `json:"file_size"`
	URL       string `json:"url"`
	ExpiresAt int64  `json:"expires_at"`
}

type SubmitProjectInput struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ImageURLs    []string `json:"image_urls"`
	DueAt        int64    `json:"due_at"`
	TargetAmount int64    `json:"target_amount"`
	Currency     string   `json:"currency"`
}

type SubmitProjectOutput struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ImageURLs    []string `json:"image_urls"`
	DueAt        int64    `json:"due_at"`
	TargetAmount int64    `json:"target_amount"`
	Currency     string   `json:"currency"`
}

type ApprovalStatusInput struct {
	Status string `json:"status"`
}

type ListProjectInput struct {
	Status  string    `json:"status"`
	Limit   int64     `json:"limit"`
	StartTs time.Time `json:"start_ts"` // jangan lupa, ini nanti Unix timestamp
	EndTs   time.Time `json:"end_ts"`   // jangan lupa, ini nanti Unix timestamp
	LastKey string    `json:"last_key"`
}

type Requester struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetProjectByIdOutput struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ImageURLs    []string  `json:"image_urls"`
	DueAt        int64     `json:"due_at"`
	TargetAmount int64     `json:"target_amount"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	Requester    Requester `json:"requester"`
}

type GetProjectOutput struct {
	Projects []GetProjectByIdOutput `json:"projects"`
	LastKey  string                 `json:"last_key"`
}

type DonateToProjectInput struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Message  string `json:"message"`
}

type ListProjectDonationInput struct {
	Limit   int64  `json:"limit"`
	LastKey string `json:"last_key"`
}

type Donor struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Donation struct {
	ID        int64  `json:"id"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	Message   string `json:"message"`
	Donor     Donor  `json:"donor"`
	CreatedAt int64  `json:"created_at"`
}

type ListProjectDonationOutput struct {
	Donations []Donation `json:"donations"`
	LastKey   string     `json:"last_key"`
}
