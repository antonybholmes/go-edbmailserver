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

	var rdb = redis.NewClient(&redis.Options{
		Addr:     consts.RedisAddr,
		Username: "edb",
		Password: consts.RedisPassword, // no password set
		DB:       0,                    // use default DB
	})

	log.Debug().Msgf("%s %s", consts.AppName, consts.RedisAddr)

	subscriber := rdb.Subscribe(ctx, mailserver.EmailQueueChannel)

	var m mailserver.MailItem

	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(msg.Payload), &m)

		if err != nil {
			log.Debug().Msgf("email error")
		}

		//log.Debug().Msgf("email %s %v", msg.Payload, qe.EmailType)

		sendEmailUsingPool(&m, pool)
	}
}
