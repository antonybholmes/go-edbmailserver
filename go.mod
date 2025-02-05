module github.com/antonybholmes/go-edb-server-mailer

go 1.23

replace github.com/antonybholmes/go-mailer => ../go-mailer

replace github.com/antonybholmes/go-sys => ../go-sys

require (
	github.com/antonybholmes/go-sys v0.0.0-20250113143747-03c4e3605208
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/aws/aws-sdk-go v1.55.6 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20250117234145-f0c83d229437
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/redis/go-redis/v9 v9.7.0
	golang.org/x/sys v0.29.0 // indirect
)
