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

func ChangeUserNameUser(db *gorm.DB, userSecretKey []byte) echo.HandlerFunc {
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
		var req models.ChangeUsernameUserRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Memeriksa apakah "currentUsername" dalam permintaan sesuai dengan username yang saat ini terotentikasi
		if req.CurrentUsername != username {
			message := "Silakan login terlebih dahulu untuk mengganti username"
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"code": http.StatusUnauthorized, "message": message, "error": true})
		}

		// Memeriksa username berdasarkan user ID
		var user models.User
		result := db.Where("user_id = ?", req.CurrentUsername).First(&user)
		if result.Error != nil {
			// Jika user tidak ditemukan, mengembalikan pesan peringatan
			message := "Pengguna tidak ditemukan"
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"code": http.StatusUnauthorized, "message": message, "error": true})
		}

		// Jika pemilik ditemukan, update nama pengguna dalam database
		user.Username = req.NewUsername
		db.Save(&user)

		message := "Silakan login kembali dengan username baru, lalu bisa masuk dengan token yang sudah disediakan"
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "message": "Username berhasil diubah", "error": false, "warning": message})
	}
}