package controllers

import (
	"fmt"
	"net/http"
	"sewakeun_project/middlewares"
	"sewakeun_project/models"
	"sewakeun_project/response"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ChangeUserNameOwner(db *gorm.DB, ownerSecretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mengambil token otorisasi dari header permintaan
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi tidak ada"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Memeriksa format token dan memperoleh username dari token
		authParts := strings.SplitN(tokenString, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Format token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}
		tokenString = authParts[1]

		// Memverifikasi token
		username, err := middlewares.OwnerVerifyToken(tokenString, ownerSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		var req models.ChangeUsernameOwnerRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Memeriksa apakah "currentUsername" dalam permintaan sesuai dengan username yang saat ini terotentikasi
		if req.CurrentUsername != username {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Silakan login terlebih dahulu untuk mengganti username"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Memeriksa username berdasarkan owner ID
		var owner models.Owner
		result := db.Where("owner_id = ?", req.CurrentUsername).Or("username = ?", req.CurrentUsername).First(&owner)
		if result.Error != nil {
			// Jika owner tidak ditemukan, mengembalikan pesan peringatan
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Owner tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Jika pemilik ditemukan, update nama pengguna dalam database
		owner.Username = req.NewUsername
		db.Save(&owner)

		// Kembalikan respons
		message := "Silakan login kembali dengan username baru, lalu bisa masuk dengan token yang sudah disediakan"
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "message": "Username berhasil diubah", "error": false, "warning": message})
	}
}

// ChangePasswordOwner mengubah kata sandi
func ChangePasswordOwner(db *gorm.DB, ownerSecretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mengambil token otorisasi dari header permintaan
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi tidak ada"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Memeriksa format token
		authParts := strings.SplitN(tokenString, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Format token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Memperoleh username dari token
		tokenString = authParts[1]

		// Memverifikasi token
		username, err := middlewares.OwnerVerifyToken(tokenString, ownerSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		ownerID := c.Param("owner_id")

		// Mengecek apakah "currentUsername" dalam permintaan sesuai dengan username yang saat ini terotentikasi
		var owner models.Owner
		result := db.Where("username = ?", username).First(&owner)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Owner tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Jika pemilik ditemukan, update kata sandi dalam database
		if ownerID != fmt.Sprint(owner.OwnerID) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "ID ini tidak berhak untuk anda ubah"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		var req models.ChangePasswordOwnerRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Memeriksa apakah "currentPassword" dalam permintaan sesuai dengan kata sandi yang saat ini terotentikasi
		err = bcrypt.CompareHashAndPassword([]byte(owner.Password), []byte(req.CurrentPassword))
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Kata sandi saat ini salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Hash kata sandi
		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal meng-hash kata sandi baru"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Update kata sandi
		owner.Password = string(hashedNewPassword)
		db.Save(&owner)

		// Kembalikan respons
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Kata sandi berhasil diperbarui, silakan cek login kembali dengan kata sandi baru"})
	}
}

// changeNameOwner mengubah nama
func ChangeNameOwner(db *gorm.DB, ownerSecretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mengambil token otorisasi dari header permintaan
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			message := "Silakan masukkan token terbaru terlebih dahulu"
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"code": http.StatusUnauthorized, "message": message, "error": true})
		}

		// Memeriksa format token dan memperoleh username dari token
		authParts := strings.SplitN(tokenString, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}
		tokenString = authParts[1]

		// Memverifikasi token
		username, err := middlewares.OwnerVerifyToken(tokenString, ownerSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		var req models.ChangeNameOwnerRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mengambil owner ID dari parameter
		ownerID := c.Param("owner_id")

		// Memeriksa username berdasarkan owner ID
		var owner models.Owner
		result := db.Where("username = ?", username).First(&owner)
		if result.Error != nil {
			// Jika owner tidak ditemukan, mengembalikan pesan peringatan
			message := "Pengguna tidak ditemukan"
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"code": http.StatusUnauthorized, "message": message, "error": true})
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		if ownerID != fmt.Sprint(owner.OwnerID) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "ID ini tidak berhak untuk anda ubah"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Update password dalam database
		owner.Nama = req.NewNama
		db.Save(&owner)

		// Mengembalikan respons
		ownerResponse := response.NewOwnerResponse(owner.Username, owner.Nama)

		// Kembalikan respons
		message := "Nama berhasil diubah"
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "message": message, "error": false, "data": ownerResponse})
	}
}
