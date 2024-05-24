package req

type RegisterReqBody struct {
	Username string `json:"username" validate:"nonzero,max=255"`
	Email    string `json:"email" validate:"nonzero,max=255,regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password string `json:"password" validate:"min=6,max=255"`
	Role     string `json:"role" validate:"regexp=^(donor|requester)$"`
}

type LoginReqBody struct {
	Username string `json:"username" validate:"nonzero,max=255"`
	Password string `json:"password" validate:"min=6,max=255"`
}

type SubmitProjectReqBody struct {
	Title        string   `json:"title" validate:"nonzero,max=255"`
	Description  string   `json:"description" validate:"nonzero,max=255"`
	ImageUrls    []string `json:"image_urls" validate:"nonzero"`
	DueAt        int64    `json:"due_at" validate:"nonzero"`
	TargetAmount int64    `json:"target_amount" validate:"nonzero"`
	Currency     string   `json:"currency" validate:"nonzero,max=255"`
}

type ReviewProjectReqBody struct {
	Status string `json:"status" validate:"regexp=^(need_review|approved|completed|rejected)$"`
}

type DonateToProjectReqBody struct {
	Amount   int64  `json:"amount" validate:"nonzero"`
	Currency string `json:"currency" validate:"nonzero,max=255"`
	Message  string `json:"message" validate:"nonzero,max=255"`
}
