FROM golang:1.14 AS dev
WORKDIR /go/src
COPY . .
RUN go mod init github.com/compose-spec/conformance-tests/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server .

FROM scratch AS run
COPY --from=dev server /
ENTRYPOINT ["/server"]
