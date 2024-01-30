package middleware

import (
	"expensesManage/database"
	"expensesManage/models"
	"expensesManage/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const UserInfoKey = "userContext"

type UserInfoMiddleware struct {
	IsLogin bool
	Id      *uint
}

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if data, err := utilities.ParseJWT(cookie); err != nil {
		userInfo := UserInfoMiddleware{
			IsLogin: false,
			Id:      nil,
		}
		c.Locals(UserInfoKey, userInfo)
	} else {
		// Convert data to uint
		userId, parseErr := strconv.ParseUint(data, 10, 64)
		if parseErr != nil {
			// Handle the error, e.g., log it, set Id to nil, etc.
			userInfo := UserInfoMiddleware{
				IsLogin: false,
				Id:      nil,
			}
			c.Locals(UserInfoKey, userInfo)
		} else {
			// Assign the converted uint value to Id
			userIdUint := uint(userId)
			var user models.UserBase
			if err := database.DB.Where("id = ?", userIdUint).First(&user).Error; err != nil {
				userInfo := UserInfoMiddleware{
					IsLogin: false,
					Id:      nil,
				}
				c.Locals(UserInfoKey, userInfo)
			}
			if !user.IsActive {
				userInfo := UserInfoMiddleware{
					IsLogin: false,
					Id:      nil,
				}
				c.Locals(UserInfoKey, userInfo)
			} else {
				userInfo := UserInfoMiddleware{
					IsLogin: true,
					Id:      &userIdUint,
				}
				c.Locals(UserInfoKey, userInfo)
			}
		}
	}
	return c.Next()
}
