package controllers

import (
	"net/http"
	"sewakeun_project/middlewares"
	"sewakeun_project/models"
	"sewakeun_project/response"
	"strings"
	"time"

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

func GetAllProperty(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Dapatkan ID pemilik properti dari token otorisasi
		userID := ownerRole.Id

		// Query properti yang dimiliki oleh pemilik properti
		var properties []models.Property
		if err := db.Where("owner_id = ?", userID).Find(&properties).Error; err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal mengambil data properti pemilik"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"properties": properties,
			"code":       http.StatusOK,
			"error":      false,
			"message":    "Daftar properti berhasil diambil",
		})
	}
}

func EditPropertyData(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Otorisasi pemilik properti
		var userRole models.User
		result := db.Where("username = ?", username).First(&userRole)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pengguna tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Parse JSON dari permintaan
		var updateData struct {
			Judul     string `json:"judul"`
			Tipe      string `json:"tipe"`
			Harga     int    `json:"harga"`
			Lokasi    string `json:"lokasi"`
			Deskripsi string `json:"deskripsi"`
		}

		if err := c.Bind(&updateData); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mendapatkan nilai ID dari parameter URL
		id := c.Param("id")

		var propertyData models.Property
		result = db.Where("id = ?", id).First(&propertyData)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Properti tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Periksa apakah pemilik properti adalah pemilik yang sah (misalnya, berdasarkan token otorisasi)
		if propertyData.OwnerId != userRole.Id {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Anda tidak memiliki izin untuk mengubah data properti ini"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Update data properti
		propertyData.Judul = updateData.Judul
		propertyData.Tipe = updateData.Tipe
		propertyData.Harga = updateData.Harga
		propertyData.Lokasi = updateData.Lokasi
		propertyData.Deskripsi = updateData.Deskripsi

		// Simpan perubahan ke dalam database
		if err := db.Save(&propertyData).Error; err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam mengupdate data properti"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"PropertyData": propertyData,
			"code":         http.StatusOK,
			"error":        false,
			"message":      "Data property berhasil diperbarui",
		})
	}
}

func DeletePropertyData(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Otorisasi pemilik properti
		var ownerRole models.User
		result := db.Where("username = ?", username).First(&ownerRole)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pengguna tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Dapatkan ID properti yang akan dihapus dari URL
		propertyID := c.Param("id")

		// Periksa apakah properti dengan ID yang sesuai ada dalam database
		var property models.Property
		result = db.Where("id = ?", propertyID).First(&property)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Properti tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Periksa apakah pemilik properti adalah pemilik yang sah (misalnya, berdasarkan token otorisasi)
		if property.OwnerId != ownerRole.Id {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Anda tidak memiliki izin untuk menghapus data properti ini"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Soft delete properti dengan mengatur deleted_at
		if err := db.Model(&property).Update("deleted_at", time.Now()).Error; err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam menghapus data properti"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Data property berhasil dihapus",
		})
	}
}
