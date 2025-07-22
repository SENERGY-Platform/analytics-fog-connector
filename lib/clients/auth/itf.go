package auth

type IDP interface {
	GetUserID() (string, error)
}
