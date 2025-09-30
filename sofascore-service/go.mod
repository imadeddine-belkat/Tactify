module github.com/imadbelkat1/sofascore-service

go 1.24

require (
	github.com/imadbelkat1/kafka v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
)

replace github.com/imadbelkat1/kafka => ../kafka
