package helpers

import (
	"os"
	"strconv"

	"github.com/go-gomail/gomail"
)

func UserSendWelcomeEmail(userEmail, nama, UserTokenVerify string) error {
	serverSmtp := os.Getenv("SMTPSERVER")
	portSmtp := os.Getenv("SMTPPORT")
	usernameSmtp := os.Getenv("SMTPUSERNAME")
	PasswordSmtp := os.Getenv("SMTPPASSWORD")

	sender := usernameSmtp
	recipient := userEmail
	subject := "Selamat Datang di Sistem Penyewaan Sewakeun"
	baseURL := "http://localhost:8000"
	tokenVerifyLink := baseURL + "/verify/user?token=" + UserTokenVerify
	emailBody := `
    <html>
    <head>
        <link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" rel="stylesheet">
        <style>
            @import url('https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,300;0,400;0,500;0,700;1,400;1,500;1,700&display=swap');
            body {
                font-family: 'Roboto', sans-serif;
                background-color: #f5f5f5;
            }
            .container {
                max-width: 600px;
                margin: 0 auto;
                padding: 20px;
                background-color: #fff;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                border-radius: 5px;
            }
            h1 {
                text-align: center;
                color: #333;
            }
            .message {
                background-color: #f9f9f9;
                padding: 15px;
                border: 1px solid #ddd;
            }
            p {
                font-size: 16px;
                margin-top: 10px;
            }
            strong {
                font-weight: bold;
            }
            .footer {
                text-align: center;
                margin-top: 20px;
                color: #666;
            }
            .btn-verify-email {
                background-color: #64CCC5;
                color: #fff;
                padding: 12px 18px;
                border-radius: 12px;
                text-decoration: none;
                display: block;
                text-align: center;
                margin: 20px auto;
            }
            .btn-verify-email:hover {
                background-color: #176B87;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Selamat Datang Pengguna 😊 Baru di Sewakeun</h1>
            <div class="message">
                <p>Selamat datang, <strong>` + nama + `</strong>,</p>
                <p>Terima kasih telah bergabung di sistem kami</p>
                <p>Jika anda butuh bantuan silakan hubungi email berikut</p>
                <p><strong>Support Team:</strong> <a href="mailto:septiandin92@gmail.com">septiandin92@gmail.com</a></p>
                <a href="` + tokenVerifyLink + `" class="btn btn-verify-email">Verifikasi akunmu disini...</a>
            </div>
            <div class="footer">
                <p>&copy; Sewakeun . 2023 - All Rights Reserved</p>
            </div>
        </div>
    </body>
    </html>
    `

	portSmtpStr, err := strconv.Atoi(portSmtp)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(serverSmtp, portSmtpStr, usernameSmtp, PasswordSmtp)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
