package aws

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/envar"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
)

func UploaderS3(bucket, filename, contentType string, contentFile []byte) (*s3manager.UploadOutput, error) {
	sess := AwsSession()
	uploader := s3manager.NewUploader(sess)

	input := &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(contentFile),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	}

	result, err := uploader.UploadWithContext(context.Background(), input)
	if err != nil {
		console.Error(err)
	}

	return result, err
}

func DeleteS3(bucket, key string) (*s3.DeleteObjectOutput, error) {
	console.Log("Bucke:", bucket)
	console.Log("key:", key)

	sess := AwsSession()
	s3client := s3.New(sess)

	request := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	result, err := s3client.DeleteObject(request)
	if err != nil {
		console.Error(err)
	}

	return result, err
}

func DownloadS3(bucket, key string) (*s3.GetObjectOutput, error) {
	sess := AwsSession()
	s3client := s3.New(sess)

	request := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	result, err := s3client.GetObject(request)
	if err != nil {
		console.Error(err)
	}

	return result, err
}

/**
*
**/
func UploaderFile(r *http.Request, folder, name string) (Json, error) {
	r.ParseMultipartForm(2000)
	file, fileInfo, err := r.FormFile("myFile")
	if err != nil {
		return Json{}, err
	}
	defer file.Close()

	contentType := fileInfo.Header.Get("Content-Type")
	ext := ExtencionFile(fileInfo.Filename)
	filename := fileInfo.Filename
	if len(name) > 0 {
		filename = Format(`%s.%s`, name, ext)
	}
	if len(folder) > 0 {
		filename = Format(`%s/%s`, folder, filename)
	}

	storageType := EnvarStr("", "STORAGE_TYPE")
	bucket := EnvarStr("", "BUCKET")
	if storageType == "S3" {
		contentFile, err := io.ReadAll(file)
		if err != nil {
			return Json{}, err
		}

		output, err := UploaderS3(bucket, filename, contentType, contentFile)
		if err != nil {
			return Json{}, err
		}

		return Json{
			"url": output.Location,
		}, nil
	} else {
		MakeFolder(bucket)
		filename := Format(`%s/%s`, bucket, filename)

		output, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return Json{}, err
		}
		defer output.Close()

		_, err = io.Copy(output, file)
		if err != nil {
			return Json{}, err
		}

		hostname := EnvarStr("", "HOST")
		url := Format(`%s/%s`, hostname, filename)

		return Json{
			"url": url,
		}, nil
	}
}

func UploaderB64(b64, filename, contentType string) (Json, error) {
	if !ValidStr(b64, 0, []string{""}) {
		return Json{}, console.ErrorF(MSG_ATRIB_REQUIRED, "b64")
	}

	if !ValidStr(filename, 0, []string{""}) {
		return Json{}, console.ErrorF(MSG_ATRIB_REQUIRED, "filename")
	}

	if !ValidStr(contentType, 0, []string{""}) {
		return Json{}, console.ErrorF(MSG_ATRIB_REQUIRED, "content-type")
	}

	storageType := EnvarStr("", "STORAGE_TYPE")
	bucket := EnvarStr("", "BUCKET")
	if storageType == "S3" {
		contentFile, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return Json{}, err
		}

		output, err := UploaderS3(bucket, filename, contentType, contentFile)
		if err != nil {
			return Json{}, err
		}

		return Json{
			"url": output.Location,
		}, nil
	} else {
		MakeFolder(bucket)
		filename := Format(`%s/%s`, bucket, filename)
		dec, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			panic(err)
		}

		output, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return Json{}, err
		}
		defer output.Close()

		if _, err := output.Write(dec); err != nil {
			return Json{}, err
		}

		if err := output.Sync(); err != nil {
			return Json{}, err
		}

		hostname := EnvarStr("", "HOST")
		url := Format(`%s/%s`, hostname, filename)

		return Json{
			"url": url,
		}, nil
	}
}

func DeleteFile(url string) (bool, error) {
	storageType := EnvarStr("", "STORAGE_TYPE")
	if storageType == "S3" {
		bucket := EnvarStr("", "BUCKET")
		key := strings.ReplaceAll(url, Format(`https://%s.s3.amazonaws.com/`, bucket), ``)
		_, err := DeleteS3(bucket, key)
		if err != nil {
			return false, err
		}

		return true, nil
	} else {
		outdel, err := RemoveFile(url)
		if err != nil {
			return false, err
		}

		return outdel, nil
	}
}

func DownloaderFile(url string) (string, error) {
	storageType := EnvarStr("", "STORAGE_TYPE")
	if storageType == "S3" {
		bucket := EnvarStr("", "BUCKET")
		key := strings.ReplaceAll(url, Format(`https://%s.s3.amazonaws.com/`, bucket), ``)
		_, err := DownloadS3(bucket, key)
		if err != nil {
			return "", err
		}

		return url, nil
	} else {

		return url, nil
	}
}
