package controllers

import (
	"net/http"
	"sewakeun_project/models"
	"text/template"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserVerifyEmail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")

		var user models.User
		result := db.Where("user_token_verify = ?", token).First(&user)
		if result.Error != nil {
			return c.String(http.StatusUnauthorized, "Gagal dalam verifikasi token")
		}

		user.DoneVerify = true
		user.UserTokenVerify = ""
		db.Save(&user)

		// Baca template HTML dari file
		tmpl, err := template.ParseFiles("helpers/useremailverif.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		// Eksekusi template dan kirimkan sebagai respons
		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error"+err.Error())
		}

		return nil
	}
}