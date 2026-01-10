module github.com/imadeddine-belkat/sofascore-service

go 1.25

require (
	github.com/chromedp/chromedp v0.14.2
	github.com/imadeddine-belkat/kafka v0.0.0-00010101000000-000000000000
	github.com/imadeddine-belkat/shared v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	golang.org/x/net v0.48.0
	golang.org/x/sync v0.19.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/chromedp/cdproto v0.0.0-20250803210736-d308e07a266d // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20251027170946-4849db3c2f7e // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.23 // indirect
	github.com/segmentio/kafka-go v0.4.49 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/imadeddine-belkat/kafka => ../kafka

replace github.com/imadeddine-belkat/shared => ../shared
