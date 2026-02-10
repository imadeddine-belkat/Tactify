.PHONY: setup gen clean

setup:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

gen:
	@echo "Generating FPL..."
	protoc --proto_path=proto --go_out=go/fpl/v1 --go_opt=paths=source_relative proto/fpl/v1/*.proto
	@echo "Generating Sofascore..."
	protoc --proto_path=proto --go_out=go/sofascore/v1 --go_opt=paths=source_relative proto/sofascore/v1/*.proto

clean:
	@echo "Cleaning..."
	-del /s /q go\*.pb.go