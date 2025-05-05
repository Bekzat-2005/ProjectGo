package models

type Product struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Name       string  `json:"name"`
	CategoryID uint    `json:"category_id"`
	Price      float64 `json:"price"`
	UserID     uint    `json:"user_id"` // 👈 Жаңа бағана
}
