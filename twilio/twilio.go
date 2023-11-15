package twilio

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/cgalvisleon/elvis/envar"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
	_ "github.com/joho/godotenv/autoload"
)

func SendWhatsApp(country, phone, message string) (e.Json, error) {
	twilioSID := envar.EnvarStr("", "TWILIO_SID")
	twilioAUT := envar.EnvarStr("", "TWILIO_AUT")
	twilioFrom := envar.EnvarStr("", "TWILIO_FROM")

	apiURL := utility.Format(`https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json`, twilioSID)
	from := utility.Format(`whatsapp:%s`, twilioFrom)
	body := message
	to := utility.Format(`whatsapp:%s%s`, country, phone)

	data := url.Values{}
	data.Set("From", from)
	data.Set("Body", body)
	data.Set("To", to)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiURL, params)
	if err != nil {
		return e.Json{}, err
	}

	req.SetBasicAuth(twilioSID, twilioAUT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return e.Json{}, err
	}

	defer resp.Body.Close()

	return e.Json{
		"status": resp.Status,
		"body":   resp.Body,
	}, nil
}
