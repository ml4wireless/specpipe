import asyncio
import os
import random
import string
from collections import deque
from io import BytesIO

import speech_recognition as sr
import uvicorn
from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from nats.aio.client import Client as NATS

app = FastAPI()

NATS_URL = os.getenv("NATS_URL", "nats://127.0.0.1:4222")
NATS_TOKEN = os.getenv("NATS_TOKEN", "mytoken")

FM_TOPIC_PREFIX = "specpipe.data.fm"
CHUNK_SIZE_SECONDS = 2


@app.get("/")
async def root():
    """Standard liveness check"""
    return {"message": "OK"}


@app.websocket("/ws/text_random")
async def ws_text_random(websocket: WebSocket):
    """Returns a random string of characters every second"""
    await websocket.accept()
    try:
        while True:
            random_text = " ".join(random.choices(string.ascii_letters, k=40))
            await websocket.send_text(random_text)
            await asyncio.sleep(1)
    except WebSocketDisconnect:
        pass


def create_wav_header(sample_rate, bits_per_sample, channels, num_samples):
    """Returns WAV file header in bytes for a given audio configuration"""
    datasize = num_samples * channels * (bits_per_sample // 8)
    o = bytes("RIFF", "ascii")  # (4byte) Marks file as RIFF
    o += (datasize + 36).to_bytes(4, "little")  # (4byte) File size in bytes
    o += bytes("WAVE", "ascii")  # (4byte) File type
    o += bytes("fmt ", "ascii")  # (4byte) Format Chunk Marker
    o += (16).to_bytes(4, "little")  # (4byte) Length of above format data
    o += (1).to_bytes(2, "little")  # (2byte) Format type (1 - PCM)
    o += (channels).to_bytes(2, "little")  # (2byte)
    o += (sample_rate).to_bytes(4, "little")  # (4byte)
    o += (sample_rate * channels * (bits_per_sample // 8)).to_bytes(4, "little")
    o += (channels * (bits_per_sample // 8)).to_bytes(2, "little")  # (2byte)
    o += (bits_per_sample).to_bytes(2, "little")  # (2byte)
    o += bytes("data", "ascii")  # (4byte) Data Chunk Marker
    o += (datasize).to_bytes(4, "little")  # (4byte) Data size in bytes
    return o


def prepend_wav_header(
    audio_data: bytes,
    sample_rate: int = 32000,
    bits_per_sample: int = 16,
    channels: int = 1,
):
    """Take raw audio data and prepend a WAV header to it, returning bytes"""
    num_samples = len(audio_data) // (channels * (bits_per_sample // 8))
    header = create_wav_header(sample_rate, bits_per_sample, channels, num_samples)
    return header + audio_data


async def async_speech_to_text(audio_data: bytes, sample_rate: int):
    """Given a chunk of bytes, return the transcribed text using the Sphinx engine"""
    audio_data_with_header = prepend_wav_header(audio_data, sample_rate=sample_rate)
    recognizer = sr.Recognizer()

    audio_file = BytesIO(audio_data_with_header)
    with sr.AudioFile(audio_file) as source:
        audio = recognizer.record(source)

    text = recognizer.recognize_sphinx(audio)
    return text


@app.websocket("/ws/fm_speech/{device_id}/{sample_rate}")
async def ws_fm_speech(
    websocket: WebSocket, device_id: str = "dev0-mock", sample_rate: int = 32000
):
    """Websocket endpoint that listens to audio data from a device and streams transcribed text"""
    await websocket.accept()

    nc = NATS()
    await nc.connect(NATS_URL, token=NATS_TOKEN)
    js = nc.jetstream()

    error_event = asyncio.Event()
    audio_buffer = deque()

    async def message_handler(msg):
        audio_buffer.extend(msg.data)
        chunk_size_bytes = sample_rate * 2 * CHUNK_SIZE_SECONDS  # two bytes per sample
        if len(audio_buffer) >= chunk_size_bytes:
            try:
                chunk_to_process = bytes(
                    [audio_buffer.popleft() for _ in range(chunk_size_bytes)]
                )
                transcribed_text = await async_speech_to_text(
                    chunk_to_process, sample_rate
                )
                await websocket.send_text(transcribed_text)
            except WebSocketDisconnect as e:
                error_event.set()
            except Exception as e:
                print(str(e))
                error_event.set()

    sub = await js.subscribe(f"{FM_TOPIC_PREFIX}.{device_id}", cb=message_handler)

    try:
        while not error_event.is_set():
            await asyncio.sleep(1)
    except Exception as e:
        print(f"An error occurred in main loop: {str(e)}")
    except BaseException as e:
        # Keyboard interrupt or system exit
        pass

    await sub.unsubscribe()
    await nc.close()
    await websocket.close()


if __name__ == "__main__":
    try:
        uvicorn.run(app, host="0.0.0.0", port=8000)
    except KeyboardInterrupt:
        pass
