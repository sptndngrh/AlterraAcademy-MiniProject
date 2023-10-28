package configs

import (
	"log"
	"sewakeun_project/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes() *echo.Echo {
	db, err := initializeDB()
	if err != nil {
		log.Fatal(err)
	}
	router := echo.New()
	router.Use(middleware.Logger())
	routes.SetupRoutes(router, db)
	return router
}
