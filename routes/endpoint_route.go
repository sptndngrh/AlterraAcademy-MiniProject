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

	e.GET("/verify", controllers.VerifyEmail(db))
	e.POST("/chatbot", func(c echo.Context) error {
		return controllers.RekomendasiPropertiChatBot(c, controllers.NewAiUsecase())
	})
	user := e.Group("/user")
	user.GET("/orders", controllers.GetAllUserOrders(db, secretKey))
	e.POST("/login", controllers.Login(db, secretKey))
	e.POST("/register", controllers.Register(db, secretKey))

	property := e.Group("/property")
	property.GET("/all", controllers.GetAllProperty(db, secretKey))
	property.POST("/create", controllers.CreatePropertyData(db, secretKey))
	property.POST("/order", controllers.CreateOrderingProperty(db, secretKey))
	property.PUT("/edit/:id", controllers.EditPropertyData(db, secretKey))
	property.DELETE("/delete/:id", controllers.DeletePropertyData(db, secretKey))

	change := e.Group("/change")
	change.PUT("/username/:id", controllers.ChangeUsername(db, secretKey))
	change.PUT("/password/:id", controllers.ChangePassword(db, secretKey))

	update := e.Group("/update")
	update.PUT("/payment/:id", controllers.UpdatePaymentStatus(db, secretKey))
}

func getSecretKeyFromEnv() string {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("Kunci rahasia tidak ditemukan di .env")
	}
	return secretKey
}
