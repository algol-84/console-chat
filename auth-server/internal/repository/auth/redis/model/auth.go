package model

// User определяет юзера в кэше
type User struct {
	ID          int64  `redis:"id"`
	Name        string `redis:"name"`
	Email       string `redis:"email"`
	Role        string `redis:"role"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at"`
}
