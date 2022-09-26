package notify

//go:generate mockgen -destination=../mocks/mock_email_notifier.go -package=mocks . Notifier

import (
	"bytes"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"kraken-dca-bot/assets"
	"kraken-dca-bot/internal/domain"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	root       = filepath.Join(filepath.Dir(b), "../..")
	newEmail   = mail.NewMSG
)

type EmailNotifier struct {
	config *domain.Config
	client *mail.SMTPServer
}

func NewEmailNotifier(config *domain.Config) Notifier {
	client := mail.NewSMTPClient()
	client.Host = config.Smtp.Host
	client.Port = config.Smtp.Port
	client.Username = config.Smtp.User
	client.Password = config.Smtp.Password
	client.Encryption = mail.EncryptionSTARTTLS

	return EmailNotifier{
		config: config,
		client: client,
	}
}

func (en EmailNotifier) NotifyFailure(transaction *domain.Transaction) error {
	t, err := template.ParseFS(assets.EmailFS, "email/transaction_failed.html")
	if err != nil {
		return fmt.Errorf("failed to parse the transaction failure template file : %w", err)
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, transaction)
	if err != nil {
		return fmt.Errorf("failed to fill the transaction failure template file : %w", err)
	}

	smtpClient, err := en.client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect the SMTP client : %w", err)
	}

	email := newEmail()
	email.SetFrom("Kraken DCA Bot <" + en.config.Smtp.From + ">").
		AddTo(en.config.Notify).
		SetSubject("Kraken DCA Bot - Transaction failure")
	email.SetBody(mail.TextHTML, buffer.String())

	if email.Error != nil {
		return email.Error
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}
