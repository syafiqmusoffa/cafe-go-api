package routes

import (
	"my-go-api/internal/controllers"
	"my-go-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/v1")

	// ✅ Auth
	api.Post("/login", controllers.Login)

	// ✅ Admin Only
	api.Post("/products", middleware.AuthMiddleware("admin"), controllers.CreateProduct)
	api.Put("/product/:id", middleware.AuthMiddleware("admin"), controllers.UpdateProduct)
	api.Patch("/product/:id", middleware.AuthMiddleware("admin"), controllers.UpdateProductStatus)
	api.Delete("/product/:id", middleware.AuthMiddleware("admin"), controllers.DeleteProduct)
	api.Post("/tables", middleware.AuthMiddleware("admin"), controllers.CreateTable)
	api.Delete("/tables/:id", middleware.AuthMiddleware("admin"), controllers.DeleteTable)
	cashierController := controllers.CashierController{DB: db}

	api.Post("/cashiers", middleware.AuthMiddleware("admin"), cashierController.CreateCashier)
	api.Put("/cashiers/:id", middleware.AuthMiddleware("admin"), cashierController.UpdateCashier)
	api.Get("/cashiers", middleware.AuthMiddleware("admin"), cashierController.GetCashiers)
	api.Patch("/cashiers/:id/status", middleware.AuthMiddleware("admin"), cashierController.UpdateCashierStatus)
	api.Delete("/cashiers/:id", middleware.AuthMiddleware("admin"), cashierController.DeleteCashier)

	// ✅ Kasir
	api.Put("/orders/:id/status", middleware.AuthMiddleware("cashier"), controllers.UpdateOrderStatus)
	api.Get("/orders/active", middleware.AuthMiddleware("cashier"), controllers.GetActiveOrders)
	api.Get("/orders", middleware.AuthMiddleware("cashier"), controllers.GetAllOrders)

	// ✅ table
	api.Get("/orders/table/:table_id", controllers.GetOrdersByTable)
	api.Post("/orders", controllers.CreateOrder)
	api.Get("/products", controllers.GetProducts)
	api.Get("/tables", controllers.GetTables)
	api.Get("/tables/:id", controllers.GetTableByID)

}
