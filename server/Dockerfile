FROM golang:1.13.3 AS dev
WORKDIR /
COPY main.go /
RUN go mod init github.com/compose-spec/compatibility-test-suite/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main /main.go

FROM ubuntu AS run
COPY --from=dev main /
ENTRYPOINT ["/main"]
