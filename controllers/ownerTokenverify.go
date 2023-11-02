package controllers

import (
	"net/http"
	"sewakeun_project/models"
	"text/template"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func OwnerVerifyEmail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")

		var owner models.Owner
		result := db.Where("owner_token_verify = ?", token).First(&owner)
		if result.Error != nil {
			return c.String(http.StatusUnauthorized, "Gagal dalam verifikasi token")
		}

		owner.DoneVerify = true
		owner.OwnerTokenVerify = ""
		db.Save(&owner)

		// Baca template HTML dari file
		tmpl, err := template.ParseFiles("helpers/html/owneremailverif.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		// Eksekusi template dan kirimkan sebagai respons
		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error" + err.Error())
		}

		return nil
	}
}
