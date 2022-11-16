package domain

import (
	"time"
)

type Base struct {
	ID        int       `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time `json:"created_at"`
}
