package consts

import (
	"os"

	"github.com/antonybholmes/go-sys/env"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const NAME = "Experiments Mail Server"
const APP_NAME = "edb-server-mailer"

var VERSION string
var COPYRIGHT string

var SESSION_SECRET string
var SESSION_NAME string
var UPDATED string
var REDIS_ADDR string
var REDIS_PASSWORD string

var URL_SIGN_IN string
var URL_VERIFY_EMAIL string
var SQS_QUEUE_URL *string

//var URL_PASSWORDLESS_SIGN_IN string

const DO_NOT_REPLY = "Please do not reply to this message. It was sent from a notification-only email address that we don't monitor."

func init() {

	env.Load("consts.env")
	env.Load("version.env")

	VERSION = os.Getenv("VERSION")
	UPDATED = os.Getenv("UPDATED")
	COPYRIGHT = os.Getenv("COPYRIGHT")
	URL_VERIFY_EMAIL = os.Getenv("URL_VERIFY_EMAIL")

	REDIS_ADDR = os.Getenv("REDIS_ADDR")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")

	SQS_QUEUE_URL = aws.String(os.Getenv("SQS_QUEUE_URL"))

	URL_SIGN_IN = os.Getenv("URL_SIGN_IN")

	//URL_PASSWORDLESS_SIGN_IN = os.Getenv("URL_PASSWORDLESS_SIGN_IN")

	//JWT_PRIVATE_KEY = []byte(os.Getenv("JWT_SECRET"))
	//JWT_PUBLIC_KEY = []byte(os.Getenv("JWT_SECRET"))
	SESSION_SECRET = os.Getenv("SESSION_SECRET")
	SESSION_NAME = os.Getenv("SESSION_NAME")

}
