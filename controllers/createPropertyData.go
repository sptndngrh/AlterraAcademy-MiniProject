package controllers

import (
	"net/http"
	"sewakeun_project/middlewares"
	"sewakeun_project/models"
	"sewakeun_project/response"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreatePropertyData(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Memeriksa token otorisasi
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi tidak ditemukan"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Format token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		tokenString := authParts[1]

		username, err := middlewares.VerifyTokenJWT(tokenString, secretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Otorisasi pengguna
		var ownerRole models.User
		result := db.Where("username = ?", username).First(&ownerRole)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Anda bukan owner!"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Menerima data input dari request
		var CreatePropertyData struct {
			Judul     string `json:"judul"`
			Tipe      string `json:"tipe"`
			Harga     int    `json:"harga"`
			Lokasi    string `json:"lokasi"`
			Deskripsi string `json:"deskripsi"`
		}

		if err := c.Bind(&CreatePropertyData); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Membuat orderedProperty dari data yang dimasukkan
		orderedProperty := models.Property{
			Judul:     CreatePropertyData.Judul,
			Tipe:      CreatePropertyData.Tipe,
			Harga:     CreatePropertyData.Harga,
			Lokasi:    CreatePropertyData.Lokasi,
			Deskripsi: CreatePropertyData.Deskripsi,
			OwnerId:   ownerRole.Id,
		}

		// Memeriksa apakah ada order sewa yang masih berlaku
		var existingOrder models.Order
		result = db.Where("properti_id = ? AND payment_status = ?", orderedProperty.Id, true).First(&existingOrder)
		if result.Error == nil {
			orderedProperty.StatusTersewa = true
		}

		// Simpan orderedProperty ke dalam database
		if err := db.Create(&orderedProperty).Error; err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam menambahkan data properti"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"PropertyData": orderedProperty,
			"code":         http.StatusOK,
			"error":        false,
			"message":      "Data property berhasil ditambahkan",
		})
	}
}
