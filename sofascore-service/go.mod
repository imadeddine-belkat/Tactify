module github.com/imadeddine-belkat/sofascore-service

go 1.26

require (
	github.com/chromedp/chromedp v0.14.2
	github.com/imadeddine-belkat/tactify-kafka v0.1.2
	github.com/imadeddine-belkat/tactify-protos v0.1.4
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	golang.org/x/sync v0.19.0
)

require (
	github.com/chromedp/cdproto v0.0.0-20250803210736-d308e07a266d // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20260214004413-d219187c3433 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.25 // indirect
	github.com/segmentio/kafka-go v0.4.50 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/imadeddine-belkat/tactify-kafka => ./../tactify-kafka

replace github.com/imadeddine-belkat/tactify-protos => ./../tactify-protos
