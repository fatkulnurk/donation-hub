package project

type FileStorage interface {
	RequestUploadUrl(mimeType string, fileSize int64) (url string, expiredAt int64, err error)
}

type DataStorage interface {
	Submit() (err error)
	ReviewByAdmin() (err error)
	Get() (err error)
	GetById() (err error)
	DonateToProject() (err error)
	GetDonationById() (err error)
}
