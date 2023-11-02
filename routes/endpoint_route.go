package routes

import (
	"log"
	"os"
	"sewakeun_project/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	e.Use(LoggerRoute())
	secretKey := []byte(getSecretKeyFromEnv())
	// GET
	e.GET("/verify", controllers.VerifyEmail(db))
	// POST
	e.POST("/login/owner", controllers.LoginOwner(db, secretKey))
	e.POST("/login/user", controllers.LoginUser(db, secretKey))
	e.POST("/register", controllers.Register(db, secretKey))
	e.POST("/property/create-property", controllers.CreatePropertyData(db, secretKey))
	e.POST("/property/order-property", controllers.CreateOrderingProperty(db, secretKey))
	// PUT
	e.PUT("/user/change-username/:id", controllers.ChangeUsername(db, secretKey))
	e.PUT("/user/change-password/:id", controllers.ChangePassword(db, secretKey))
	e.PUT("/updatePaymentStatus", controllers.UpdatePaymentStatus(db, secretKey))
	// DELETE
}

func getSecretKeyFromEnv() string {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("Kunci rahasia tidak ditemukan di .env")
	}
	return secretKey
}
