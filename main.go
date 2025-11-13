package main

import (
	"net/mail"

	"github.com/antonybholmes/go-edbmailserver/consts"

	"github.com/antonybholmes/go-mailserver/sesmailserver"
	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/antonybholmes/go-sys/log"
	"github.com/panjf2000/ants"
)

// func initLogger() {
// 	fileLogger := &lumberjack.Logger{
// 		Filename:   fmt.Sprintf("logs/%s.log", consts.AppId),
// 		MaxSize:    10,   // Max size in MB before rotating
// 		MaxBackups: 3,    // Keep 3 backup files
// 		MaxAge:     7,    // Retain files for 7 days
// 		Compress:   true, // Compress old log files
// 	}

// 	multiWriter := io.MultiWriter(os.Stderr, fileLogger)

// 	logger := zerolog.New(multiWriter).With().Timestamp().Logger()

// 	// we use != development because it means we need to set the env variable in order
// 	// to see debugging work. The default is to assume production, in which case we use
// 	// lumberjack
// 	if os.Getenv("APP_ENV") != "development" {
// 		//zerolog.SetGlobalLevel(zerolog.InfoLevel)
// 		logger = zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)).With().Timestamp().Logger()
// 	}

// 	log.Logger = logger
// }

func init() {
	log.SetAppName(consts.AppId)
	//initLogger()

	env.Ls()

	from := sys.Must(mail.ParseAddress(consts.EmailFrom))

	sesmailserver.InitSesMailer(from)
}

func main() {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	// make a thread pool
	pool := sys.Must(ants.NewPool(10))

	defer pool.Release()

	log.Debug().Msgf("%s %s", consts.AppId, consts.Version)

	consumer := NewSQSConsumer(pool)

	//ConsumeRedis(pool)
	consumer.ConsumeSQS()
	//consumeKafka(pool)
}

// func consumeKafka(pool *ants.Pool) {

// 	var ctx = context.Background()

// 	r := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:   []string{"localhost:9094"}, // Kafka broker
// 		Topic:     mailserver.QUEUE_EMAIL_CHANNEL, // Topic name
// 		Partition: 0,
// 		MaxBytes:  10e6, // 10MB
// 	})

// 	defer r.Close()

// 	var qe mailserver.QueueEmail

// 	for {
// 		m, err := r.ReadMessage(ctx)

// 		if err != nil {
// 			fmt.Printf("%s", err)
// 			break
// 		}

// 		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

// 		err = json.Unmarshal([]byte(m.Value), &qe)

// 		if err != nil {
// 			log.Debug().Msgf("email error")
// 		}

// 		sendEmail(&qe, pool)
// 	}
// }

// func sendEmailUsingPool(m *mailserver.MailItem, pool *ants.Pool) {
// 	pool.Submit(func() { sendEmail(m) })
// }
