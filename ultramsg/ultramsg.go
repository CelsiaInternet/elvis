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

func SendWhatsApp(country, phone, message string) (Json, error) {
	ultramsgID := EnvarStr("ULTRAMSG_ID")
	ultramsgToken := EnvarStr("ULTRAMSG_TOKEN")
	apiURL := Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := Format(`%s%s`, country, phone)
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

/**
* .jpg
**/
func SendWhatsAppImage(country, phone, image, caption string) (Json, error) {
	ultramsgID := EnvarStr("ULTRAMSG_ID")
	ultramsgToken := EnvarStr("ULTRAMSG_TOKEN")
	apiURL := Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := Format(`%s%s:@g.us`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("image", image)
	data.Set("caption", caption)

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

/**
* .webp
**/
func SendWhatsAppSticker(country, phone, sticker string) (Json, error) {
	ultramsgID := EnvarStr("ULTRAMSG_ID")
	ultramsgToken := EnvarStr("ULTRAMSG_TOKEN")
	apiURL := Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := Format(`%s%s:@g.us`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("sticker", sticker)

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

/**
* .pdf
**/
func SendWhatsAppDocument(country, phone, filename, document, caption string) (Json, error) {
	ultramsgID := EnvarStr("ULTRAMSG_ID")
	ultramsgToken := EnvarStr("ULTRAMSG_TOKEN")
	apiURL := Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := Format(`%s%s:@g.us`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("filename", filename)
	data.Set("document", document)
	data.Set("caption", caption)

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

/**
* .mp3
**/
func SendWhatsAppAudio(country, phone, audio string) (Json, error) {
	ultramsgID := EnvarStr("ULTRAMSG_ID")
	ultramsgToken := EnvarStr("ULTRAMSG_TOKEN")
	apiURL := Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := Format(`%s%s:@g.us`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("audio", audio)	

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

/**
* .ogg
**/
func SendWhatsAppVoice(country, phone, audio string) (Json, error) {
	ultramsgID := EnvarStr("ULTRAMSG_ID")
	ultramsgToken := EnvarStr("ULTRAMSG_TOKEN")
	apiURL := Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := Format(`%s%s:@g.us`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("audio", audio)	

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

/**
* .mp4
**/
func SendWhatsAppVideo(country, phone, video, caption string) (Json, error) {
	ultramsgID := EnvarStr("ULTRAMSG_ID")
	ultramsgToken := EnvarStr("ULTRAMSG_TOKEN")
	apiURL := Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := Format(`%s%s:@g.us`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("video", video)
	data.Set("caption", caption)

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