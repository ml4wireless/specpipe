# SpecPipe - Distributed Data Pipeline for Spectrum
SpecPipe is an end-to-end distributed data pipeline that leverages software-defined radio (SDR) to capture, process, and stream radio spectrum data in near real-time. It consists of three primary components:
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
- Portable - packages complex dependencies into a single container for easy edge deployment.
- Intuitive - user-friendly CLI with robust configuration options via file, CLI args, or environment variables.
- Lightweight - small binaries with low memory footprint.

## System Architecture

<img width="1200" alt="system_architecture" src="images/system_architecture.png">

## Getting Started
Follow the tutorial on the [specpipe website](https://ml4wireless.github.io/specpipe)