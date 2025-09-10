package main

import (
	"net/mail"

	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/sesmailserver"
	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog/log"
)

func init() {
	from := sys.Must(mail.ParseAddress(env.GetStr("FROM", "")))

	sesmailserver.Init(from)
}

func main() {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	env.Ls()

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
// 		Topic:     mailer.QUEUE_EMAIL_CHANNEL, // Topic name
// 		Partition: 0,
// 		MaxBytes:  10e6, // 10MB
// 	})

// 	defer r.Close()

// 	var qe mailer.QueueEmail

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

func sendEmailUsingPool(qe *mailer.QueueEmail, pool *ants.Pool) {
	pool.Submit(func() { sendEmail(qe) })
}

func sendEmail(qe *mailer.QueueEmail) {

	switch qe.EmailType {
	case mailer.QUEUE_EMAIL_TYPE_VERIFY:
		SendVerifyEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_VERIFIED:
		SendVerifiedEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_PASSWORDLESS:
		SendPasswordlessSigninEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_PASSWORD_RESET:
		SendPasswordResetEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_PASSWORD_UPDATED:
		SendPasswordUpdatedEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_EMAIL_RESET:
		SendEmailResetEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_EMAIL_UPDATED:
		SendEmailUpdatedEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_ACCOUNT_CREATED:
		SendAccountCreatedEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_ACCOUNT_UPDATED:
		SendAccountUpdatedEmail(qe)
	case mailer.QUEUE_EMAIL_TYPE_TOTP:
		SendTOTPEmail(qe)
	default:
		log.Debug().Msgf("invalid email type: %s", qe.EmailType)
	}

}
