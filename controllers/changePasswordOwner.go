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

func ChangePasswordOwner(db *gorm.DB, ownerSecretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token otorisasi tidak ada"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		authParts := strings.SplitN(tokenString, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Format token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		tokenString = authParts[1]

		username, err := middlewares.OwnerVerifyToken(tokenString, ownerSecretKey)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Token tidak valid"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		ownerID := c.Param("owner_id")

		var owner models.Owner
		result := db.Where("username = ?", username).First(&owner)
		if result.Error != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusNotFound, Message: "Owner tidak ditemukan"}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		if ownerID != fmt.Sprint(owner.OwnerID) {
			errorResponse := response.ErrorResponse{Code: http.StatusForbidden, Message: "ID ini tidak berhak untuk anda ubah"}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		var req models.ChangePasswordOwnerRequest
		if err := c.Bind(&req); err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		err = bcrypt.CompareHashAndPassword([]byte(owner.Password), []byte(req.CurrentPassword))
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusUnauthorized, Message: "Kata sandi saat ini salah"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Gagal meng-hash kata sandi baru"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		owner.Password = string(hashedNewPassword)
		db.Save(&owner)

		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Kata sandi berhasil diperbarui, silakan cek login kembali dengan kata sandi baru"})
	}
}
