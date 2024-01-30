package models

import "golang.org/x/crypto/bcrypt"

type UserBase struct {
	Id       uint   `json:"id"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
	IsActive bool   `json:"is_active"`
	IsVerify bool   `json:"is_verify"`
	Username string `json:"username" gorm:"unique"`
}

type UserProfile struct {
	Id          uint     `json:"id"`
	UserId      uint     `json:"user_id"`
	User        UserBase `json:"-" gorm:"one2one;unique"`
	FirstName   string   `json:"first_name" gorm:"not null"`
	LastName    string   `json:"last_name" gorm:"not null"`
	ProfilePic  string   `json:"profile_pic"`
	Profession  string   `json:"profession"`
	Description string   `json:"description" gorm:"type:text"`
}

func (user *UserBase) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func CreateHashedPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword)
}

func (user *UserBase) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
