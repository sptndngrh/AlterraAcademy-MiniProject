package controllers

import (
	"net/http"
	"sewakeun_project/models"
	"sewakeun_project/response"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// UpdatePaymentStatus mengizinkan pengguna untuk memperbarui PaymentStatus pesanan properti.
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

		// Mendapatkan data JSON raw body
		var requestBody struct {
			OrderID       uint `json:"id"`
			PaymentStatus bool `json:"payment_status"`
		}

		if err := c.Bind(&requestBody); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   true,
				"message": "Gagal dalam mengurai JSON body.",
			})
		}

		// Mengidentifikasi pesanan properti
		var order models.Order
		if err := db.First(&order, requestBody.OrderID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   true,
				"message": "Pesanan properti tidak ditemukan.",
			})
		}

		// Memeriksa apakah pengguna yang melakukan permintaan adalah pemilik properti
		isOwner := IsPropertyOwner(db, c, order.PropertiId)
		if !isOwner {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   true,
				"message": "Anda tidak diizinkan untuk mengubah status pembayaran.",
			})
		}

		// Memperbarui PaymentStatus pesanan properti
		order.PaymentStatus = requestBody.PaymentStatus
		if err := db.Save(&order).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   true,
				"message": "Gagal dalam memperbarui PaymentStatus.",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":   false,
			"message": "PaymentStatus pesanan properti diperbarui.",
		})
	}
}

// Fungsi untuk memeriksa apakah pengguna adalah pemilik properti berdasarkan OwnerRole
func IsPropertyOwner(db *gorm.DB, c echo.Context, propertyID uint) bool {
	var user models.User
	username := c.Get("username").(string)

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}

	// Periksa apakah pengguna memiliki OwnerRole yang sama dengan properti yang akan diubah
	return user.OwnerRole && user.Id == propertyID
}
