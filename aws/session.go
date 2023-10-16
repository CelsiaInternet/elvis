package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	. "github.com/cgalvisleon/elvis/envar"
)

/**
* AWS Session
**/
func AwsSession() *session.Session {
	region := EnvarStr("AWS_REGION")
	id := EnvarStr("AWS_ACCESS_KEY_ID")
	secret := EnvarStr("AWS_SECRET_ACCESS_KEY")
	token := EnvarStr("AWS_SESSION_TOKEN")

	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			id,
			secret,
			token,
		),
	}))
}
