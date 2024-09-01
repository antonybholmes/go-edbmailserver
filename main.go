package main

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-mailer"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	subscriber := redisClient.Subscribe(ctx, "email")

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
	}
}
