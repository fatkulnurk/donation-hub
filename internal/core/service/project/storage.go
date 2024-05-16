package project

type FileStorage interface {
	RequestUploadUrl(mimeType string, fileSize int) (url string, err error)
}

type DataStorage interface {
	Submit() (err error)
	ReviewByAdmin() (err error)
	Get() (err error)
	GetById() (err error)
	DonateToProject() (err error)
	GetDonationById() (err error)
}
