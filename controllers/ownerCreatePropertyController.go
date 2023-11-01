package controllers

import (
	"net/http"
	"sewakeun_project/middlewares"
	"sewakeun_project/models"
	"sewakeun_project/response"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func OwnerCreateProperty(db *gorm.DB, ownerSecretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Autentikasi dan verifikasi token
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			// Peringatan jika token tidak ada
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi tidak ada"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		authParts := strings.SplitN(tokenString, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			// Peringatan jika format token tidak valid
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Format token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		tokenString = authParts[1]

		// Verifikasi token pemilik
		username, err := middlewares.OwnerVerifyToken(tokenString, ownerSecretKey)
		if err != nil {
			// Peringatan jika token tidak valid
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Mendapatkan owner_id dari URL
		ownerID, err := strconv.Atoi(c.Param("owner_id"))
		if err != nil {
			// Peringatan jika owner ID tidak valid
			errorResponse := response.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "ID pemilik tidak valid",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mendapatkan dan memproses data properti
		var req models.Property
		if err := c.Bind(&req); err != nil {
			// Peringatan jika terjadi kesalahan dalam pemrosesan permintaan
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mencari pemilik
		var owner models.Owner
		result := db.Where("owner_id = ?", ownerID).First(&owner)
		if result.Error != nil {
			// Peringatan jika pemilik tidak ditemukan
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pemilik tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Verifikasi bahwa pemilik diotorisasi untuk membuat properti
		if username != owner.Username {
			// Peringatan jika pemilik tidak memiliki hak akses ke properti ini
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Pemilik tidak memiliki hak akses ke properti ini"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Mengambil harga properti sebagai integer
		harga := req.Harga

		// Membuat properti dengan ID pemilik
		propertyCreated := models.Property{
			Judul:     req.Judul,
			Tipe:      req.Tipe,
			Harga:     harga,
			Lokasi:    req.Lokasi,
			Deskripsi: req.Deskripsi,
			OwnerID:   owner.OwnerID,
		}

		if err := db.Create(&propertyCreated).Error; err != nil {
			// Peringatan jika gagal membuat properti
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal membuat properti"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Respon sukses
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":         http.StatusOK,
			"error":        false,
			"message":      "Properti berhasil dibuat",
			"propertyData": propertyCreated,
		})
	}
}
