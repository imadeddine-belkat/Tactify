module github.com/imadeddine-belkat/fpl-service

go 1.25

require (
	github.com/imadeddine-belkat/kafka v0.0.0-00010101000000-000000000000
	github.com/imadeddine-belkat/shared v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.49 // indirect
)

replace github.com/imadeddine-belkat/kafka => ../kafka

replace github.com/imadeddine-belkat/shared => ../shared
