package entity

type ProjectImage struct {
	ID        int64  `db:"id"`
	ProjectID int64  `db:"project_id"`
	Url       string `db:"url"`
}
