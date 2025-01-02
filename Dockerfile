FROM golang:alpine AS build

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -mod=readonly -ldflags="-s -w" -o /tmp/gmail

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=build /tmp/gmail /usr/bin/gmail

CMD ["/usr/bin/gmail"]
