# Mock Audio Streaming

This example program streams the content of `mock_audio.wav` file to NATS JetStream circularly to simulate continuous audio collection. This allows creating mock audio data for testing without requiring an actual antenna or audio source. The circular streaming approach generates a continuous audio data source.
## Docker Build
```bash
docker build -t minghsu0107/specpipe-audio-mock .
```
## Usage
First, start NATS JetStream container with proper configurations as shown in the [deployment section](../../README.md#deployment).

To simulate an FM device named `dev1` streaming audio to NATS JetStream:
```bash
NATS_URL="nats://mytoken@127.0.0.1:4222" DEVICE="dev1" go run main.go 
```
You can also run in Docker:
```bash
docker run --rm -e NATS_URL="nats://mytoken@host.docker.internal:4222" -e DEVICE="dev1" minghsu0107/specpipe-audio-mock
```
