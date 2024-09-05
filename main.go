package main

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-sys/env"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func main() {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	env.Ls()

	var ctx = context.Background()
	log.Debug().Msgf("start rdb %s", consts.REDIS_ADDR)

	var rdb = redis.NewClient(&redis.Options{
		Addr:     consts.REDIS_ADDR,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//mailserver.Init()

	log.Debug().Msgf("edb-server-mailer %s", consts.REDIS_ADDR)

	subscriber := rdb.Subscribe(ctx, mailer.REDIS_EMAIL_CHANNEL)

	var qe mailer.RedisQueueEmail

	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(msg.Payload), &qe)

		if err != nil {
			log.Debug().Msgf("email error")
		}

		log.Debug().Msgf("email this %s", msg.Payload)

		switch qe.EmailType {
		case mailer.REDIS_EMAIL_TYPE_VERIFY:
			go SendVerifyEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_VERIFIED:
			go SendVerifiedEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_PASSWORDLESS:
			go SendPasswordlessSigninEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_PASSWORD_RESET:
			go SendPasswordResetEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_PASSWORD_UPDATED:
			go SendPasswordUpdatedEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_EMAIL_RESET:
			go SendEmailResetEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_EMAIL_UPDATED:
			go SendEmailUpdatedEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_ACCOUNT_CREATED:
			go SendAccountCreatedEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_ACCOUNT_UPDATED:
			go SendAccountUpdatedEmail(&qe)
		default:
			log.Debug().Msgf("invalid email type: %s", qe.EmailType)
		}

	}
}
