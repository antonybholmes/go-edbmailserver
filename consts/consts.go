package consts

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/antonybholmes/go-sys/env"
)

var NAME string
var APP_NAME string
var APP_URL string
var VERSION string
var COPYRIGHT string

var SESSION_SECRET string
var SESSION_NAME string
var UPDATED string

const DO_NOT_REPLY = "Please do not reply to this message. It was sent from a notification-only email address that we don't monitor."

func Load() {
	env.Load()

	godotenv.Load("consts.env")
	godotenv.Load("version.env")

	NAME = os.Getenv("NAME")
	APP_NAME = os.Getenv("APP_NAME")
	APP_URL = os.Getenv("APP_URL")
	VERSION = os.Getenv("VERSION")
	UPDATED = os.Getenv("UPDATED")
	COPYRIGHT = os.Getenv("COPYRIGHT")
	//JWT_PRIVATE_KEY = []byte(os.Getenv("JWT_SECRET"))
	//JWT_PUBLIC_KEY = []byte(os.Getenv("JWT_SECRET"))
	SESSION_SECRET = os.Getenv("SESSION_SECRET")
	SESSION_NAME = os.Getenv("SESSION_NAME")

}
