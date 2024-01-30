package models

import (
	"time"

	"gorm.io/gorm"
)

type ExpensesTable struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title" gorm:"not null"`
	Amount     float64    `json:"amount" gorm:"not null"`
	OnDate     time.Time  `json:"on_date" gorm:"not null"`
	Note       string     `json:"note"`
	CreatedAt  time.Time  `json:"-" gorm:"type:time`
	UpdatedAt  time.Time  `json:"-" gorm:"type:time"`
	UserId     uint       `json:"user_id" gorm:"not null"`
	User       UserBase   `json :"-" gorm:"foreignKey:UserId"`
	CategoryId uint       `json:"category_id" gorm:"not null"`
	Category   Categories `json :"-" gorm:"foreignKey:CategoryId"`
	Total      float64    `json:"total"`
}

// Utilities
func (expenses *ExpensesTable) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&expenses).Count(&total)
	return total
}

func (expens *ExpensesTable) Take(db *gorm.DB, limit int, offset int, userId uint) interface{} {
	var expenses []ExpensesTable
	db.Where("user_id = ?", userId).Preload("Category").Preload("User").Offset(offset).Limit(limit).Find(&expenses)
	return expenses
}

// Expenses Table
func (expenses *ExpensesTable) BeforeCreate(tx *gorm.DB) error {
	expenses.CreatedAt = time.Now()
	expenses.UpdatedAt = time.Now()
	return nil
}

func (expenses *ExpensesTable) BeforeUpdate(tx *gorm.DB) error {
	expenses.UpdatedAt = time.Now()
	return nil
}
