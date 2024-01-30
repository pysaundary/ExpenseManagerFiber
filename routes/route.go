package routes

import (
	"expensesManage/controllers"
	"expensesManage/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Use(middleware.IsAuthenticated)
	app.Static("api/uploads/", "./uploads")
	app.Post("/apis/register-user/", controllers.RegisterUser)
	app.Post("/apis/login-user/", controllers.LoginUser)
	app.Post("/apis/logout-user/", controllers.Logout)
	app.Post("/apis/forget-password-request/", controllers.ForgetPassword)
	app.Post("/apis/change-password/", controllers.ChangePassword)
	app.Get("/apis/user-profile/", controllers.GetMyProfile)
	app.Put("/apis/user-profile/", controllers.UpdateMyProfile)
	app.Delete("/apis/user-profile/", controllers.DeleteMyProfile)

	// Expense and Income
	app.Post("/apis/add-income/", controllers.AddIncome)
	app.Get("/apis/incomes/", controllers.GetIncomes)
	app.Get("/apis/user-income/:id/", controllers.GetSingleIncome)
	app.Put("/apis/user-income/:id/", controllers.UpdateIncome)
	app.Delete("/apis/user-income/:id/", controllers.DeleteIncome)

	app.Post("/apis/add-expense/", controllers.AddExpense)
	app.Get("/apis/expenses/", controllers.GetExpenses)
	app.Get("/apis/user-expense/:id/", controllers.GetSingleExpense)
	app.Put("/apis/user-expense/:id/", controllers.UpdateExpense)
	app.Delete("/apis/user-expense/:id/", controllers.DeleteExpense)
}
