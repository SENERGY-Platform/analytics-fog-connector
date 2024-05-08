package auth

type MockAuth struct {
	ReturnError bool;
	Error error;
}

func (m *MockAuth) GetUserID(username string, password string) (string, error) {
	if m.ReturnError {
		return "", m.Error
	}
	return "", nil
}