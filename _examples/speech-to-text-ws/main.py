"""Speech to text web server demo"""
import asyncio
import os
import random
import string
from io import BytesIO

import speech_recognition as sr
import uvicorn
from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from nats.aio.client import Client as NATS

app = FastAPI()

NATS_URL = os.getenv("NATS_URL", "nats://127.0.0.1:4222")
NATS_TOKEN = os.getenv("NATS_TOKEN", "mytoken")

FM_TOPIC_PREFIX = "specpipe.data.fm"
BUFFERED_CHUNKS = 16


@app.get("/")
async def root() -> dict[str, str]:
    """Standard liveness check"""
    return {"message": "OK"}


@app.websocket("/ws/text_random")
async def ws_text_random(websocket: WebSocket) -> None:
    """Returns a random string of characters every second"""
    await websocket.accept()
    try:
        while True:
            random_text = " ".join(random.choices(string.ascii_letters, k=40))
            await websocket.send_text(random_text)
            await asyncio.sleep(1)
    except (KeyboardInterrupt, WebSocketDisconnect):
        pass


def create_wav_header(
    sample_rate: int, num_samples: int, bits_per_sample: int, channels: int
) -> bytes:
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
    sample_rate: int,
    bits_per_sample: int = 16,
    channels: int = 1,
) -> bytes:
    """Take raw audio data and prepend a WAV header to it, returning bytes"""
    num_samples = len(audio_data) // (channels * (bits_per_sample // 8))
    header = create_wav_header(sample_rate, num_samples, bits_per_sample, channels)
    return header + audio_data


async def async_speech_to_text(audio_data: bytes, sample_rate: int) -> str:
    """Given a chunk of bytes, return the transcribed text using the Sphinx engine"""
    audio_data_with_header = prepend_wav_header(audio_data, sample_rate)
    recognizer = sr.Recognizer()

    audio_file = BytesIO(audio_data_with_header)
    with sr.AudioFile(audio_file) as source:
        audio = recognizer.record(source)

    text = recognizer.recognize_sphinx(audio)
    return str(text)


@app.websocket("/ws/fm_speech/{device_id}/{sample_rate}")
async def ws_fm_speech(
    websocket: WebSocket, device_id: str = "dev0-mock", sample_rate: int = 32000
) -> None:
    """Websocket endpoint that listens to audio data from a device and streams transcribed text"""
    await websocket.accept()

    nc = NATS()
    await nc.connect(NATS_URL, token=NATS_TOKEN)
    js = nc.jetstream()
    sub = await js.subscribe(f"{FM_TOPIC_PREFIX}.{device_id}")

    try:
        while True:
            audio_buffer = bytearray()
            try:
                for _ in range(BUFFERED_CHUNKS):
                    msg = await sub.next_msg(timeout=1)
                    audio_buffer.extend(msg.data)
            except asyncio.TimeoutError:
                pass
            transcribed_text = await async_speech_to_text(audio_buffer, sample_rate)
            await websocket.send_text(transcribed_text)
    except (KeyboardInterrupt, WebSocketDisconnect):
        pass
    finally:
        await sub.unsubscribe()
        await nc.close()
        try:
            await websocket.close()
        except RuntimeError:
            pass


if __name__ == "__main__":
    try:
        uvicorn.run(app, host="0.0.0.0", port=8000)
    except KeyboardInterrupt:
        pass
