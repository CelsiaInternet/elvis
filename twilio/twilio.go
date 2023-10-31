package twilio

import (
	"bytes"
	"net/http"
	"net/url"

	. "github.com/cgalvisleon/elvis/envar"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
	_ "github.com/joho/godotenv/autoload"
)

func SendWhatsApp(country, phone, message string) (Json, error) {
	twilioSID := EnvarStr("", "TWILIO_SID")
	twilioAUT := EnvarStr("", "TWILIO_AUT")
	twilioFrom := EnvarStr("", "TWILIO_FROM")

	apiURL := Format(`https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json`, twilioSID)
	from := Format(`whatsapp:%s`, twilioFrom)
	body := message
	to := Format(`whatsapp:%s%s`, country, phone)

	data := url.Values{}
	data.Set("From", from)
	data.Set("Body", body)
	data.Set("To", to)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiURL, params)
	if err != nil {
		return Json{}, err
	}

	req.SetBasicAuth(twilioSID, twilioAUT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return Json{}, err
	}

	defer resp.Body.Close()

	return Json{
		"status": resp.Status,
		"body":   resp.Body,
	}, nil
}
