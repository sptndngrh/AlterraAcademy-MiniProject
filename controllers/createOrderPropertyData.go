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

func CreateOrderingProperty(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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
		var userRole models.User
		result := db.Where("username = ?", username).First(&userRole)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pengguna tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Memeriksa peran pengguna
		if userRole.OwnerRole {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "Anda adalah pemilik properti dan tidak dapat membuat pesanan"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		var CreateOrderData struct {
			Bulanan    int  `json:"bulanan"`
			PropertiID uint `json:"properti_id"`
		}

		if err := c.Bind(&CreateOrderData); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		var propertyData models.Property
		result = db.Where("id = ?", CreateOrderData.PropertiID).First(&propertyData)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Properti tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Memeriksa apakah properti sudah tersewa
		if propertyData.StatusTersewa {
			errorResponse := response.ErrorResponse{Code: http.StatusConflict, Message: "Maaf, properti sudah disewakan, silakan cari lagi menurut pilihan anda"}
			return c.JSON(http.StatusConflict, errorResponse)
		}

		// Menghitung PaymentTotal
		bulanan := CreateOrderData.Bulanan
		hargaProperti := propertyData.Harga
		PaymentTotal := bulanan * hargaProperti

		// Menghitung start date (waktu saat ini) dan end date berdasarkan bulanan
		startDate := time.Now()
		endDate := startDate.AddDate(0, bulanan, 0)

		// Mengatur StatusTersewa properti yang dipesan menjadi true
		propertyData.StatusTersewa = true
		db.Save(&propertyData)

		// Membuat orderedProperty dari data yang dihitung
		orderData := models.Order{
			Bulanan:       bulanan,
			PaymentTotal:  PaymentTotal,
			PaymentStatus: false,
			StartDate:     startDate,
			EndDate:       endDate,
			UserId:        userRole.Id,
			PropertiId:    propertyData.Id,
		}

		// Simpan orderedProperty ke dalam database
		if err := db.Create(&orderData).Error; err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam menambahkan data properti"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Kirim email notifikasi
		if err := sendPaymentConfirmationEmail(userRole.Email, orderData); err != nil {
			// Handle kesalahan pengiriman email
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   true,
				"message": "Gagal mengirim email konfirmasi pembayaran",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"PropertyData": propertyData,
			"OrderData":    orderData,
			"code":         http.StatusOK,
			"error":        false,
			"message":      "Data property berhasil ditambahkan",
		})
	}
}
