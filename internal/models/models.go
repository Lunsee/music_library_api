package models

import (
	"time"
)

type Song struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Group       string    `json:"group" gorm:"column:group_name"`
	Song        string    `json:"song" gorm:"column:song"`
	ReleaseDate time.Time `json:"releaseDate" gorm:"column:release_date"`
	Text        string    `json:"text" gorm:"column:text"`
	Link        string    `json:"link" gorm:"column:link"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}
