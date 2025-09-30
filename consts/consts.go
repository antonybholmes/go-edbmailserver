package consts

import (
	"os"

	"github.com/antonybholmes/go-sys/env"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	Name    = "Experiments Mail Server"
	AppName = "edb-mail-server"

	TextDoNotReply = "Please do not reply to this message. It was sent from a notification-only email address that we don't monitor."
)

var (
	Version   string
	Copyright string

	SessionSecret string
	SessionName   string
	Updated       string
	RedisAddr     string
	RedisPassword string

	UrlSignIn      string
	UrlVerifyEmail string
	SqsQueueURL    *string
)

func init() {

	env.Load("consts.env")
	env.Load("version.env")

	Version = os.Getenv("VERSION")
	Updated = os.Getenv("UPDATED")
	Copyright = os.Getenv("COPYRIGHT")
	UrlVerifyEmail = os.Getenv("URL_VERIFY_EMAIL")

	RedisAddr = os.Getenv("REDIS_ADDR")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	SqsQueueURL = aws.String(os.Getenv("SQS_QUEUE_URL"))

	UrlSignIn = os.Getenv("URL_SIGN_IN")

	//URL_PASSWORDLESS_SIGN_IN = os.Getenv("URL_PASSWORDLESS_SIGN_IN")

	//JWT_PRIVATE_KEY = []byte(os.Getenv("JWT_SECRET"))
	//JWT_PUBLIC_KEY = []byte(os.Getenv("JWT_SECRET"))
	SessionSecret = os.Getenv("SESSION_SECRET")
	SessionName = os.Getenv("SESSION_NAME")

}
