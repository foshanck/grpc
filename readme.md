# gen files

protoc -I /usr/include  --proto_path=. --go_out=.  --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative deviceinfo/deviceInfo.proto