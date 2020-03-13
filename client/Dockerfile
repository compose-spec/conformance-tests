FROM golang:1.14 AS dev
WORKDIR /go/src
COPY . .
RUN go mod init github.com/compose-spec/conformance-tests/client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /client .

FROM scratch AS run
COPY --from=dev client /
ENTRYPOINT ["/client"]
