FROM golang:alpine

# RUN apk add --no-cache git openssh openssl curl
RUN apk add --no-cache git curl
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/ark120202/easy-elo-ranker/backend

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . ./

RUN go build -o app

EXPOSE 8080

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /go/src/github.com/ark120202/easy-elo-ranker/backend/app .
ENTRYPOINT ["./app"]
