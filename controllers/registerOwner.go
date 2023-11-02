package controllers

import (
	"errors"
	"net/http"
	"sewakeun_project/helpers"
	"sewakeun_project/middlewares"
	"sewakeun_project/models"
	"sewakeun_project/response"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterOwner(db *gorm.DB, ownerSecretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var owner models.Owner
		if err := c.Bind(&owner); err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		var existingOwner models.Owner
		result := db.Where("username = ?", owner.Username).First(&existingOwner)
		if result.Error == nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Nama pengguna sudah ada",
			}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Gagal dalam mengecek username",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		result = db.Where("email = ?", owner.Email).First(&existingOwner)
		if result.Error == nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Email sudah ada",
			}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Gagal dalam mengecek email",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		result = db.Where("no_telp = ?", owner.NoTelp).First(&existingOwner)
		if result.Error == nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Nomor telepon sudah ada",
			}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Gagal dalam mengecek nomor telepon",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(owner.Password), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Gagal dalam menghash password",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		uniqueToken := helpers.GenerateUniqueToken()
		owner.OwnerTokenVerify = uniqueToken

		owner.Password = string(hashedPassword)
		db.Create(&owner)
		owner.Password = ""

		tokenString, err := middlewares.OwnerGenerateToken(owner.Username, ownerSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Gagal dalam mengenerate token",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		if err := helpers.OwnerSendWelcomeEmail(owner.Email, owner.Nama, uniqueToken); err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Gagal dalam mengirim email"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		return c.JSON(
			http.StatusOK,
			map[string]interface{}{"code": http.StatusOK,
				"error":   false,
				"message": "Pengguna berhasil didaftarkan, Silakan cek email untuk verifikasi lebih lanjut",
				"token":   tokenString,
				"id":      owner.Id})
	}
}
