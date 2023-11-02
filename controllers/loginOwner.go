package controllers

import (
	"errors"
	"net/http"
	"sewakeun_project/middlewares"
	"sewakeun_project/models"
	"sewakeun_project/response"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginOwner(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if err := c.Bind(&user); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mengecek apakah username ada dalam database
		var existingUser models.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Nama pengguna atau password salah"}
				return c.JSON(http.StatusUnauthorized, errorResponse)
			} else {
				errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal memeriksa username"}
				return c.JSON(http.StatusInternalServerError, errorResponse)
			}
		}

		// Membandingkan password yang dimasukkan dengan password yang di-hash
		err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Nama pengguna atau password yang di hash salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Mengecek apakah pengguna adalah owner
		if !existingUser.OwnerRole {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "Anda tidak memiliki akses untuk login sebagai pengguna"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Mengecek apakah pengguna sudah melakukan verifikasi email
		if !existingUser.DoneVerify {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Akun anda belum verified. Silahkan lakukan verifikasi email terlebih dahulu"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Generate JWT token
		tokenString, err := middlewares.GenerateToken(existingUser.Username, secretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal generate token"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Menyertakan ID pengguna dalam respons
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Sukses Login Sebagai Owner", "token": tokenString, "id": existingUser.Id})
	}
}
