package routes

import (
	"log"
	"net/http"
	"os"
	"sewakeun_project/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	e.Use(LoggerRoute())
	secretKey := []byte(getSecretKeyFromEnv())
	e.POST("/register/owner", controllers.RegisterOwner(db, secretKey))
	e.POST("/register/user", controllers.RegisterUser(db, secretKey))
	e.GET("/verify/:context", func(c echo.Context) error {
		context := c.Param("context")
		if context == "user" {
			return controllers.UserVerifyEmail(db)(c)
		} else if context == "owner" {
			return controllers.OwnerVerifyEmail(db)(c)
		} else {
			return c.String(http.StatusNotFound, "Rute tidak ditemukan")
		}
	})
	e.POST("/login/user", controllers.LoginUser(db, secretKey))
	e.POST("/login/owner", controllers.LoginOwner(db, secretKey))
	e.PUT("/change-username-owner/owner/:owner_id", controllers.ChangeUserNameOwner(db, secretKey))
	e.PUT("/change-username-user/user/:user_id", controllers.ChangeUserNameUser(db, secretKey))
	e.PUT("/change-password-owner/owner/:owner_id", controllers.ChangePasswordOwner(db, secretKey))
	e.PUT("/change-password-user/user/:user_id", controllers.ChangePasswordUser(db, secretKey))
	e.PUT("/change-name-owner/owner/:owner_id", controllers.ChangeNameOwner(db, secretKey))
	e.PUT("/change-name-user/user/:user_id", controllers.ChangeNameUser(db, secretKey))

}

func getSecretKeyFromEnv() string {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("Kunci rahasia tidak ditemukan di .env")
	}
	return secretKey
}
