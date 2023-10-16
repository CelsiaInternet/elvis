package ultramsg

import (
	"bytes"
	"net/http"
	"net/url"

	. "github.com/cgalvisleon/elvis/envar"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
	_ "github.com/joho/godotenv/autoload"
)

func SendWhatsapp(country, phone, message string) (Json, error) {
	ultramsgID := EnvarStr("ULTRAMSG_ID")
	ultramsgToken := EnvarStr("ULTRAMSG_TOKEN")
	apiURL := Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := Format(`%s%s:@g.us`, country, phone)
	body := message
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("body", body)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiURL, params)
	if err != nil {
		return Json{}, err
	}

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
