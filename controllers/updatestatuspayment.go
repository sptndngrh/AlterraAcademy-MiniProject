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

func UpdatePaymentStatus(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		var UpdatePaymentData struct {
			PaymentStatus bool `json:"payment_status"`
		}

		if err := c.Bind(&UpdatePaymentData); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Memeriksa apakah nilai payment_status adalah "true"
		if UpdatePaymentData.PaymentStatus != true {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: "Nilai payment_status harus berisi 'true'"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mendapatkan nilai ID dari parameter URL
		id := c.Param("id")

		var orderData models.Order
		result = db.Where("id = ?", id).First(&orderData)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pesanan tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Mendapatkan detail properti yang sesuai dengan pesanan
		var propertyData models.Property
		result = db.Where("id = ?", orderData.PropertiId).First(&propertyData)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Properti tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Memeriksa apakah pemilik properti yang sedang masuk adalah pemilik properti dari pesanan yang sesuai
		if propertyData.OwnerId != userRole.Id {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "Anda bukan pemilik properti ini dan tidak memiliki izin untuk mengupdate status pembayaran"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Memeriksa apakah status pembayaran sudah "true"
		if orderData.PaymentStatus == true {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: "Status pembayaran sudah true dan tidak dapat diupdate lagi"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mengupdate status pembayaran
		orderData.PaymentStatus = true
		db.Save(&orderData)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"OrderData": orderData,
			"code":      http.StatusOK,
			"error":     false,
			"message":   "Status pembayaran berhasil diupdate",
		})
	}
}
