package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"strings"
)

const (
	HOST    = "smtp.gmail.com"
	PORT    = "587"
	SUBJECT = "Food Delivery"
)

func SendEmail(logger logging.Logger, post *model.Post) {
	auth := smtp.PlainAuth("", os.Getenv("POST_FROM"), os.Getenv("POST_PASSWORD"), HOST)
	from := os.Getenv("POST_FROM")
	smtpHost := HOST
	smtpPort := PORT

	msg := fmt.Sprintf("Уважаемый клиент, Ваш текущий пароль: %s.", post.Password)
	message := strings.Replace("From: "+from+"~To: "+post.Email+"~Subject: "+SUBJECT+"~~", "~", "\r\n", -1) + msg
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{post.Email}, []byte(message))
	if err != nil {
		logger.Errorf("Error while sending email to %s:%s", post.Email, err)
		return
	}
	logger.Infof("Email for %s Sent Successfully!", post.Email)
}
