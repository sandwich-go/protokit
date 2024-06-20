../../protokitgo/bin/osx/protoc_3.14.0 -I../../protokitgo/sdk/proto_google_3.14.0 \
--experimental_allow_proto3_optional \
-Iprotos --go_out=module=github.com/sandwich-go/protokit/option:. \
--plugin=protoc-gen-go=../../protokitgo/bin/osx/protoc-gen-go  \
protos/protokit/orm.proto \
protos/protokit/rpc.proto