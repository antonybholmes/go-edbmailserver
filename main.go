package main

import (
	"net/mail"

	"github.com/antonybholmes/go-mailer/sesmailserver"
	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/panjf2000/ants"
)

func init() {
	from := sys.Must(mail.ParseAddress(env.GetStr("FROM", "")))

	sesmailserver.Init(from)
}

func main() {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	env.Ls()

	// make a thread pool
	pool := sys.Must(ants.NewPool(10))

	ConsumeRedis(pool)

	//ConsumeKafka(pool)
}
