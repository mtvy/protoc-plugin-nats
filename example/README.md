## Example

```sh
protoc --nats_out=. --plugin=../bin/protoc-gen-nats ./proto/hookah.proto

protoc --proto_path=. --go_out=. ./proto/hookah.proto
```