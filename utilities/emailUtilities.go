package utilities

import (
	"expensesManage/database"
	"expensesManage/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Slugify(input string) string {
	// Replace non-alphanumeric characters and hyphens with a space
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	slug := re.ReplaceAllString(input, " ")

	// Trim leading and trailing spaces
	slug = strings.TrimSpace(slug)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Convert to lowercase
	slug = strings.ToLower(slug)

	return slug
}

func IsValidEmail(email string) bool {
	// Define a regular expression pattern for a basic email validation
	// This is a simplified version and may not cover all edge cases
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

func GenerateSlugFromEmail(email string) (string, error) {
	// Split email at the "@" symbol
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email format")
	}

	// Create a slug from the first part of the email
	slug := Slugify(parts[0])

	// Check if the slug already exists in the database
	var count int64
	database.DB.Model(&models.UserBase{}).Where("username = ?", slug).Count(&count)
	if count != 0 {
		// User exists, get the latest ID
		var latestID uint
		database.DB.Model(&models.UserBase{}).Select("id").Order("id desc").Limit(1).Scan(&latestID)
		slug = slug + "_" + strconv.FormatUint(uint64(latestID), 10)
	}

	return slug, nil
}
