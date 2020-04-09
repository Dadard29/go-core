package connectors

import "github.com/go-gomail/gomail"

type SmtpConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type MailConnector struct {
	smtpConfig SmtpConfig
	from       string
	dialer     *gomail.Dialer
}

func NewMailConnector(from string, smtpConfig SmtpConfig) *MailConnector {
	return &MailConnector{
		smtpConfig: smtpConfig,
		from: from,
		dialer: gomail.NewDialer(
			smtpConfig.Host,
			smtpConfig.Port,
			smtpConfig.User,
			smtpConfig.Password),
	}
}

func (mc MailConnector) SendMail (to string, subject string, html string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mc.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	if err := mc.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
