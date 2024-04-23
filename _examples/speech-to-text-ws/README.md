#### README

### Overview
Demo application for speech to text applications.

A simple FastAPI web server fetches FM data from NATS, decodes it into .wav
chunks suitable for the SpeechRecognition model, and streams the text to
any listeners via websocket. A demo frontend application is included which
listens to the websocket and outputs the live text to a scrollable text box.

### Prerequisites
NATS should be running with at least one device streaming to it. See main
project README for information on how to set this up, or use the mock_fm 
example to stream sample FM data without a physical device.

#### Setup
`python3.11 -m venv venv`
`source venv/bin/activate`
`pip install -r requirements.txt`

### Run Demo

#### Terminal 1
`python3 main.py`
or
`uvicorn main:app --reload`

#### Terminal 2
`cd frontend && npm install && npm start`