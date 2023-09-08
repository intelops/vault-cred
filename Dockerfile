FROM  ghcr.io/sheril5/golang:1.19.4 AS builder
WORKDIR /workspace

COPY . ./

# Download dependencies
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o vault-cred cmd/main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /workspace/vault-cred vault-cred

USER 65532:65532
ENTRYPOINT ["./vault-cred"]
