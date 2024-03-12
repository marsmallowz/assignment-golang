package model

type Item struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `json:"itemCode"`
	Quantity    uint
	Description string
	OrderID     uint
}
