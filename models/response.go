package models

type ResponseWithEmail struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type Response struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Message    string `json:"message"`
	Token      string `json:"token"`
	ExpireTime string `json:"expire_time"`
}

type UserGetResponse struct {
	Id       string `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
