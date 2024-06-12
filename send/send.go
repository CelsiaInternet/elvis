package send

import (
	"fmt"
	"net/smtp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
)

func SMS(country string, mobile string, message string) (bool, interface{}, error) {
	var result bool

	phoneNumber := country + mobile
	sess := AwsSession()
	svc := sns.New(sess)

	message = strs.RemoveAcents(message)
	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phoneNumber),
	}

	output, err := svc.Publish(params)
	if err != nil {
		return result, output, console.Error(err)
	}

	return true, output, nil
}

func Email(to string, subject string, message string) (bool, error) {
	from := "sgo@celsia.com"
	port := envar.EnvarInt(25, "SMTP_PORT")
	host := envar.EnvarStr("relayappann.celsia.local", "SMTP_HOST")
	addr := strs.Format(`%s:%d`, host, port)
	c, err := smtp.Dial(addr)
	if err != nil {
		return false, err
	}

	if err := c.Mail(from); err != nil {
		return false, err
	}

	if err := c.Rcpt(to); err != nil {
		return false, err
	}

	wc, err := c.Data()
	if err != nil {
		return false, err
	}

	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s", to, from, subject, message)
	if _, err = wc.Write([]byte(msg)); err != nil {
		return false, err
	}

	if err := wc.Close(); err != nil {
		return false, err
	}

	return true, nil
}

/**
* VerifyMobile
* Send sms message a code to six digit from validate user identity
**/
func VerifyMobile(app string, device string, country string, phoneNumber string) error {
	code := utility.GetCodeVerify(6)
	cache.SetVerify(device, country+phoneNumber, code)

	message := strs.Format(msg.MSG_MOBILE_VALIDATION, app, code)
	_, _, err := SMS(country, phoneNumber, message)
	if err != nil {
		return err
	}

	return nil
}

/**
* CheckMobile
* Check code in cache db
**/
func CheckMobile(device string, country string, mobile string, code string) (bool, error) {
	val, err := cache.GetVerify(device, country+mobile)
	if err != nil {
		return false, err
	}

	result := val == code
	if result {
		cache.DelVerify(device, country+mobile)
	}

	return result, nil
}

func VerifyEmail(app, device, email string) error {
	key := strs.Format("%s-%s", device, email)
	code := utility.GetCodeVerify(6)
	cache.SetVerify(device, key, code)

	message := strs.Format(msg.MSG_MOBILE_VALIDATION, app, code)
	_, err := Email(email, "Validacion de cuenta", message)
	if err != nil {
		return err
	}

	return nil
}

func CheckEmail(device, email, code string) (bool, error) {
	key := strs.Format("%s-%s", device, email)
	val, err := cache.GetVerify(device, key)
	if err != nil {
		return false, err
	}

	result := val == code
	if result {
		cache.DelVerify(device, key)
	}

	return result, nil
}
