package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // например: "admin" или "user"
}
