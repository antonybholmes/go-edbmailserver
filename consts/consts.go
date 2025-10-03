package consts

import (
	"os"

	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	AppName        = "Experiments Mail Server"
	AppId          = "edb-mail-server"
	EmailFrom      = "Experiments Team <donotreply@rdf-lab.org>"
	ProductName    = "Experiments"
	TextDoNotReply = "Please do not reply to this email as we are not able to respond to messages sent to this address." //"Please do not reply to this message. It was sent from a notification-only email address that we don't monitor."
)

var (
	Version sys.VersionInfo

	Updated       string
	RedisAddr     string
	RedisPassword string

	UrlSignIn      string
	UrlVerifyEmail string
	SqsQueueURL    *string
)

func init() {

	env.Load("consts.env")
	//env.Load("version.env")

	//Version = os.Getenv("VERSION")
	//Updated = os.Getenv("UPDATED")
	//Copyright = os.Getenv("COPYRIGHT")
	UrlVerifyEmail = os.Getenv("URL_VERIFY_EMAIL")

	RedisAddr = os.Getenv("REDIS_ADDR")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	SqsQueueURL = aws.String(os.Getenv("SQS_QUEUE_URL"))

	UrlSignIn = os.Getenv("URL_SIGN_IN")

	Version = sys.Must(sys.LoadVersionInfo("version.json"))
}
