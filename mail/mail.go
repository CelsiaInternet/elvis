package mail

import (
	"context"
	"os"

	"github.com/cgalvisleon/elvis/cache"
	. "github.com/cgalvisleon/elvis/envar"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
	_ "github.com/joho/godotenv/autoload"
	mail "github.com/xhit/go-simple-mail/v2"
)

func Send(ctx context.Context, from string, to string, subject string, html string) (bool, error) {
	port := EnvarInt("EMAIL_PORT")
	server := mail.NewSMTPClient()
	server.Host = EnvarStr("EMAIL_HOST")
	server.Port = port
	server.Username = EnvarStr("EMAIL")
	server.Password = EnvarStr("EMAIL_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		return false, err
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(from)
	email.AddTo(to)
	email.SetSubject(subject)

	email.SetBody(mail.TextHTML, html)

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		return false, err
	}

	return true, nil
}

func SendVerify(ctx context.Context, to string, subject string, title string, email string, code string) (bool, error) {
	company := EnvarStr("COMPANY")
	fromEmail := EnvarStr("EMAIL")
	project := EnvarStr("PROJECT")
	from := Format("%s account team <%s>", Titlecase(project), fromEmail)

	css, err := os.ReadFile("./assets/template/style.txt")
	if err != nil {
		return false, err
	}

	template, err := os.ReadFile("./assets/template/mailVerify.txt")
	if err != nil {
		return false, err
	}

	html := Format(string(template), css, title, email, code, company)
	return Send(ctx, from, to, subject, html)
}

func SendAlert(ctx context.Context, to string, subject string, title string, subtitle string, message string, button string, href string, thanks string) (bool, error) {
	company := EnvarStr("COMPANY")
	fromEmail := EnvarStr("EMAIL")
	project := EnvarStr("PROJECT")
	from := Format("%s account team <%s>", Titlecase(project), fromEmail)

	css, err := os.ReadFile("./assets/template/style.txt")
	if err != nil {
		return false, err
	}

	template, err := os.ReadFile("./assets/template/mailAlert.txt")
	if err != nil {
		return false, err
	}

	html := Format(string(template), css, title, subtitle, message, href, button, thanks, company)
	return Send(ctx, from, to, subject, html)
}

func SendAction(ctx context.Context, to string, subject string, title string, message string, button string, href string) (bool, error) {
	company := EnvarStr("COMPANY")
	fromEmail := EnvarStr("EMAIL")
	project := EnvarStr("PROJECT")
	from := Format("%s account team <%s>", Titlecase(project), fromEmail)

	css, err := os.ReadFile("./assets/template/style.txt")
	if err != nil {
		return false, err
	}

	logo, err := os.ReadFile("./assets/template/logo.txt")
	if err != nil {
		return false, err
	}

	template, err := os.ReadFile("./assets/template/mailAction.txt")
	if err != nil {
		return false, err
	}

	html := Format(string(template), css, logo, title, message, href, button, company)
	return Send(ctx, from, to, subject, html)
}

func VerifyMail(ctx context.Context, device string, name string, email string) error {
	code := GetCodeVerify(6)
	cache.SetVerify(device, email, code)

	to := Format("%s <%s>", name, email)
	_, err := SendVerify(ctx, to, MSG_MAIL_001, MSG_MAIL_002, email, code)
	if err != nil {
		return err
	}

	return nil
}

func CheckMail(ctx context.Context, device string, email string, code string) (bool, error) {
	val, err := cache.GetVerify(device, email)
	if err != nil {
		return false, err
	}

	return val == code, nil
}
