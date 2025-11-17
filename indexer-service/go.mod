module github.com/imadeddine-belkat/indexer-service

go 1.25

replace github.com/imadeddine-belkat/kafka => ../kafka

replace github.com/imadeddine-belkat/shared => ../shared

replace github.com/imadeddine-belkat/fpl-service => ../fpl-service

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/imadeddine-belkat/kafka v0.0.0-00010101000000-000000000000
	github.com/imadeddine-belkat/shared v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.10.9
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.49 // indirect
)
