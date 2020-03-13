FROM golang:1.14 AS dev
WORKDIR /
COPY main.go /
RUN go mod init github.com/compose-spec/compatibility-test-suite/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main /main.go

FROM scratch AS run
COPY --from=dev main /
ENTRYPOINT ["/main"]
