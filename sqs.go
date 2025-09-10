package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/antonybholmes/go-edb-server-mailer/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog/log"
)

const MAX_MESSAGES = 5
const WAIT_TIME_SECONDS = 10
const SLEEP_DURATION_SECONDS = 5 * time.Second

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
			MaxNumberOfMessages: MAX_MESSAGES,
			WaitTimeSeconds:     WAIT_TIME_SECONDS,
		})

		if err != nil {
			log.Printf("Failed to receive messages: %v", err)
			time.Sleep(SLEEP_DURATION_SECONDS)
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

			log.Debug().Msgf("Message %s deleted successfully\n", *message.ReceiptHandle)

			//log.Debug().Msgf("email %v %v", message.Body, qe.EmailType)

			//sendEmailUsingPool(&qe, pool)
			sendEmail(&qe)
		}
	}
}
