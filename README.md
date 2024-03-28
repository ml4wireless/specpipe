# SpecPipe - Distributed Data Pipeline for Spectrum
In today's interconnected world, radio spectrum signals surround us, yet there exist noticeable limitations in the data systems created to access, monitor, perform AI experiments, and contribute to this analog data.

To democratize the access and usage of spectrum data, we have built SpecPipe, a modern scalable data pipeline for spectrum. This platform’s core values of accessibility, extensibility and scalability ensure that individual users can start to work with radio data with inexpensive hardware, minimal configuration, and a smooth onboarding process.

We have accomplished this goal of improving access to spectrum data by building SpecPipe as an open source project free for people to access and use, with easy to follow documentation, and a plethora of startup examples that allow users to understand our framework interactively.

SpecPipe leverages software-defined radio (SDR) to capture, process, and stream radio spectrum data in near real-time. It consists of three primary components:
- `rtl_rpcd` - runs on edge nodes, facilitating remote access to SDR hardware for `sp-edge`.
- `sp-edge` - operates within containers on edge nodes, managing SDR hardware to capture spectrum data. It processes and streams the data to the cloud.
- `sp-server` - provides a centralized control plane in the cloud to seamlessly orchestrate and manage edge devices at scale. It enables monitoring and management of devices and their data streams.

The system is designed for efficient, resilient spectrum data collection and processing across distributed edge nodes. It provides two main data pipelines:
- **IQ data pipeline**: for high-throughput raw IQ data retrieved from the RF front-end. This pipeline produces large amounts of raw data, which is temporarily stored in the cloud.
- **demodulated data pipeline**: has lower throughput for processed and demoulated data, saving bandwidth and storage compared to raw IQ data. Enables more cost-efficient data collection.

The dual pipelines allow flexible handling of diverse spectrum data types and throughput requirements across the distributed network.

Key Features:
- Transparent - server tracks edge node status, performs health checks, and exposes REST APIs for node management.
- Fault tolerant - automatic reconnections and timeouts between edge and cloud.
- Horizontally scalable - seamlessly orchestrates and manages a large, scalable number of edge devices.
- Dynamic configuration - allows dynamic device configuration on the fly, enabling flexibility for operating large clusters.
- Direct data stream forwarding - allows IQ edge devices to directly forward their data streams to specified servers through gRPC, without going through NATS. This reduces load and disk space used on NATS servers and provides more flexibility in routing data streams directly between devices.
- Portable - packages complex dependencies into a single container for easy edge and cloud deployment.
- Intuitive - user-friendly CLI with robust configuration options via file, CLI args, or environment variables.
- Lightweight - small binaries with low memory footprint.

The following is an overview of the system:

<img width="795" alt="image" src="https://github.com/ml4wireless/specpipe/assets/50090692/b95311ec-d323-495b-a22b-775c94d34138">

## Getting Started
Follow the tutorial on the [specpipe website](https://ml4wireless.github.io/specpipe)