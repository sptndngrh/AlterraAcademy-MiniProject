package controllers

import (
	"html/template"
	"net/http"
	"os"
	"sewakeun_project/helpers"
	"sewakeun_project/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func VerifyEmail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")

		var user models.User
		result := db.Where("jwt_token_verify = ?", token).First(&user)
		if result.Error != nil {
			return c.String(http.StatusUnauthorized, "Gagal mengonfirmasi email")
		}

		user.DoneVerify = true
		user.JWTTokenVerify = "" // Setelah verifikasi, hapus token verifikasi
		db.Save(&user)

		// Determine the user type based on the OwnerRole field
		userType := "user" // Default to "user"
		if user.OwnerRole {
			userType = "owner"
		}

		// Define a map for user types and corresponding template file paths
		templatePaths := map[string]string{
			"user":  "helpers/html/verification_user.html",
			"owner": "helpers/html/verification_owner.html",
		}

		// Check if the user type is valid
		templatePath, ok := templatePaths[userType]
		if !ok {
			return c.String(http.StatusInternalServerError, "Invalid user type")
		}

		// Attempt to read the template file
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to load template: "+err.Error())
		}

		// Eksekusi template dan kirimkan sebagai respons
		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server Internal Sedang Error")
		}

		return nil
	}
}

func sendPaymentConfirmationEmail(recipientEmail string, orderData models.Order) error {
	// Konfigurasi email
	email := gomail.NewMessage()
	email.SetHeader("From", "your-email@gmail.com") // Ganti dengan alamat email Anda
	email.SetHeader("To", recipientEmail)
	email.SetHeader("Subject", "Konfirmasi Pembayaran")

	// Membuat tautan WhatsApp
	whatsappNumber := "+6289660515237" // Ganti dengan nomor WhatsApp yang sesuai
	whatsappURL := "https://wa.me/" + whatsappNumber

	// Sisipkan konten HTML ke dalam email
	emailBody := `
	<html>
		<head>
			<style>
				body {
					font-family: 'Roboto', sans-serif;
                	background-color: #f5f5f5;
					font-size: 16px;
					margin: 0;
					padding: 0;
				}

				.container {
					background-color: #f4f4f4;
					padding: 20px;
				}

				.content {
					background-color: #ffffff;
					padding: 20px;
				}

				h1 {
					color: #333;
					font-size: 24px;
				}

				ul {
					list-style: none;
					padding: 0;
				}

				li {
					margin-bottom: 10px;
				}

				.whatsapp-button {
					display: inline-block;
					background-color: #25d366;
					color: #fff;
					font-size: 16px;
					text-decoration: none;
					padding: 10px 20px;
					border-radius: 5px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="content">
					<h1>Terima kasih atas pembayaran Anda</h1>
					<p>Detail pesanan:</p>
					<ul>
						<li>ID Pesanan: ` + strconv.Itoa(int(orderData.Id)) + `</li>
						<li>Total Harga: Rp ` + strconv.Itoa(int(orderData.PaymentTotal)) + `</li>
						<li>Status Pembayaran: ` + helpers.GetPaymentStatusString(orderData.PaymentStatus) + `</li>
					</ul>
					<p>Untuk mengirimkan bukti transaksi, silakan hubungi kami di <a href="` + whatsappURL + `" class="whatsapp-button">Hubungi via WhatsApp</a> dan lampirkan bukti transaksi Anda.</p>
				</div>
			</div>
		</body>
	</html>
`

	email.SetBody("text/html", emailBody)

	// Konfigurasi SMTP menggunakan variabel lingkungan
	smtpHost := os.Getenv("SMTPSERVER")
	smtpPortStr := os.Getenv("SMTPPORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}
	smtpUsername := os.Getenv("SMTPUSERNAME")
	smtpPassword := os.Getenv("SMTPPASSWORD")

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Kirim email
	if err := dialer.DialAndSend(email); err != nil {
		return err
	}

	return nil
}
