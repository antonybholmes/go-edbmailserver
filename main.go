package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/sesmailserver"
	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/panjf2000/ants"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
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

	consumeRedis(pool)
	//ConsumeKafka(pool)
}

func consumeRedis(pool *ants.Pool) {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	var ctx = context.Background()

	log.Debug().Msgf("start rdb %s", consts.REDIS_PASSWORD)

	var rdb = redis.NewClient(&redis.Options{
		Addr:     consts.REDIS_ADDR,
		Username: "edb",
		Password: consts.REDIS_PASSWORD, // no password set
		DB:       0,                     // use default DB
	})

	log.Debug().Msgf("%s %s", consts.APP_NAME, consts.REDIS_ADDR)

	subscriber := rdb.Subscribe(ctx, mailer.QUEUE_EMAIL_CHANNEL)

	var qe mailer.QueueEmail

	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(msg.Payload), &qe)

		if err != nil {
			log.Debug().Msgf("email error")
		}

		//log.Debug().Msgf("email %s %v", msg.Payload, qe.EmailType)

		sendEmail(&qe, pool)
	}
}

func consumeKafka(pool *ants.Pool) {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	var ctx = context.Background()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9094"}, // Kafka broker
		Topic:     mailer.QUEUE_EMAIL_CHANNEL, // Topic name
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	defer r.Close()

	var qe mailer.QueueEmail

	for {
		m, err := r.ReadMessage(ctx)

		if err != nil {
			fmt.Printf("%s", err)
			break
		}

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		err = json.Unmarshal([]byte(m.Value), &qe)

		if err != nil {
			log.Debug().Msgf("email error")
		}

		//log.Debug().Msgf("email %s %v", msg.Payload, qe.EmailType)

		sendEmail(&qe, pool)
	}

}

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
	default:
		log.Debug().Msgf("invalid email type: %s", qe.EmailType)
	}

}
