package controllers

import (
	"fmt"
	"os"
	"testing"
)

func TestSendMail(t *testing.T) {
	SMTP_HOST := os.Getenv("DB_HOST")
	SMTP_PORT := os.Getenv("DB_PORT")

	fmt.Println("SMTP HERE ", SMTP_HOST+":"+SMTP_PORT)

	recipent_list := []string{"rohanchavan@nimapinfotech.com"}
	message := "Hello bro how are you"
	SendMail(recipent_list, message, "test")
}
