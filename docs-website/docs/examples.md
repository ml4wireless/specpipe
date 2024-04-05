---
title: Example Applications
sidebar_position: 3
---

# Example SpecPipe Applications

### 1. FM Audio
[This example](https://github.com/ml4wireless/specpipe/tree/main/_examples/audio_play) demonstrates how to subscribe to live audio data streams captured and processed by SpecPipe and play the audio chunks locally.

### 2. IQ Visualization & IQEngine Integration
[This example](https://github.com/ml4wireless/specpipe/tree/main/_examples/plot_iq) demonstrates how to visualize the IQ spectrum data captured by SpecPipe.

### 3. Speech to Text
[This example](https://github.com/ml4wireless/specpipe/tree/main/_examples/speech2text) demonstrates converting speech received via FM to text.

### 4. Audio Data Mocking & Prometheus Exporter
[This example](https://github.com/ml4wireless/specpipe/tree/main/_examples/mock_fm) streams the content of mock_audio.wav file (sampled at 32 KHz) to NATS JetStream circularly to simulate continuous audio collection. To monitor this data, this example further runs the Prometheus exporter.
