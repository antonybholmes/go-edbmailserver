package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/antonybholmes/go-mailer"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func ConsumeKafka(pool *ants.Pool) {
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

		SendEmail(&qe, pool)
	}

}
