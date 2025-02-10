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
	log.Debug().Msgf("start rdb %s", consts.REDIS_PASSWORD)

	var rdb = redis.NewClient(&redis.Options{
		Addr:     consts.REDIS_ADDR,
		Username: "edb",
		Password: consts.REDIS_PASSWORD, // no password set
		DB:       0,                     // use default DB
	})

	//mailserver.Init()

	log.Debug().Msgf("%s %s", consts.APP_NAME, consts.REDIS_ADDR)

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

		log.Debug().Msgf("email this %s %v", msg.Payload, qe.EmailType)

		switch qe.EmailType {
		case mailer.REDIS_EMAIL_TYPE_VERIFY:
			SendVerifyEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_VERIFIED:
			go SendVerifiedEmail(&qe)
		case mailer.REDIS_EMAIL_TYPE_PASSWORDLESS:
			err := SendPasswordlessSigninEmail(&qe)
			if err != nil {
				log.Debug().Msgf("%s", err)
			}
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
