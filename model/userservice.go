package model

type ChangePasswordRequest struct {
	UserID      string `json:"-"`
	Password    string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

type ChangePasswordResponse struct {
	Code         int
	Message      string
	SessionToken *SessionToken
}

type ChangeProfileRequest struct {
	UserID      string `json:"-"`
	Displayname string `json:"displayname"`
}

type ChangeProfileResponse struct {
	Code         int
	Message      string
	SessionToken *SessionToken
}

type ChangeUsernameRequest struct {
	UserID      string `json:"-"`
	Password    string `json:"password"`
	NewUsername string `json:"newusername"`
}

type ChangeUsernameResponse struct {
	Code         int
	Message      string
	SessionToken *SessionToken
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code         int
	Message      string
	SessionToken *SessionToken
}
