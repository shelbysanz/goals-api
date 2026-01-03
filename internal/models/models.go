package models

import "time"

type MonthGoal struct {
	ID        uint      `gorm:"primaryKey"`
	Month     time.Time `gorm:"type:date;not null;index"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Completed bool      `gorm:"not null; default: false"`
	Notes     string    `gorm:"type:text"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type WeekGoal struct {
	ID        uint      `gorm:"primaryKey"`
	WeekStart time.Time `gorm:"type:date;not null;index"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Completed bool      `gorm:"not null; default: false"`
	Notes     string    `gorm:"type:text"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type WeekTodo struct {
	ID        uint      `gorm:"primaryKey"`
	WeekStart time.Time `gorm:"type:date;not null;index"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Completed bool      `gorm:"not null; default: false"`
	Notes     string    `gorm:"type:text"`
	UpdatedAt time.Time
	CreatedAt time.Time
}
