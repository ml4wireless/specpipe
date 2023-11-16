FROM golang:1.20-alpine AS go_build

RUN mkdir -p /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/mock-audio .

FROM alpine:3.18.4
RUN apk update && apk add --no-cache ca-certificates

RUN mkdir -p /app

COPY --from=go_build /bin/mock-audio /bin/mock-audio
COPY --from=go_build /app/mock_audio.wav /app/mock_audio.wav

WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
ENTRYPOINT ["/bin/mock-audio"]
CMD ["/app/mock_audio.wav"]