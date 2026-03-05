module github.com/imadeddine-belkat/fpl-service

go 1.25

require (
	github.com/imadeddine-belkat/tactify-kafka v0.1.2
	github.com/imadeddine-belkat/tactify-protos v0.1.4
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
)

require (
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.25 // indirect
	github.com/segmentio/kafka-go v0.4.50 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/imadeddine-belkat/tactify-kafka => ../tactify-kafka

replace github.com/imadeddine-belkat/tactify-protos => ../tactify-protos
