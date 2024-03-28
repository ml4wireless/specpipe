# Audio Data Mocking & Prometheus Exporter

This example program streams the content of `mock_audio.wav` file (sampled at 32 KHz) to NATS JetStream circularly to simulate continuous audio collection. This allows creating mock audio data for testing without requiring an actual antenna or audio source. The circular streaming approach generates a continuous audio data source.
## Docker Build
```bash
docker build -t minghsu0107/specpipe-audio-mock .
```
## Usage
First, start NATS JetStream container with proper configurations as shown in [README.md](../../README.md).

To simulate an FM device named `dev1` streaming audio to NATS JetStream:
```bash
NATS_URL="nats://mytoken@127.0.0.1:4222" DEVICE="dev1" go run main.go mock_audio.wav
```
You can also run in Docker:
```bash
docker run --rm -e NATS_URL="nats://mytoken@host.docker.internal:4222" -e DEVICE="dev1" minghsu0107/specpipe-audio-mock
```
## Prometheus Exporter
To monitor the FM audio data size, you can run the Prometheus exporter.

First, install dependencies:
```bash
pip3 install -r requirements.txt
```

Then, run the exporter:
```bash
export NATS_URL="nats://mytoken@host.docker.internal:4222"
export DEVICE="dev1"
python3 prom.py
```
The counter metrics `fm_data_bytes_total` will be available at `http://localhost:6060/metrics`.
