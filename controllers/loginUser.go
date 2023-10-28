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

func LoginUser(db *gorm.DB, userSecretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User

		if err := c.Bind(&user); err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mengecek persamaan username dalam database
		var existingUser models.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorResponse := response.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: "User tidak ditemukan nih...",
				}
				return c.JSON(http.StatusUnauthorized, errorResponse)
			} else {
				errorResponse := response.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Gagal dalam mengecek username",
				}
				return c.JSON(http.StatusInternalServerError, errorResponse)
			}
		}

		// Cek apakah role ini sudah verifikasi
		if !existingUser.DoneVerify {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Email belum terverifikasi. Silakan verifikasi email Anda terlebih dahulu.",
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Bandingkan password yang diinput dengan yang dihash
		err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
		if err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Password salah",
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Generate JWT token
		tokenString, err := middlewares.UserGenerateToken(existingUser.Username, userSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal untuk generate token"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Return the token and user ID
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Login user berhasil", "token": tokenString, "id": existingUser.UserID})
	}
}
