# FM Audio Stream on NATS JetStream

This example demonstrates how to subscribe to live audio data streams captured and processed by SpecPipe and play the audio chunks locally.

## Dependencies
To enable real-time local playback of FM audio for demonstration purpose, the following dependencies are required.

For Debian / Ubuntu Linux:
```bash
apt-get install -y pkg-config portaudio19-dev
```
For OS X:
```bash
brew install pkg-config portaudio
```
## Usage
To run this example, first follow the deployment instructions in [README.md](../../README.md#deployment) to set up SpecPipe.

Start the FM audio subscriber on your host machine to listen to the audio stream in real time. The subscriber will connect to the NATS JetStream container and play the FM audio processed by SpecPipe.
```bash
NATS_URL="nats://mytoken@127.0.0.1:4222" DEVICE="dev1" go run main.go
```