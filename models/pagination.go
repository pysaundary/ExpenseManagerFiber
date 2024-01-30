package models

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Entity interface {
	Count(db *gorm.DB) int64
	Take(db *gorm.DB, limit int, offset int, userId uint) interface{}
}

func PaginateData(db *gorm.DB, entity Entity, limit int, page int, userId uint) fiber.Map {
	var total int64
	offset := (page - 1) * limit
	data := entity.Take(db, limit, offset, userId)
	total = entity.Count(db)
	return fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(int(total) / limit)),
		},
	}
}
