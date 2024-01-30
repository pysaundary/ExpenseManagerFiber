package utilities

import (
	"expensesManage/config"
	"fmt"
	"net/smtp"
)

func SendEmail(toEmail string, emailSubject string, emailBody string) bool {

	smtpHost := config.SMTP_HOST
	smtpPort := config.SMTP_PORT
	smtpUsername := config.SMTP_Username
	smtpPassword := config.SMTP_Password

	// Recipient email address
	to := toEmail

	// Sender email address
	from := config.FromEmail

	// Message content
	subject := emailSubject
	body := emailBody

	// Compose the email message
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, []string{to}, []byte(message))

	if err != nil {
		fmt.Println("Error sending email:", err)
		return false
	}

	fmt.Println("Email sent successfully.")
	return true
}
