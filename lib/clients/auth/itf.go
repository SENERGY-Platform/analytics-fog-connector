package auth

type AuthClient interface {
	GetUserID(string, string) (string, error)
}