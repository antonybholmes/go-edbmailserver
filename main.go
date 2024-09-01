package main

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/mailserver"
	"github.com/antonybholmes/go-sys/env"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     consts.REDIS_ADDR,
	Password: "", // no password set
	DB:       0,  // use default DB
})

func main() {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	env.Ls()

	mailserver.Init()

	log.Debug().Msgf("edb-server-mailer %s", consts.REDIS_ADDR)

	subscriber := rdb.Subscribe(ctx, "email")

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

		switch qe.EmailType {
		case mailer.REDIS_EMAIL_TYPE_PASSWORDLESS:
			go SendPasswordlessSigninEmail(&qe)
		default:
			log.Debug().Msgf("invalid email type: %s", qe.EmailType)
		}

		log.Debug().Msgf("email %v", qe)
	}
}
