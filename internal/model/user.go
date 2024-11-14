package model

// UserInfo модель юзера, которая содержит все данные для сохранения в payload jwt токена
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserLogin struct {
	Username string
	Password string
}
