rm -rf gen
protokitgo export --config=setting.yaml --log_level=0
rm -rf protos/rawdata
rm -rf gen/rawdata
rm -rf gen/golang/message_registry
rm -rf gen/golang/protokit/*.hack.go
go mod tidy -v