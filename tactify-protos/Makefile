.PHONY: setup gen clean

setup:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

gen:
	@echo "Generating FPL..."
	protoc --proto_path=proto --go_out=go/ --go_opt=paths=source_relative proto/fpl/v1/*.proto
	@echo "Generating PL..."
	protoc --proto_path=proto --go_out=go/ --go_opt=paths=source_relative proto/pl/v1/*.proto

clean:
	@echo "Cleaning..."
	-del /s /q go\*.pb.go/