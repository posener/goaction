
FROM golang:1.14.1-alpine3.11
RUN apk add git

COPY cmd/goaction /home/src/
WORKDIR /home/src
RUN go mod init action
RUN go build -o /bin/action

FROM alpine:3.11
RUN apk add git
COPY --from=0 /bin/action /bin/action

ENTRYPOINT [ "/bin/action" ]
