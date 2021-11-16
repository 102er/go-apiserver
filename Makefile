GOPATH:=$(shell go env GOPATH)

.PHONY: proto1
# generate go protoc code
gp:
	protoc  --go_out=paths=source_relative:. \
			 api/v1/hello_world.proto

ggp:
	protoc  --gogo_out=paths=source_relative:. \
			 api/v1/hello_world.proto