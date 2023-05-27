FROM golang:1.19.4 AS builder
WORKDIR /workspace

# Set the GOPRIVATE environment variable
ARG GITHUB_TOKEN
RUN git config --global --add url."https://${GITHUB_TOKEN}:@github.com/intelops/go-common".insteadOf "https://github.com/intelops/go-common"
ENV GOPRIVATE=github.com/intelops/go-common

COPY . ./

# Download dependencies
RUN GOPRIVATE=$GOPRIVATE go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o vault-cred cmd/main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /workspace/vault-cred vault-cred

USER 65532:65532
ENTRYPOINT ["./vault-cred"]