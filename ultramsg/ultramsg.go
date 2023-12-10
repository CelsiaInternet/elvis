package ultramsg

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/strs"
	_ "github.com/joho/godotenv/autoload"
)

func SendWhatsApp(country, phone, message string) (e.Json, error) {
	ultramsgID := envar.EnvarStr("", "ULTRAMSG_ID")
	ultramsgToken := envar.EnvarStr("", "ULTRAMSG_TOKEN")
	apiurl := strs.Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := strs.Format(`%s%s`, country, phone)
	body := message
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("body", body)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiurl, params)
	if err != nil {
		return e.Json{}, err
	}

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

/**
* .jpg
**/
func SendWhatsAppImage(country, phone, image, caption string) (e.Json, error) {
	ultramsgID := envar.EnvarStr("", "ULTRAMSG_ID")
	ultramsgToken := envar.EnvarStr("", "ULTRAMSG_TOKEN")
	apiurl := strs.Format(`https://api.ultramsg.com/%s/messages/image`, ultramsgID)
	to := strs.Format(`%s%s`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("image", image)
	data.Set("caption", caption)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiurl, params)
	if err != nil {
		return e.Json{}, err
	}

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

/**
* .webp
**/
func SendWhatsAppSticker(country, phone, sticker string) (e.Json, error) {
	ultramsgID := envar.EnvarStr("", "ULTRAMSG_ID")
	ultramsgToken := envar.EnvarStr("", "ULTRAMSG_TOKEN")
	apiurl := strs.Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := strs.Format(`%s%s`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("sticker", sticker)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiurl, params)
	if err != nil {
		return e.Json{}, err
	}

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

/**
* .pdf
**/
func SendWhatsAppDocument(country, phone, filename, document, caption string) (e.Json, error) {
	ultramsgID := envar.EnvarStr("", "ULTRAMSG_ID")
	ultramsgToken := envar.EnvarStr("", "ULTRAMSG_TOKEN")
	apiurl := strs.Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := strs.Format(`%s%s`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("filename", filename)
	data.Set("document", document)
	data.Set("caption", caption)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiurl, params)
	if err != nil {
		return e.Json{}, err
	}

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

/**
* .mp3
**/
func SendWhatsAppAudio(country, phone, audio string) (e.Json, error) {
	ultramsgID := envar.EnvarStr("", "ULTRAMSG_ID")
	ultramsgToken := envar.EnvarStr("", "ULTRAMSG_TOKEN")
	apiurl := strs.Format(`https://api.ultramsg.com/%s/messages/chat`, ultramsgID)
	to := strs.Format(`%s%s`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("audio", audio)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiurl, params)
	if err != nil {
		return e.Json{}, err
	}

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

/**
* .ogg
**/
func SendWhatsAppVoice(country, phone, audio string) (e.Json, error) {
	ultramsgID := envar.EnvarStr("", "ULTRAMSG_ID")
	ultramsgToken := envar.EnvarStr("", "ULTRAMSG_TOKEN")
	apiurl := strs.Format(`https://api.ultramsg.com/%s/messages/voice`, ultramsgID)

	to := strs.Format(`%s%s`, country, phone)
	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("audio", audio)

	client := &http.Client{}
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", apiurl, params)
	if err != nil {
		return e.Json{}, err
	}

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

/**
* .mp4
**/
func SendWhatsAppVideo(country, phone, video, caption string) (e.Json, error) {
	ultramsgID := envar.EnvarStr("", "ULTRAMSG_ID")
	ultramsgToken := envar.EnvarStr("", "ULTRAMSG_TOKEN")
	apiurl := strs.Format(`https://api.ultramsg.com/%s/messages/video`, ultramsgID)
	to := strs.Format(`%s%s`, country, phone)

	data := url.Values{}
	data.Set("token", ultramsgToken)
	data.Set("to", to)
	data.Set("video", video)
	data.Set("caption", caption)

	payload := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", apiurl, payload)
	if err != nil {
		return e.Json{}, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return e.Json{}, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	result := e.Json{
		"to":     to,
		"status": res.Status,
		"body":   string(body),
	}

	console.Log(result.ToString())

	return result, nil
}
