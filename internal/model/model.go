package model

type User struct {
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password"`
}

type UserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
