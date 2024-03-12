package model

import (
	"time"
)

type Order struct {
	ID           uint `gorm:"primaryKey"`
	CustomerName string
	Items        []Item
	OrderedAt    time.Time
}
