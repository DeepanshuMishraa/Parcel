package tools

import (
	"context"
	"log"

	"github.com/DeepanshuMishraa/mini-job-queue/config"
	"github.com/resend/resend-go/v3"
)

func SendEmail(cfg *config.Config, to string) error {
	ctx := context.Background()
	client := resend.NewClient(cfg.RESEND_API_KEY)

	params := &resend.SendEmailRequest{
		From:    cfg.FROM_EMAIL,
		To:      []string{to},
		Subject: "hello world",
		Html:    "<p>Hello from job queue</p>",
	}

	sent, err := client.Emails.SendWithContext(ctx, params)

	if err != nil {
		return err
	}

	log.Println("Sent email with id: ", sent.Id)
	return nil

}
