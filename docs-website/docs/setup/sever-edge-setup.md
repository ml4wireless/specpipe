---
title: Setup SpecPipe Edge and Server
sidebar_position: 1
---

# Setup SpecPipe Edge Device, Server, Controller Plane and Dashboards

Note: This guide only applies if you want to setup to send Radio Data as an [Edge Node](../architecture#edge-nodes). If you only want to build an application that accesses data, you don't need to have an SDR.

### 1. Install Dependencies for librtlsdr
In order to extract raw data from the SDR hardware, the `librtlsdr` binaries have to be installed on the host machine. Before we install these binaries, we need to install `gcc`, `g++`, `make`, `cmake` and `libusb`.

#### For Linux Users
```bash
sudo apt-get update
sudo apt-get -y install build-essential cmake libusb-1.0-0-dev
```

#### For Mac Users using Apple silicon
Install `cmake` and `libusb` via Homebrew.
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

### 2. Install librtlsdr
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
cmake -DCMAKE_HOST_SYSTEM_PROCESSOR:STRING=arm64 -DLIBUSB_INCLUDE_DIR=/opt/homebrew/Cellar/libusb/1.0.26/include/libusb-1.0 -DLIBUSB_LIBRARIES=/opt/homebrew/lib/libusb-1.0.dylib ../
sudo make && sudo make install
```

For Mac M1, an example command is the following. Note that the path to the libusb-1.0 is different for the machines.
```bash
cmake -DCMAKE_HOST_SYSTEM_PROCESSOR:STRING=arm64 -DLIBUSB_INCLUDE_DIR=/usr/local/Cellar/libusb/1.0.26/include/libusb-1.0 -DLIBUSB_LIBRARIES=/usr/local/lib/libusb-1.0.dylib ../
sudo make && sudo make install
```

After building and installing librtlsdr, the files are located in the following directories:
- Header files are installed to `/usr/local/include`
- Library files are installed to `/usr/local/lib`
- Executable binaries are installed to `/usr/local/bin`

### 3. Build Docker Image (Optional)
Building the Docker images locally is optional since prebuilt images are available on DockerHub.
Navigate to the root of the `specpipe` project and build the docker image.
```bash
export DOCKER_DEFAULT_PLATFORM=linux/amd64
make docker VERSION=v0.2.1
```

### 4. Allow remote access of SDR hardware within container
Start a `rtl_rpcd` daemon on the host machine, which allows remote access of SDR hardware at `127.0.0.1:40000` via `librtlsdr` command-line tools.

```bash
RTLSDR_RPC_SERV_ADDR=127.0.0.1 RTLSDR_RPC_SERV_PORT=40000 rtl_rpcd >> rtlrpcd.log 2>&1 &
```

### 5. Setting up SpecPipe server

Start NATS JetStream container and create stream `specpipe`, `specpipe-iq` and KV store `specpipe` respectively. This will also create Prometheus, Grafana, and NATS exporter containers, which are used for monitoring and alerting.

```bash
export DOCKER_DEFAULT_PLATFORM=linux/amd64
docker-compose up -d
```

#### NATS JetStream Subjects
- Demodulated data pipeline: `specpipe.data.<sdr_type>.<device_name>`
  - Supported SDR types: `iq`, `fm`
- IQ data pipeline: `specpipe-iq.data.iq.<device_name>`
- Cluster commands: `specpipe-cluster.<sdr_type>.<device_name>.<cmd>`
  - Supported commands:
    - `health` - For health checks
    - `watchcfg` - For dynamic configuration

### 6. Running SpecPipe Edge Device Device
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

### 7. Setting up Controller Plane
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
##### Cloud APIs
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

### 8. Viewing Graphana Dashboard
Open [http://localhost:3000](http://localhost:3000) in your browser to access Grafana. The default username and password are `admin` and `admin`. You could add the Prometheus datasource at [http://prometheus:9090](http://prometheus:9090)
