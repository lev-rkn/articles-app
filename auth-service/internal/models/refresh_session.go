package models

type RefreshSession struct {
	RefreshToken string
	Fingerprint  string
	UserEmail    string
	AppId        int32
}
