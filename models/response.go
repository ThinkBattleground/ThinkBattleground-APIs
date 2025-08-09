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
	Id         string `json:"id"`
	UserId     string `json:"user_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}
