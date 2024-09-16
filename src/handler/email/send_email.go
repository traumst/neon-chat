package email

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

type email struct {
	sender    string
	receivers []string
	subject   string
	body      string
}

func (e email) sendAs(sender string, password string) error {
	auth := smtp.PlainAuth("", sender, password, "smtp.gmail.com")

	headers := make(map[string]string)
	headers["From"] = sender
	headers["To"] = strings.Join(e.receivers, ",")
	headers["Subject"] = e.subject
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	//headers["CC"] = "cc@example.com"
	//headers["BCC"] = "bcc@example.com"
	headers["Reply-To"] = "do-not-reply@please.com"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	msg := ""
	for k, v := range headers {
		msg += k + ": " + v + "\r\n"
	}
	msg += "\r\n" + e.body

	err := smtp.SendMail("smtp.gmail.com:587", auth, sender, e.receivers, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email about [%s] to [%v], %s", e.subject, e.receivers, err.Error())
	}

	return nil
}
