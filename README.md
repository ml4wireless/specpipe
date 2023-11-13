# SpecPipe - Distributed Data Pipeline for Spectrum
## Install Dependencies
In order to extract IQ data from the SDR hardware, the `librtlsdr` binaries have to be installed on the host machine.

First, have `gcc`, `g++`, and `make` installed. 

```bash
sudo apt-get update
sudo apt-get -y install build-essential
```

Then install `cmake` and `libusb`.

```bash
sudo apt-get -y install cmake libusb-1.0-0-dev
```

Next, build and install `librtlsdr` binaries and libraries.

```bash
git clone https://github.com/minghsu0107/librtlsdr
cd librtlsdr
mkdir build && cd build
cmake ../
sudo make && sudo make install
```
After building and installing librtlsdr, the files are located in the following directories:
- Header files are installed to `/usr/local/include`
- Library files are installed to `/usr/local/lib`
- Executable binaries are installed to `/usr/local/bin`
### For Mac Users (Apple Chips)
Install  `cmake` and `libusb` via Homebrew.
```bash
brew install cmake libusb
```
Check the version and library paths of `libusb`.
```bash
brew ls libusb
```
Build and install the librtlsdr binaries and libraries, setting the appropriate configuration and library paths for the system. For example, on Mac M2 with `libusb` version `1.0.26`:
```bash
git clone https://github.com/minghsu0107/librtlsdr
cd librtlsdr
mkdir build && cd build
cmake -DCMAKE_HOST_SYSTEM_PROCESSOR:STRING=arm64 -DLIBUSB_INCLUDE_DIR=/opt/homebrew/Cellar/libusb/1.0.26/include/libusb-1.0 -DLIBUSB_LIBRARY=/opt/homebrew/lib/libusb-1.0.dylib ../
sudo make && sudo make install
```
Another example when using Mac M1 with `libusb` version `1.0.26`:
```bash
git clone https://github.com/minghsu0107/librtlsdr
cd librtlsdr
mkdir build && cd build
cmake -DCMAKE_HOST_SYSTEM_PROCESSOR:STRING=arm64 -DLIBUSB_INCLUDE_DIR=/usr/local/Cellar/libusb/1.0.26/include/libusb-1.0 -DLIBUSB_LIBRARY=/usr/local/lib/libusb-1.0.dylib ../
sudo make && sudo make install
```
## Build Docker Image
```bash
make docker
```
## Getting Started
Start a `rtl_rpcd` daemon on the host machine, which allows remote access of SDR hardware at `127.0.0.1:40000` via `librtlsdr` command-line tools.

```bash
RTLSDR_RPC_SERV_ADDR=127.0.0.1 RTLSDR_RPC_SERV_PORT=40000 rtl_rpcd >> rtlrpcd.log 2>&1 &
```

Start NATS JetStream container and create stream `specpipe` and KV store `specpipe` respectively.

```bash
docker-compose up -d
```
Run `specpipe-edge` container, which retrieves IQ data remotely from the `rtl_rpcd` daemon on the host machine and streams demodulized data to JetStream (take `fm` as example).

```bash
docker run --rm -d minghsu0107/specpipe-edge fm \
    --rpc-server-addr=host.docker.internal \
    --rpc-server-port=40000 \
    --nats-url=nats://mytoken@host.docker.internal:4222 \
    --device-name=dev1 \
    --freq=99700000
```

Start the API server at `localhost:8888`, which serves as the control plane enabling viewing of registered services and management of device configurations.

```bash
docker run --rm -p 80:8888 -d minghsu0107/specpipe-server \
    --http-server-port=8888 \
    --nats-url=nats://mytoken@host.docker.internal:4222
```

View configurations of all registered FM devices.
```bash
curl http://localhost/v0/fm/devices
```
Example response:
```
{"devices":[{"freq":"99700000","latitude":0,"longitude":0,"name":"dev1","sample_rate":"170k"}]}
```

View configuration of a registered FM device.
```bash
curl http://localhost/v0/fm/devices/<device_name>
```
Example response:
```
{"device":{"freq":"99700000","latitude":0,"longitude":0,"name":"dev1","sample_rate":"170k"}}
```
Update configuration of a registered FM device.
```bash
curl -X PUT http://localhost/v0/fm/devices/<device_name> --data '{"freq":"94100000"}'
```

You can optionally run the Swagger UI to view all APIs in your browser at `http://localhost:5555`. Before running the following command, you should modify `server/openapi/main.yaml#/servers.url` from `/v0` to `http://localhost/v0` in order to make API's `Try it out` works.

```bash
docker run --rm -d -p 5555:8080 -e API_URL=api/main.yaml -v $(PWD)/server/openapi:/usr/share/nginx/html/api swaggerapi/swagger-ui
```
