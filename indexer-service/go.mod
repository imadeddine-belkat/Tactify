module github.com/imadbelkat1/indexer-service

go 1.24

replace github.com/imadbelkat1/kafka => ../kafka

replace github.com/imadbelkat1/fpl-service => ../fpl-service

require (
	github.com/imadbelkat1/fpl-service v0.0.0-00010101000000-000000000000
	github.com/imadbelkat1/kafka v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.10.9
)
