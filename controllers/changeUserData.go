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

func ChangeNameUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")

		var user models.User
		result := db.Where("username = ?", user.Username).First(&user)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pengguna tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		if userID != fmt.Sprint(user.Id) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "Akses ditolak"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		var req models.ChangeUsernameRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		user.Username = req.NewUsername
		db.Save(&user)

		// Kembalikan data pengguna yang diperbarui
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Nama berhasil diperbarui",
			"user":    user, // Sertakan data pengguna yang diperbarui dalam respons
		})
	}
}

func ChangeUsername(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi tidak ditemukan"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		authParts := strings.SplitN(tokenString, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Format token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		tokenString = authParts[1]

		username, err := middlewares.VerifyTokenJWT(tokenString, secretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		userID := c.Param("id")

		var user models.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pengguna tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		if userID != fmt.Sprint(user.Id) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "Akses ditolak"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		var req models.ChangeUsernameRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		user.Username = req.NewUsername
		db.Save(&user)

		// Kembalikan data pengguna yang diperbarui
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Username berhasil diperbarui",
			"user":    user, // Sertakan data pengguna yang diperbarui dalam respons
		})
	}
}

func ChangePassword(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi tidak ditemukan"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		authParts := strings.SplitN(tokenString, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Format token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		tokenString = authParts[1]

		username, err := middlewares.VerifyTokenJWT(tokenString, secretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		userID := c.Param("id")

		var user models.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Pengguna tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		if userID != fmt.Sprint(user.Id) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "Akses ditolak"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		var req models.ChangePasswordRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Password saat ini salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam meng-hash password baru"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		user.Password = string(hashedNewPassword)
		db.Save(&user)

		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Password berhasil diperbarui"})
	}
}
