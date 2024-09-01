package main

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"
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
	log.Debug().Msgf("go-edb-server-mailer %s", consts.REDIS_ADDR)

	subscriber := rdb.Subscribe(ctx, "email")

	var email mailer.RedisQueueEmail

	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(msg.Payload), &email)

		if err != nil {
			log.Debug().Msgf("email error")
		}

		log.Debug().Msgf("email %v", email)
	}
}
