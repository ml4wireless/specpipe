# IQ Data Mocking and ADSB

## Streaming Mock IQ Data

This example program streams the content of a file (`exampledata.bin`) to NATS Jetstream containg IQ Data collected at the frequency of 1090000000, sampling rate of 2000000 Hz, and gain of 50 collected for 10 seconds on a loop.

This circular streaming generates continuous IQ Data without needing an actual antenna.

Side note: The IQ data was collected using the following call (which requires `rtl_sdr` installed - not part of this code but added for more explanation):
```bash
rtl_sdr -f 1090000000 -s 2000000 -g 50 output.bin
```

All of this code is contained in the `mock_iq` folder.

## Decoding IQ to ADS-B

A client `main.py` has been created that fetches this IQ Data from Jetstream, processes it via [Dump1090](https://github.com/antirez/dump1090), decodes the output and prints valid ADSB signals to the console.

## Setup

### 1) Starting Mock Data Streamer (Terminal Window 1)

```bash
cd mock_iq/
NATS_URL="nats://mytoken@host.docker.internal:4222" DEVICE="dev1" go run main.go exampledata.bin
```

We setup `NATS_URL` and `DEVICE` as environment variables that are used by `main.go` and pass `exampledata.bin` to `main.go` so that's the data it plays on a loop.

### 2) Processing Raw IQ to produce decoded ADSB Signals (Terminal Window 2)

```bash
docker build -t specpipe-adsb-demo .
docker run -e NATS_URL="nats://mytoken@host.docker.internal:4222" -e DEVICE="dev1" -e PYTHONUNBUFFERED=1 --rm specpipe-adsb-demo
```