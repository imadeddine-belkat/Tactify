module github.com/imadbelkat1/indexer-service

go 1.24

replace github.com/imadbelkat1/kafka => ../kafka

replace github.com/imadbelkat1/shared => ../shared

replace github.com/imadbelkat1/fpl-service => ../fpl-service

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/imadbelkat1/kafka v0.0.0-00010101000000-000000000000
	github.com/imadbelkat1/shared v0.0.0-00010101000000-000000000000
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
