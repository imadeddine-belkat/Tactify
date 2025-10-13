module github.com/imadbelkat1/sofascore-service

go 1.25

require (
	github.com/chromedp/chromedp v0.14.2
	github.com/imadbelkat1/kafka v0.0.0-00010101000000-000000000000
	github.com/imadbelkat1/shared v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	golang.org/x/sync v0.17.0
)

require (
	github.com/chromedp/cdproto v0.0.0-20250803210736-d308e07a266d // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20250910080747-cc2cfa0554c3 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/segmentio/kafka-go v0.4.49 // indirect
	golang.org/x/sys v0.37.0 // indirect
)

replace github.com/imadbelkat1/kafka => ../kafka

replace github.com/imadbelkat1/shared => ../shared
