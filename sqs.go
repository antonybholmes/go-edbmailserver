package main

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-edbmailserver/consts"
	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-sys"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog/log"
)

// each SQS receive can get up to 5 messages at a time
const (
	MaxMessages  = 5
	WaitTimeSecs = 10
)

type SQSConsumer struct {
	pool    *ants.Pool
	backoff *sys.Backoff
}

func NewSQSConsumer(pool *ants.Pool) *SQSConsumer {
	return &SQSConsumer{
		pool:    pool,
		backoff: sys.NewDefaultBackoff(),
	}
}

// ConsumeSQS starts consuming messages from the SQS queue.
// This is a long-running function that continuously polls the queue for messages.
func (c *SQSConsumer) ConsumeSQS() {

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Msgf("Unable to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	log.Debug().Msgf("start sqs %s", *consts.SqsQueueURL)

	for {
		resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            consts.SqsQueueURL,
			MaxNumberOfMessages: MaxMessages,
			WaitTimeSeconds:     WaitTimeSecs,
		})

		if err != nil {
			log.Printf("failed to receive messages: %v", err)
			c.backoff.Sleep()
			continue
		}

		err = c.pool.Submit(func() {
			err := processMessages(client, ctx, resp)
			if err != nil {
				log.Debug().Msgf("error processing messages: %v", err)
			}
		})

		if err != nil {
			log.Debug().Msgf("failed to submit processMessages: %v", err)
		}

	}
}

func processMessages(client *sqs.Client, ctx context.Context, resp *sqs.ReceiveMessageOutput) error {
	var m mailserver.MailItem

	for _, message := range resp.Messages {

		err := json.Unmarshal([]byte(*message.Body), &m)

		if err != nil {
			log.Debug().Msgf("error reading email json: %v", err)
		}

		handle := message.ReceiptHandle

		_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      consts.SqsQueueURL,
			ReceiptHandle: handle,
		})

		if err != nil {
			log.Debug().Msgf("failed to delete message: %v", err)
		}

		log.Debug().Msgf("message deleted: %s", *handle)

		//log.Debug().Msgf("email %v %v", message.Body, qe.EmailType)

		//sendEmailUsingPool(&qe, pool)
		err = edbmail.SendEmail(&m)

		if err != nil {
			log.Debug().Msgf("error sending email: %v", err)
		}
	}

	return nil
}
