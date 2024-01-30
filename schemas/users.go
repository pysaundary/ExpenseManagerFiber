package schemas

import "expensesManage/models"

type UserBase struct {
	Email    string `json : "email gorm:"unique;not null"`
	Password string `json : "password"`
	IsActive bool   `json : "is_active`
	IsVerify bool   `json : "is_verify`
	Username string `json:"username gorm:"unique;not null"`
}

type ChangePassword struct {
	Password        string `json:"confirm_password"`
	ConfirmPassword string `json:"password"`
}

type UserProfile struct {
	UserBase    models.UserBase    `json:"user_base_data"`
	UserProfile models.UserProfile `json:"user_profile"`
}
