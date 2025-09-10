package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog/log"
)

func ConsumeSQS(pool *ants.Pool) {

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Msgf("Unable to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	var qe mailer.QueueEmail

	for {
		resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            consts.SQS_QUEUE_URL,
			MaxNumberOfMessages: 5,
			WaitTimeSeconds:     10, // long polling
		})

		if err != nil {
			log.Printf("Failed to receive messages: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, message := range resp.Messages {

			err = json.Unmarshal([]byte(*message.Body), &qe)

			if err != nil {
				log.Debug().Msgf("email error")
			}

			_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      consts.SQS_QUEUE_URL,
				ReceiptHandle: message.ReceiptHandle,
			})

			if err != nil {
				log.Debug().Msgf("Failed to delete message: %v", err)
			}

			fmt.Println("Message deleted successfully")

			log.Debug().Msgf("email %v %v", message.Body, qe.EmailType)

			sendEmail(&qe, pool)
		}
	}
}
