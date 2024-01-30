package models

import (
	"time"

	"gorm.io/gorm"
)

type IncomeTable struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title" gorm:"not null"`
	Amount     float64    `json:"amount" gorm:"not null"`
	OnDate     time.Time  `json:"on_date" gorm:"not null;type:time"`
	Note       string     `json:"note"`
	CreatedAt  time.Time  `json:"-" gorm:"type:date"`
	UpdatedAt  time.Time  `json:"-" gorm:"type:date"`
	UserId     uint       `json:"user_id" gorm:"not null"`
	User       UserBase   `json :"-" gorm:"foreignKey:UserId"`
	CategoryId uint       `json:"category_id" gorm:"not null"`
	Category   Categories `json :"-" gorm:"foreignKey:CategoryId"`
	Total      float64    `json:"total"`
}

// Utilities
func (incomes *IncomeTable) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&incomes).Count(&total)
	return total
}

func (income *IncomeTable) Take(db *gorm.DB, limit int, offset int, userId uint) interface{} {
	var incomes []IncomeTable
	db.Where("user_id = ?", userId).Preload("Category").Preload("User").Offset(offset).Limit(limit).Find(&incomes)
	return incomes
}

// Hooks
// Income table
func (incomeTable *IncomeTable) BeforeCreate(tx *gorm.DB) error {
	incomeTable.CreatedAt = time.Now()
	incomeTable.UpdatedAt = time.Now()
	return nil
}

func (incomeTable *IncomeTable) BeforeUpdate(tx *gorm.DB) error {
	incomeTable.UpdatedAt = time.Now()
	return nil
}
