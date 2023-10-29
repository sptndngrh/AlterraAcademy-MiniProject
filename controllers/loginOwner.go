package controllers

import (
	"errors"
	"net/http"
	"sewakeun_project/models"
	"sewakeun_project/response"

	"sewakeun_project/middlewares"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginOwner(db *gorm.DB, ownerSecretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var owner models.Owner

		if err := c.Bind(&owner); err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mengecek persamaan username dalam database
		var existingOwner models.Owner
		result := db.Where("username = ?", owner.Username).First(&existingOwner)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorResponse := response.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: "Owner tidak ditemukan nih...",
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
		if !existingOwner.DoneVerify {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Email belum terverifikasi. Silakan verifikasi email Anda terlebih dahulu.",
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Bandingkan password yang diinput dengan yang dihash
		err := bcrypt.CompareHashAndPassword([]byte(existingOwner.Password), []byte(owner.Password))
		if err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Password salah",
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Generate JWT token
		tokenString, err := middlewares.OwnerGenerateToken(existingOwner.Username, ownerSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal untuk generate token"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Return the token and owner ID
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Login owner berhasil", "token": tokenString, "id": existingOwner.OwnerID})
	}
}
