import asyncio

import websockets

DEVICE_ID = "dev0-mock"
SAMPLE_RATE = 64000

WS_URL = f"ws://0.0.0.0:8000/ws/fm_speech/{DEVICE_ID}/{SAMPLE_RATE}"
WS_TEST_URL = f"ws://0.0.0.0:8000/ws/text_random"


async def listen():
    """Listens to a websocket and prints the outputted text"""
    async with websockets.connect(WS_TEST_URL) as websocket:
        while True:
            message = await websocket.recv()
            print(message)


if __name__ == "__main__":
    try:
        asyncio.run(listen())
    except:
        pass
