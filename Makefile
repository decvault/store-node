.PHONY: install-protoc
install-protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: api
api: install-protoc
	find api -name "*.proto" | xargs protoc \
  		--go_out=./internal/pb \
		--go-grpc_out=./internal/pb
