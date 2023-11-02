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

func Register(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if err := c.Bind(&user); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		var existingUser models.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error == nil {
			errorResponse := response.ErrorResponse{Code: http.StatusConflict, Message: "Nama pengguna sudah terdaftar"}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam mengecek username"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		result = db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error == nil {
			errorResponse := response.ErrorResponse{Code: http.StatusConflict, Message: "Email sudah terdaftar"}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam mengecek email"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		result = db.Where("no_telp = ?", user.NoTelp).First(&existingUser)
		if result.Error == nil {
			errorResponse := response.ErrorResponse{Code: http.StatusConflict, Message: "Nomor telepon sudah terdaftar"}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam mengecek nomor telepon"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam hash password"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		uniqueToken := helpers.GenerateUniqueToken()
		user.JWTTokenVerify = uniqueToken

		user.Password = string(hashedPassword)
		db.Create(&user)
		user.Password = ""

		tokenString, err := middlewares.GenerateToken(user.Username, secretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal dalam mengirim token"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		if err := helpers.SendWelcomeEmail(user.Email, user.Nama, uniqueToken); err != nil {
			errorResponse := response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Gagal dalam mengirim email"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		response := map[string]interface{}{
			"id":      user.Id,
			"code":    http.StatusOK,
			"message": "Sukses melakukan registrasi, silakan cek email untuk verifikasi",
			"token":   tokenString,
		}

		return c.JSON(http.StatusOK, response)
	}
}
