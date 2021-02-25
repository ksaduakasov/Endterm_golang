FROM golang
COPY . /app
WORKDIR /app
RUN go mod download
ENV GO111MODULE=on\
        CGO_ENABLED=0\
        GOOS=linux\
        GOARCH=amd64
RUN go build /app/greet/greet_server && chmod +x /app/greet/greet_server
ENTRYPOINT ["./greet_server"]
FROM alpine:latest
RUN apk --no cache add ca-certificates
COPY --from=0 /app .

CMD ["./greet_server"]

EXPOSE 4001