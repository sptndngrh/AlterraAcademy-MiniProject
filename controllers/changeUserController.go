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

// ChangeUserNameUser mengubah nama pengguna
func ChangeUserNameUser(db *gorm.DB, userSecretKey []byte) echo.HandlerFunc {
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
		username, err := middlewares.UserVerifyToken(tokenString, userSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		var req models.ChangeUsernameUserRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Memeriksa apakah "currentUsername" dalam permintaan sesuai dengan username yang saat ini terotentikasi
		if req.CurrentUsername != username {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Silakan login terlebih dahulu untuk mengganti username"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		userID := c.Param("user_id")

		// Memeriksa username berdasarkan user ID
		var user models.User
		result := db.Where("user_id = ?", req.CurrentUsername).Or("username = ?", req.CurrentUsername).First(&user)
		if result.Error != nil {
			// Jika user tidak ditemukan, mengembalikan pesan peringatan
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Owner tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Jika pemilik ditemukan, update kata sandi dalam database
		if userID != fmt.Sprint(user.UserID) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "ID ini tidak berhak untuk anda ubah"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Jika pemilik ditemukan, update nama pengguna dalam database
		user.Username = req.NewUsername
		db.Save(&user)

		// Mengembalikan respons
		message := "Silakan login kembali dengan username baru, lalu bisa masuk dengan token yang sudah disediakan"
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "message": "Username berhasil diubah", "error": false, "warning": message})
	}
}

// ChangePasswordUser mengubah kata sandi
func ChangePasswordUser(db *gorm.DB, userSecretKey []byte) echo.HandlerFunc {
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

		// Memverifikasi token
		tokenString = authParts[1]

		// Memverifikasi token
		username, err := middlewares.UserVerifyToken(tokenString, userSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		userID := c.Param("user_id")

		// Memeriksa apakah "currentUsername" dalam permintaan sesuai dengan username yang saat ini terotentikasi
		var user models.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pengguna tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Jika pemilik ditemukan, update kata sandi dalam database
		if userID != fmt.Sprint(user.UserID) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "ID ini tidak berhak untuk anda ubah"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		var req models.ChangePasswordUserRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Memeriksa apakah "currentPassword" dalam permintaan sesuai dengan kata sandi yang saat ini terotentikasi
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Kata sandi saat ini salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Meng-hash kata sandi baru
		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal meng-hash kata sandi baru"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Memperbarui kata sandi
		user.Password = string(hashedNewPassword)
		db.Save(&user)

		// Mengembalikan respons
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Kata sandi berhasil diperbarui, silakan cek login kembali dengan kata sandi baru"})
	}
}

// ChangeNameUser mengubah nama
func ChangeNameUser(db *gorm.DB, userSecretKey []byte) echo.HandlerFunc {
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
		username, err := middlewares.UserVerifyToken(tokenString, userSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Membaca data JSON yang dikirimkan dalam permintaan
		var req models.ChangeNameUserRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mengambil owner ID dari parameter
		userID := c.Param("user_id")

		// Memeriksa username berdasarkan user ID
		var user models.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			// Jika user tidak ditemukan, mengembalikan pesan peringatan
			message := "Pengguna tidak ditemukan"
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"code": http.StatusUnauthorized, "message": message, "error": true})
		}

		// Jika pemilik ditemukan, update kata sandi dalam database
		if userID != fmt.Sprint(user.UserID) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "ID ini tidak berhak untuk anda ubah"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Update password dalam database
		user.Nama = req.NewNama
		db.Save(&user)

		// Mengembalikan respons
		userResponse := response.NewUserResponse(user.Username, user.Nama)

		// Kembalikan respons
		message := "Nama berhasil diubah"
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "message": message, "error": false, "data": userResponse})
	}
}
