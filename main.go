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

func sendEmail(qe *mailer.QueueEmail, pool *ants.Pool) {

	switch qe.EmailType {
	case mailer.QUEUE_EMAIL_TYPE_VERIFY:
		pool.Submit(func() { SendVerifyEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_VERIFIED:
		pool.Submit(func() { SendVerifiedEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_PASSWORDLESS:
		pool.Submit(func() { SendPasswordlessSigninEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_PASSWORD_RESET:
		pool.Submit(func() { SendPasswordResetEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_PASSWORD_UPDATED:
		pool.Submit(func() { SendPasswordUpdatedEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_EMAIL_RESET:
		pool.Submit(func() { SendEmailResetEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_EMAIL_UPDATED:
		pool.Submit(func() { SendEmailUpdatedEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_ACCOUNT_CREATED:
		pool.Submit(func() { SendAccountCreatedEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_ACCOUNT_UPDATED:
		pool.Submit(func() { SendAccountUpdatedEmail(qe) })
	case mailer.QUEUE_EMAIL_TYPE_TOTP:
		pool.Submit(func() { SendTOTPEmail(qe) })
	default:
		log.Debug().Msgf("invalid email type: %s", qe.EmailType)
	}

}
