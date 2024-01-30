package controllers

import (
	"expensesManage/config"
	"expensesManage/database"
	"expensesManage/middleware"
	"expensesManage/models"
	"expensesManage/schemas"
	"expensesManage/utilities"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	userBaseValues := make(map[string]string)
	userProfileValues := make(map[string]interface{})
	for key, item := range form.Value {
		if len(item) > 0 {
			if key == "email" {
				userBaseValues[key] = item[0]
			} else if key == "password" {
				userBaseValues[key] = item[0]
			} else {
				userProfileValues[key] = item[0]
			}
		}
	}
	uploadPath := config.ProfilePicPath

	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	files := form.File["profile_pic"]
	filename := ""
	for _, file := range files {
		filename = file.Filename
		if err := c.SaveFile(file, config.ProfilePicPath+filename); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": err.Error(),
			})
		}
	}
	userProfileValues["profile_pic"] = filename
	if utilities.IsValidEmail(userBaseValues["email"]) {
		username, err := utilities.GenerateSlugFromEmail(userBaseValues["email"])
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": err.Error(),
			})
		}
		userBaseValues["username"] = username
		user := models.UserBase{
			Email:    userBaseValues["email"],
			Username: username,
			IsActive: true,
			IsVerify: false,
		}
		// fmt.Println(userBaseValues["password"])
		user.SetPassword(userBaseValues["password"])
		if err := database.DB.Create(&user).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": err.Error(),
			})
		}
		userProfileValues["user_id"] = user.Id
		if err := database.DB.Model(&models.UserProfile{}).Create(userProfileValues).Error; err != nil {

			return c.Status(400).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": err.Error(),
			})
		}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   "success",
		"reason": "register successfully",
	})
}

func LoginUser(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed 1",
			"reason": err.Error(),
		})
	}

	var userBase models.UserBase
	if err := database.DB.Where("email=?", data["email"]).First(&userBase).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed 2",
			"reason": err.Error(),
		})
	}
	if err := userBase.CheckPassword(data["password"]); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed 3",
			"reason": err.Error(),
		})
	}
	claim, err := utilities.GenerateJWTToken(strconv.Itoa(int(userBase.Id)))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": false,
			"data":   "failed 4",
			"reason": err.Error(),
		})
	}
	token := claim

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"data":    "passed",
		"message": "success",
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"data":    "passed",
		"message": "success",
	})
}

func ForgetPassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}

	if len(data["email"]) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "please pass the email",
		})
	}
	var user models.UserBase
	var count int64
	database.DB.Model(&models.UserBase{}).Where("email = ?", data["email"]).Count(&count)
	if count == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "no user with email existed",
		})
	} else {
		database.DB.Model(&models.UserBase{}).Where("email = ?", data["email"]).First(&user)
		jwtToken, err := utilities.GenerateForgetPassWordJWTToken(strconv.FormatUint(uint64(user.Id), 10))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": "please pass the email",
			})
		}
		var linkToSend string = config.FrontEndPath + "/retrieveMyPassword/?token=" + jwtToken
		isPassed := utilities.SendEmail(user.Email, "Welcome To Expenses Mangement", "Hi please click on link for restore password "+linkToSend)
		if isPassed {
			// fmt.Println(user.Email)
			return c.Status(http.StatusCreated).JSON(fiber.Map{
				"status":  true,
				"data":    "passed",
				"message": "success",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "failed to send email",
		})
	}
}

func ChangePassword(c *fiber.Ctx) error {
	// Retrieve user information from the context
	userInfo, ok := c.Locals(middleware.UserInfoKey).(middleware.UserInfoMiddleware)

	if !ok || !userInfo.IsLogin || userInfo.Id == nil {
		// Handle the case when user information is not present or not logged in
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "please login before try",
		})
	}

	// Convert the user ID to int for database query
	id := int(*userInfo.Id)

	var user models.UserBase
	database.DB.Where("id = ?", id).First(&user)
	var data schemas.ChangePassword

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "please send password and confirm_password",
		})
	}
	if data.Password != data.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "password and confirm_password not matched",
		})
	} else {
		user.SetPassword(data.Password)
		// fmt.Println(user.Email)
		if err := database.DB.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": err.Error(),
			})
		} else {
			return c.Status(http.StatusCreated).JSON(fiber.Map{
				"status":  true,
				"data":    "passed",
				"message": "success",
			})
		}
	}
}

func GetMyProfile(c *fiber.Ctx) error {
	userInfo, ok := c.Locals(middleware.UserInfoKey).(middleware.UserInfoMiddleware)
	if !ok || !userInfo.IsLogin || userInfo.Id == nil {
		// Handle the case when user information is not present or not logged in
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "please login before try",
		})
	} else {
		id := int(*userInfo.Id)

		var user models.UserBase
		database.DB.Where("id = ?", id).First(&user)
		var userProfile models.UserProfile
		database.DB.Where("user_id = ?", user.Id).First(&userProfile)
		userProfile.ProfilePic = config.BASE_URL + config.ProfilePicUrl + userProfile.ProfilePic
		return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data": fiber.Map{
				"user":         user,
				"user_profile": userProfile,
			},
			"reason": "success",
		})
	}
}

func UpdateMyProfile(c *fiber.Ctx) error {
	userInfo, ok := c.Locals(middleware.UserInfoKey).(middleware.UserInfoMiddleware)
	if !ok || !userInfo.IsLogin || userInfo.Id == nil {
		// Handle the case when user information is not present or not logged in
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "please login before try",
		})
	} else {
		id := int(*userInfo.Id)

		var user models.UserBase
		database.DB.Where("id = ?", id).First(&user)
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": err.Error(),
			})
		}
		userProfileValues := make(map[string]interface{})
		for key, item := range form.Value {
			if len(item) > 0 {
				userProfileValues[key] = item[0]
			}
		}
		uploadPath := config.ProfilePicPath

		err = os.MkdirAll(uploadPath, os.ModePerm)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": err.Error(),
			})
		}
		files := form.File["profile_pic"]
		filename := ""
		for _, file := range files {
			filename = file.Filename
			if err := c.SaveFile(file, config.ProfilePicPath+filename); err != nil {
				return c.Status(400).JSON(fiber.Map{
					"status": false,
					"data":   "failed",
					"reason": err.Error(),
				})
			}
		}
		userProfileValues["profile_pic"] = filename
		var userProfile models.UserProfile
		database.DB.Where("user_id = ?", user.Id).First(&userProfile)

		// Update UserProfile fields
		if err := database.DB.Model(&userProfile).Updates(userProfileValues).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": false,
				"data":   "failed",
				"reason": err.Error(),
			})
		}
		database.DB.Where("user_id = ?", user.Id).First(&userProfile)
		userProfile.ProfilePic = config.BASE_URL + config.ProfilePicUrl + userProfile.ProfilePic
		return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data": fiber.Map{
				"user":         user,
				"user_profile": userProfile,
			},
			"reason": "success",
		})
	}
}

func DeleteMyProfile(c *fiber.Ctx) error {
	userInfo, ok := c.Locals(middleware.UserInfoKey).(middleware.UserInfoMiddleware)
	if !ok || !userInfo.IsLogin || userInfo.Id == nil {
		// Handle the case when user information is not present or not logged in
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": "please login before try",
		})
	} else {
		id := int(*userInfo.Id)

		var user models.UserBase
		database.DB.Where("id = ?", id).First(&user).Update("IsActive", false)
		return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data":   "your account is now deactivate, after 40 days your account will be complete deleted from server",
			"reason": "success",
		})
	}
}
