# Raw IQ Data Mocking and ADSB

## Streaming Mock Raw IQ Data

This example program streams the content of a file (`exampledata.bin`) to NATS Jetstream containg raw IQ Data collected at the frequency of 1090000000, sampling rate of 2000000 Hz, and gain of 50 collected for 10 seconds on a loop.

This circular streaming generates continuous Raw IQ Data without needing an actual antenna.

Side note: The Raw IQ data was collected using the following call (which requires rtl_sdr installed - not part of this code but added for more explanation):
```bash
rtl_sdr -f 1090000000 -s 2000000 -g 50 output.bin
```

All of this code is contained in the `mock_iq` folder

## Decoding Raw IQ to ADS-B

A client `main.py` has been created that fetches this Raw IQ Data from Jetstream, processes it via [Dump1090](https://github.com/antirez/dump1090), decodes the output and prints valid ADSB signals to the console.

## Setup

### 1) Starting Mock Data Streamer (Terminal Window 1)
```
NATS_URL="nats://mytoken@127.0.0.1:4222"
DEVICE="dev1"
go run main.go exampledata.bin
```

We setup the NATS URL and DEVICE as environment variables that are used by `main.go`
We pass `exampledata.bin` to `main.go` so that's the data it plays on a loop

### 2) Processing Raw IQ to produce decoded ADSB Signals (Terminal Window 2)

```
pip install requirements.txt
NATS_URL="nats://127.0.0.1:4222"
DEVICE="dev1"
TOKEN="mytoken"
python3 main.py
```