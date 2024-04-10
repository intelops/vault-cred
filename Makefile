APP_NAME := vault-cred
BUILD := 0.1.1

gen-protoc:
	mkdir -p proto/pb/vaultcredpb
	protoc --go_out=proto/pb/vaultcredpb --go_opt=paths=source_relative \
    		--go-grpc_out=proto/pb/vaultcredpb --go-grpc_opt=paths=source_relative \
    		--proto_path=./proto vault-cred.proto

build:
	CGO_ENABLED=0 go build -o vault-cred cmd/main.go

docker-build:
	docker build -f Dockerfile -t ${APP_NAME}:${BUILD} .

clean:
	rm vault-cred
