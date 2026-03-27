module github.com/imadeddine-belkat/indexer-service

go 1.25

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/imadeddine-belkat/tactify-kafka v0.0.0
	github.com/imadeddine-belkat/tactify-protos v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.11.2
)

require (
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pierrec/lz4/v4 v4.1.26 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/segmentio/kafka-go v0.4.50 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/imadeddine-belkat/tactify-kafka => ../tactify-kafka
	github.com/imadeddine-belkat/tactify-protos => ../tactify-protos
)