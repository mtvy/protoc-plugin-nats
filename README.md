# protoc-plugin-nats

## Description
Protoc plugin golang code generator for NATS (https://nats.io/)

## Build
```sh
go build -C ./ -o bin/protoc-gen-nats
```

## Generate
```sh
protoc --nats_out=. --plugin=./protoc-gen-nats/bin/protoc-gen-nats ./proto/your_file.proto
```