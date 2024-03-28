FROM golang:1.20-alpine AS go_build

RUN mkdir -p /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$VERSION -w -s" -o /bin/sp-server ./cmd/server

FROM alpine:3.18.4
RUN apk update && apk add --no-cache ca-certificates

COPY --from=go_build /bin/sp-server /bin/sp-server

RUN mkdir -p /app
WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
ENTRYPOINT ["/bin/sp-server"]