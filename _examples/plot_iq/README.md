# IQ Visualization & IQEngine Integration
This example demonstrates how to visualize the IQ spectrum data captured by SpecPipe.
## Prerequisites
Follow the deployment instructions in [README.md](../../README.md#deployment) to set up SpecPipe.
## Capture IQ Data
Capture IQ data at 99.7 MHz with a 2.5 MHz sampling rate:

```bash
docker run --rm -d minghsu0107/specpipe-edge iq \
    --rpc-server-addr=host.docker.internal \
    --rpc-server-port=40000 \
    --nats-url=nats://mytoken@host.docker.internal:4222 \
    --device-name=dev1 \
    --sample-rate=2500000 \
    --freq=99700000
```
### Plot Spectrum
Install dependencies:
```bash
pip3 install -r requirements.txt
```
Run the plot script:
```bash
NATS_URL="nats://mytoken@127.0.0.1:4222" DEVICE="dev1" python3 plot.py
```
The spectrum visualization will be saved to `iq.png` by default.
## IQEngine Integration
You can export the IQ data in [SigMF](https://github.com/sigmf/SigMF) format that can be imported into [IQEngine](https://iqengine.org/browser) for further visualization as follows:

<img width="1466" alt="image" src="https://github.com/minghsu0107/specpipe/assets/50090692/00976b65-0811-48ce-8730-09f796eab229">

For example, the following command exports the latest 50 MB of IQ data from the device `dev1`:
```bash
CONTROLLER_API_BASE_URL="http://localhost/v0" NATS_URL="nats://mytoken@127.0.0.1:4222" DEVICE="dev1" python3 export.py
```
The exported data file will be saved to `iq_example.sigmf-data`, and the metadata file will be saved to `iq_example.sigmf-meta`. An example SigMF file pair is available at `example_data/` for reference.
