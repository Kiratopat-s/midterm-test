package model

import "github.com/Kiratopat-s/workflow/internal/constant"

type Item struct {
	ID       uint           `gorm:"primaryKey;autoIncrement" json:"id"`
    Title    string         `gorm:"size:255;not null" json:"title"`
    Amount   int            `gorm:"not null" json:"amount"`
    Quantity int            `gorm:"not null" json:"quantity"`
    Status   constant.ItemStatus         `gorm:"size:20;not null" json:"status"`
    OwnerID  int            `gorm:"not null" json:"owner_id"`
}

