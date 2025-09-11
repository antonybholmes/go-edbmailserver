package main

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-edbmailserver/consts"
	mailserver "github.com/antonybholmes/go-mailserver"
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

	subscriber := rdb.Subscribe(ctx, mailserver.QUEUE_EMAIL_CHANNEL)

	var qe mailserver.QueueEmail

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

		sendEmailUsingPool(&qe, pool)
	}
}
