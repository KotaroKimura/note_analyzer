package model

import (
	"time"
)

type Article struct {
	ID        uint `xorm:"name 'id'"`
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Article) TableName() string {
	return "articles"
}
