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
	subject := "Selamat Datang di Sewakeun"
	baseURL := "http://localhost:8000"
	tokenVerifyLink := baseURL + "/verify/user?token=" + UserTokenVerify
	emailBody := `
    <html>
    <head>
        <link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" rel="stylesheet">
        <style>
            body {
                font-family: Arial, sans-serif;
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
                background-color: #ff6600;
                color: #fff;
                padding: 10px 20px;
                border-radius: 5px;
                text-decoration: none;
                display: block;
                text-align: center;
                margin: 20px auto;
            }
            .btn-verify-email:hover {
                background-color: #ff3300;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Selamat Datang Pengguna Baru di Sewakeun</h1>
            <div class="message">
                <p>Selamat datang, <strong>` + nama + `</strong>,</p>
                <p>Terima kasih telah bergabung di sistem kami</p>
                <p>Jika anda butuh bantuan silakan hubungi email berikut</p>
                <p><strong>Support Team:</strong> <a href="mailto:septiandin92@gmail.com">septiandin92@gmail.com</a></p>
                <a href="` + tokenVerifyLink + `" class="btn btn-verify-email">Verify Email</a>
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
