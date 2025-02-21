package main

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/panjf2000/ants"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func ConsumeRedis(pool *ants.Pool) {
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

		SendEmail(&qe, pool)
	}

}

func SendEmail(qe *mailer.QueueEmail, pool *ants.Pool) {

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
