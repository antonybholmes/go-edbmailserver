package main

import (
	"fmt"
	"io"
	"net/mail"
	"os"

	"github.com/antonybholmes/go-edbmailserver/consts"
	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/sesmailserver"
	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLogger() {
	fileLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", consts.AppName),
		MaxSize:    10,   // Max size in MB before rotating
		MaxBackups: 3,    // Keep 3 backup files
		MaxAge:     7,    // Retain files for 7 days
		Compress:   true, // Compress old log files
	}

	multiWriter := io.MultiWriter(os.Stderr, fileLogger)

	logger := zerolog.New(multiWriter).With().Timestamp().Logger()

	// we use != development because it means we need to set the env variable in order
	// to see debugging work. The default is to assume production, in which case we use
	// lumberjack
	if os.Getenv("APP_ENV") != "development" {
		//zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)).With().Timestamp().Logger()
	}

	log.Logger = logger
}

func init() {
	initLogger()

	env.Ls()

	from := sys.Must(mail.ParseAddress(env.GetStr("SMTP_FROM", "")))

	sesmailserver.InitSesMailer(from)
}

func main() {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	// make a thread pool
	pool := sys.Must(ants.NewPool(10))

	//ConsumeRedis(pool)
	ConsumeSQS(pool)
	//consumeKafka(pool)
}

// func consumeKafka(pool *ants.Pool) {

// 	var ctx = context.Background()

// 	r := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:   []string{"localhost:9094"}, // Kafka broker
// 		Topic:     mailserver.QUEUE_EMAIL_CHANNEL, // Topic name
// 		Partition: 0,
// 		MaxBytes:  10e6, // 10MB
// 	})

// 	defer r.Close()

// 	var qe mailserver.QueueEmail

// 	for {
// 		m, err := r.ReadMessage(ctx)

// 		if err != nil {
// 			fmt.Printf("%s", err)
// 			break
// 		}

// 		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

// 		err = json.Unmarshal([]byte(m.Value), &qe)

// 		if err != nil {
// 			log.Debug().Msgf("email error")
// 		}

// 		sendEmail(&qe, pool)
// 	}
// }

func sendEmailUsingPool(m *mailserver.MailItem, pool *ants.Pool) {
	pool.Submit(func() { sendEmail(m) })
}

func sendEmail(m *mailserver.MailItem) {

	//log.Debug().Msgf("send email %s %s", m.To, m.EmailType)

	switch m.EmailType {
	case edbmail.EmailQueueTypeVerify:
		edbmail.SendVerifyEmail(m)
	case edbmail.EmailQueueTypeVerified:
		edbmail.SendVerifiedEmail(m)
	case edbmail.EmailQueueTypePasswordless:
		edbmail.SendPasswordlessSigninEmail(m)
	case edbmail.EmailQueueTypePasswordReset:
		edbmail.SendPasswordResetEmail(m)
	case edbmail.EmailQueueTypePasswordUpdated:
		edbmail.SendPasswordUpdatedEmail(m)
	case edbmail.EmailQueueTypeEmailReset:
		edbmail.SendEmailResetEmail(m)
	case edbmail.EmailQueueTypeEmailUpdated:
		edbmail.SendEmailUpdatedEmail(m)
	case edbmail.EmailQueueTypeAccountCreated:
		edbmail.SendAccountCreatedEmail(m)
	case edbmail.EmailQueueTypeAccountUpdated:
		edbmail.SendAccountUpdatedEmail(m)
	case edbmail.EmailQueueTypeOTP:
		edbmail.SendOTPEmail(m)
	default:
		log.Debug().Msgf("invalid email type: %s", m.EmailType)
	}

}
