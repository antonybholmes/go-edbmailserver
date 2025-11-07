module github.com/antonybholmes/go-edbmailserver

go 1.25

replace github.com/antonybholmes/go-mailserver => ../go-mailserver

replace github.com/antonybholmes/go-sys => ../go-sys

require (
	github.com/antonybholmes/go-sys v0.0.0-20251017225059-9b05304a0fc0
	github.com/aws/aws-sdk-go-v2 v1.39.6
	github.com/aws/aws-sdk-go-v2/config v1.31.17
	github.com/panjf2000/ants v1.3.0
	github.com/rs/zerolog v1.34.0
	golang.org/x/text v0.30.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.18.21 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.54.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.39.1 // indirect
	github.com/aws/smithy-go v1.23.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/segmentio/kafka-go v0.4.49 // indirect
	github.com/tiendc/go-deepcopy v1.7.1 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/excelize/v2 v2.10.0 // indirect
	github.com/xuri/nfp v0.0.2-0.20250530014748-2ddeb826f9a9 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/net v0.46.0 // indirect
)

require (
	github.com/antonybholmes/go-mailserver v0.0.0-20251022134000-989a996f1eac
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.13
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/redis/go-redis/v9 v9.16.0
	golang.org/x/sys v0.37.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
