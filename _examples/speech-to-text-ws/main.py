import asyncio
import os
from collections import deque
from io import BytesIO

from fastapi import FastAPI, WebSocket
from nats.aio.client import Client as NATS
import speech_recognition as sr

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "OK"}

@app.websocket("/text_random")
async def text_random(websocket: WebSocket):
    await websocket.accept()
    while True:
        random_bytes = str(os.urandom(10))
        await websocket.send_text(str(random_bytes))
        await asyncio.sleep(1)

def create_wav_header(sample_rate, bits_per_sample, channels, num_samples):
    """
    Create a WAV file header for a given set of audio properties, 
    assuming linear PCM audio.
    """
    datasize = num_samples * channels * (bits_per_sample // 8)
    o = bytes("RIFF",'ascii')                                                   # (4byte) Marks file as RIFF
    o += (datasize + 36).to_bytes(4, 'little')                                  # (4byte) File size in bytes excluding this and RIFF marker
    o += bytes("WAVE",'ascii')                                                  # (4byte) File type
    o += bytes("fmt ",'ascii')                                                  # (4byte) Format Chunk Marker
    o += (16).to_bytes(4, 'little')                                             # (4byte) Length of above format data
    o += (1).to_bytes(2, 'little')                                              # (2byte) Format type (1 - PCM)
    o += (channels).to_bytes(2, 'little')                                       # (2byte)
    o += (sample_rate).to_bytes(4, 'little')                                    # (4byte)
    o += (sample_rate * channels * (bits_per_sample // 8)).to_bytes(4, 'little')# (4byte)
    o += (channels * (bits_per_sample // 8)).to_bytes(2, 'little')              # (2byte)
    o += (bits_per_sample).to_bytes(2, 'little')                                # (2byte)
    o += bytes("data",'ascii')                                                  # (4byte) Data Chunk Marker
    o += (datasize).to_bytes(4, 'little')                                       # (4byte) Data size in bytes
    return o

def prepend_wav_header(audio_data, sample_rate=32000, bits_per_sample=16, channels=1):
    num_samples = len(audio_data) // (channels * (bits_per_sample // 8))
    header = create_wav_header(sample_rate, bits_per_sample, channels, num_samples)
    return header + audio_data

async def async_speech_to_text(audio_data):
    audio_data_with_header = prepend_wav_header(audio_data)
    recognizer = sr.Recognizer()
    audio_file = BytesIO(audio_data_with_header)
    with sr.AudioFile(audio_file) as source:
        audio = recognizer.record(source)

    text = recognizer.recognize_sphinx(audio)
    return text


@app.websocket("/ws_fm_speech")
async def ws_fm(websocket: WebSocket):
    await websocket.accept()

    nc = NATS()
    await nc.connect("nats://127.0.0.1:4222", token="mytoken")
    js = nc.jetstream()

    error_event = asyncio.Event()

    audio_buffer = deque()

    async def message_handler(msg):
        audio_buffer.extend(msg.data)
        if len(audio_buffer) >= 256000:
            try:
                chunk_to_process = bytes([audio_buffer.popleft() for _ in range(256000)])
                transcribed_text = await async_speech_to_text(chunk_to_process)
                await websocket.send_text(transcribed_text)
            except Exception as e:
                print(str(e))
                error_event.set()

    sub = await js.subscribe("specpipe.data.fm.dev1", cb=message_handler)

    try:
        while not error_event.is_set():
            await asyncio.sleep(1)
    except Exception as e:
        print(f"An error occurred in main loop: {e}")
    
    await sub.unsubscribe()
    await nc.close()

    try:
      await websocket.close()
    except RuntimeError:
        pass



if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
