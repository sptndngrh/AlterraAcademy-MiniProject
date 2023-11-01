package configs

import (
	"log"
	"os"
	"sewakeun_project/models"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	DBHost     string
	DBPort     int
	DBUsername string
	DBPassword string
	DBName     string
}

func initializeDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbConfig := DBConfig{
		DBHost: os.Getenv("DBHOST"),
	}
	portStr := os.Getenv("DBPORT")
	dbConfig.DBPort, _ = strconv.Atoi(portStr)
	dbConfig.DBUsername = os.Getenv("DBUSER")
	dbConfig.DBPassword = os.Getenv("DBPASS")
	dbConfig.DBName = os.Getenv("DBNAME")

	dsn := dbConfig.DBUsername + ":" + dbConfig.DBPassword + "@tcp(" + dbConfig.DBHost + ":" + portStr + ")/" + dbConfig.DBName + "?charset=utf8&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{}, &models.Owner{}, models.Property{}, models.Ticket{})
	return db, nil
}
