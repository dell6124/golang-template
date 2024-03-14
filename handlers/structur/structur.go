package structur

type CreateUserRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DeleteUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ChangePasswordRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}