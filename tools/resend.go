package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/DeepanshuMishraa/Parcel/config"
	"github.com/resend/resend-go/v3"
)

func SendEmail(cfg *config.Config, to string, subject string, body string) error {
	ctx := context.Background()
	client := resend.NewClient(cfg.RESEND_API_KEY)

	_body := fmt.Sprintf("<p>%s</p>", body)
	params := &resend.SendEmailRequest{
		From:    cfg.FROM_EMAIL,
		To:      []string{to},
		Subject: subject,
		Html:    _body,
	}

	sent, err := client.Emails.SendWithContext(ctx, params)

	if err != nil {
		return err
	}

	log.Println("Sent email with id: ", sent.Id)
	return nil

}
