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

<img width="821" alt="image" src="https://github.com/ml4wireless/specpipe/assets/50090692/5bd11fa2-d512-4a03-b0da-80f6fa2484ff">

## Getting Started
### Install Dependencies for librtlsdr
In order to extract raw data from the SDR hardware, the `librtlsdr` binaries have to be installed on the host machine. Before we install these binaries, we need to install `gcc`, `g++`, `make`, `cmake` and `libusb`.

#### For Linux Users
```bash
sudo apt-get update
sudo apt-get -y install build-essential cmake libusb-1.0-0-dev
```

#### For Mac Users using Apple silicon
Install  `cmake` and `libusb` via Homebrew.
```bash
brew install cmake libusb
```
Find the version and library paths of `libusb` so that they can be used in the next step
```bash
brew ls libusb
```
The two paths that we will need here are:
1) The path to the folder that contains `libusb.h` e.g.
```
/opt/homebrew/Cellar/libusb/1.0.26/include/libusb-1.0
```
2) The path to the `.dylib` file e.g.
```
/opt/homebrew/Cellar/libusb/1.0.26/lib/libusb-1.0.0.dylib
```


### Install librtlsdr
Build and install the librtlsdr binaries and libraries

#### Setup for installing librtlsdr

```bash
git clone https://github.com/minghsu0107/librtlsdr
cd librtlsdr
mkdir build && cd build
```

#### For Linux Users
```bash
cmake ../
sudo make && sudo make install
```
#### For Mac Users using Apple silicon
We need to set the appropriate configuration and library paths for the system when calling `cmake`. Specifically, we'll need to set the options:
1) `DLIBUSB_INCLUDE_DIR` to the path to the folder that contains `libusb.h` found above
2) `DLIBUSB_LIBRARY` to the path to the `.dylib` file found above

For Mac M2, an example command is:

```bash
cmake -DCMAKE_HOST_SYSTEM_PROCESSOR:STRING=arm64 -DLIBUSB_INCLUDE_DIR=/opt/homebrew/Cellar/libusb/1.0.26/include/libusb-1.0 -DLIBUSB_LIBRARY=/opt/homebrew/lib/libusb-1.0.dylib ../
sudo make && sudo make install
```

For Mac M1, an example command is the following. Note that the path to the libusb-1.0 is different for the machines.
```bash
cmake -DCMAKE_HOST_SYSTEM_PROCESSOR:STRING=arm64 -DLIBUSB_INCLUDE_DIR=/usr/local/Cellar/libusb/1.0.26/include/libusb-1.0 -DLIBUSB_LIBRARY=/usr/local/lib/libusb-1.0.dylib ../
sudo make && sudo make install
```

After building and installing librtlsdr, the files are located in the following directories:
- Header files are installed to `/usr/local/include`
- Library files are installed to `/usr/local/lib`
- Executable binaries are installed to `/usr/local/bin`

### Build Docker Image (Optional)
Building the Docker images locally is optional since prebuilt images are available on DockerHub.
Navigate to the root of the `specpipe` project and build the docker image.
```bash
export DOCKER_DEFAULT_PLATFORM=linux/amd64
make docker VERSION=v0.2.1
```
### Deployment

Start a `rtl_rpcd` daemon on the host machine, which allows remote access of SDR hardware at `127.0.0.1:40000` via `librtlsdr` command-line tools.

```bash
RTLSDR_RPC_SERV_ADDR=127.0.0.1 RTLSDR_RPC_SERV_PORT=40000 rtl_rpcd >> rtlrpcd.log 2>&1 &
```

Start NATS JetStream container and create stream `specpipe`, `specpipe-iq` and KV store `specpipe` respectively. This will also create Prometheus, Grafana, and NATS exporter containers, which are used for monitoring and alerting.

```bash
export DOCKER_DEFAULT_PLATFORM=linux/amd64
docker-compose up -d
```
Run `specpipe-edge` container, which retrieves raw data remotely from the `rtl_rpcd` daemon on the host machine and streams demoulated data to JetStream (take `fm` as example).

```bash
docker run --rm -d minghsu0107/specpipe-edge fm \
    --rpc-server-addr=host.docker.internal \
    --rpc-server-port=40000 \
    --nats-url=nats://mytoken@host.docker.internal:4222 \
    --device-name=dev1 \
    --sample-rate=170k \
    --freq=99700000
```
Note that `host.docker.internal` is used to access the host machine from the container, but only works on Mac and Windows. If you are using Linux, you can add argument `--network=host` and use `localhost` instead.

Start the API server (controller), which serves as the control plane enabling viewing of registered services and management of device configurations.

```bash
docker run --rm -p 80:8888 -d minghsu0107/specpipe-server controller \
    --http-server-port=8888 \
    --nats-url=nats://mytoken@host.docker.internal:4222
```

Start the API server healthcheck routine.

```bash
docker run --rm -d minghsu0107/specpipe-server health \
    --nats-url=nats://mytoken@host.docker.internal:4222
```
#### Cloud APIs
View configurations of all registered FM devices.
```bash
curl http://localhost/v0/fm/devices
```
Example response:
```
{"devices":[{"freq":"99700000","latitude":0,"longitude":0,"name":"dev1","register_ts":1708125826204,"resample_rate":"32k","sample_rate":"170k","specpipe_version":"v0.2.1"}]}
```

View configuration of a registered FM device.
```bash
curl http://localhost/v0/fm/devices/<device_name>
```
Example response:
```
{"device":{"freq":"99700000","latitude":0,"longitude":0,"name":"dev1","register_ts":1708125826204,"resample_rate":"32k","sample_rate":"170k","specpipe_version":"v0.2.1"}}
```

Update configuration of a registered FM device. For example, you could tune a device to frequency 94100000 with samping rate 200k on the fly.
```bash
curl -X PUT http://localhost/v0/fm/devices/<device_name> --data '{"freq":"94100000","sample_rate": "200k"}'
```

You could optionally run the Swagger UI to view all APIs in your browser at `http://localhost:5555`. Before running the following command, you should modify `server/openapi/main.yaml#/servers.url` from `/v0` to `http://localhost/v0` in order to make API's `Try it out` works.

```bash
docker run --rm -d -p 5555:8080 -e API_URL=api/main.yaml -v $(PWD)/server/openapi:/usr/share/nginx/html/api swaggerapi/swagger-ui
```
#### Grafana
Open `http://localhost:3000` in your browser to access Grafana. The default username and password are `admin` and `admin`. You could add the Prometheus datasource at `http://prometheus:9090`.
#### NATS JetStream Subjects
- Demodulated data pipeline: `specpipe.data.<sdr_type>.<device_name>`
  - Supported SDR types: `iq`, `fm`
- IQ data pipeline: `specpipe-iq.data.iq.<device_name>`
- Cluster commands: `specpipe-cluster.<sdr_type>.<device_name>.<cmd>`
  - Supported commands:
    - `health` - For health checks
    - `watchcfg` - For dynamic configuration
### More Examples
- [FM Audio Stream on NATS JetStream](./_examples/audio_play)
- [IQ Visualization & IQEngine Integration](./_examples/plot_iq)
- [Speech to Text](./_examples/speech2text)
- [Audio Data Mocking & Prometheus Exporter](./_examples/mock_fm)
