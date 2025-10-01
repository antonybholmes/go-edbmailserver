package main

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-edbmailserver/consts"
	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
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

	log.Debug().Msgf("%s %s", consts.AppId, consts.RedisAddr)

	subscriber := rdb.Subscribe(ctx, mailserver.EmailQueueChannel)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		err = pool.Submit(func() {
			err := processMessage(msg)

			if err != nil {
				log.Debug().Msgf("error sending email: %v", err)
			}
		})

		if err != nil {
			log.Debug().Msgf("failed to submit task: %v", err)
		}

	}
}

func processMessage(msg *redis.Message) error {
	var m mailserver.MailItem

	err := json.Unmarshal([]byte(msg.Payload), &m)

	if err != nil {
		log.Debug().Msgf("email error")
	}

	//log.Debug().Msgf("email %s %v", msg.Payload, qe.EmailType)

	return edbmail.SendEmail(&m)
}
