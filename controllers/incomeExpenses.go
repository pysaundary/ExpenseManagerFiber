package controllers

import (
	"expensesManage/database"
	"expensesManage/middleware"
	"expensesManage/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddIncome(c *fiber.Ctx) error {
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
	}
	var data models.IncomeTable
	if err := c.BodyParser(&data); err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed 1",
			"reason": err.Error(),
		})
	}
	data.UserId = *userInfo.Id
	if err := database.DB.Create(&data).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	var user models.UserBase
	if err := database.DB.Where("id = ?", data.UserId).First(&user).Error; err != nil {
		// Handle the error
		return err
	}
	data.User = user

	var category models.Categories
	if err := database.DB.Where("id = ?", data.CategoryId).First(&category).Error; err != nil {
		// Handle the error
		return err
	}
	data.Category = category
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   &data,
		"reason": "created",
	})
}

func GetIncomes(c *fiber.Ctx) error {
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
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("page_size", "5"))
	return c.JSON(models.PaginateData(database.DB, &models.IncomeTable{}, limit, page, *userInfo.Id))
}

func GetSingleIncome(c *fiber.Ctx) error {
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
	}
	var incomeRecord models.IncomeTable
	incomeId := c.Params("id")
	if err := database.DB.Where("id = ?", incomeId).Preload("User").Preload("Category").First(&incomeRecord).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   incomeRecord,
		"reason": "retreived",
	})
}

func UpdateIncome(c *fiber.Ctx) error {
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
	}
	incomeId := c.Params("id")
	var incomeRecord models.IncomeTable
	var data models.IncomeTable
	if err := c.BodyParser(&data); err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed 1",
			"reason": err.Error(),
		})
	}
	database.DB.Where("id = ?", incomeId).First(&incomeRecord)
	if err := database.DB.Model(&incomeRecord).Updates(data).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	var user models.UserBase
	if err := database.DB.Where("id = ?", incomeRecord.UserId).First(&user).Error; err != nil {
		// Handle the error
		return err
	}
	incomeRecord.User = user
	var category models.Categories
	if err := database.DB.Where("id = ?", incomeRecord.CategoryId).First(&category).Error; err != nil {
		// Handle the error
		return err
	}
	incomeRecord.Category = category
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   incomeRecord,
		"reason": "Updated",
	})
}

func DeleteIncome(c *fiber.Ctx) error {
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
	}
	var incomeRecord models.IncomeTable
	incomeId := c.Params("id")
	if err := database.DB.Where("id = ?", incomeId).Preload("User").Preload("Category").Delete(&incomeRecord).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   "deleted",
		"reason": "successful",
	})
}

// Expenses

func AddExpense(c *fiber.Ctx) error {
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
	}
	var data models.ExpensesTable
	if err := c.BodyParser(&data); err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed 1",
			"reason": err.Error(),
		})
	}
	data.UserId = *userInfo.Id
	if err := database.DB.Create(&data).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	var user models.UserBase
	if err := database.DB.Where("id = ?", data.UserId).First(&user).Error; err != nil {
		// Handle the error
		return err
	}
	data.User = user

	var category models.Categories
	if err := database.DB.Where("id = ?", data.CategoryId).First(&category).Error; err != nil {
		// Handle the error
		return err
	}
	data.Category = category
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   &data,
		"reason": "created",
	})
}

func GetExpenses(c *fiber.Ctx) error {
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
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("page_size", "5"))
	return c.JSON(models.PaginateData(database.DB, &models.ExpensesTable{}, limit, page, *userInfo.Id))
}

func GetSingleExpense(c *fiber.Ctx) error {
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
	}
	var eRecord models.ExpensesTable
	id := c.Params("id")
	if err := database.DB.Where("id = ?", id).Preload("User").Preload("Category").First(&eRecord).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   eRecord,
		"reason": "retreived",
	})
}

func UpdateExpense(c *fiber.Ctx) error {
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
	}
	id := c.Params("id")
	var eRecord models.ExpensesTable
	var data models.ExpensesTable
	if err := c.BodyParser(&data); err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed 1",
			"reason": err.Error(),
		})
	}
	database.DB.Where("id = ?", id).First(&eRecord)
	if err := database.DB.Model(&eRecord).Updates(data).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	var user models.UserBase
	if err := database.DB.Where("id = ?", eRecord.UserId).First(&user).Error; err != nil {
		// Handle the error
		return err
	}
	eRecord.User = user
	var category models.Categories
	if err := database.DB.Where("id = ?", eRecord.CategoryId).First(&category).Error; err != nil {
		// Handle the error
		return err
	}
	eRecord.Category = category
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   eRecord,
		"reason": "Updated",
	})
}

func DeleteExpense(c *fiber.Ctx) error {
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
	}
	var eRecord models.ExpensesTable
	id := c.Params("id")
	if err := database.DB.Where("id = ?", id).Preload("User").Preload("Category").Delete(&eRecord).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": false,
			"data":   "failed",
			"reason": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data":   "deleted",
		"reason": "successful",
	})
}
