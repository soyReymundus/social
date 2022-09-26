package emailservice

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	Auth smtp.Auth
}

func (e *EmailService) Open() error {
	//e.Auth = smtp.PlainAuth("", os.Getenv("SMTPUsername"), os.Getenv("SMTPPassword"), os.Getenv("SMTPHost"))
	return nil
}

func (e *EmailService) NoReply(To []string, Subject, Message string) error {
	/*
		to := ""

		for i := 0; i < len(To); i++ {
			to += "To: " + To[i] + "\r\n"
		}

		message := []byte(to + "Subject:" + Subject + "\r\n\r\n" + Message + "\r\n")

		err := smtp.SendMail(os.Getenv("SMTPAddresses"), e.Auth, "no-reply@"+os.Getenv("SMTPHost"), To, message)

		if err != nil {
			return errors.New("Error sending email")
		}
	*/
	fmt.Println(Subject + "\n" + Message)

	return nil
}
