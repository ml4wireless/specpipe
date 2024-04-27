"""Sample websocket listener for speech-to-text demo app"""
import asyncio
import sys

import websockets
from websockets import WebSocketException

DEVICE_ID = "dev0-mock"
SAMPLE_RATE = 32000

WS_URL = f"ws://0.0.0.0:8000/ws/fm_speech"
WS_TEST_URL = f"ws://0.0.0.0:8000/ws/text_random"


async def listen() -> None:
    """Listens to a websocket and prints the outputted text"""
    ws_url = WS_TEST_URL

    # Use device_id if provided
    if len(sys.argv) >= 2:
        ws_url = WS_URL + "/" + sys.argv[1]

    # Use sample_rate if provided
    if len(sys.argv) >= 3:
        ws_url = WS_URL + "/" + sys.argv[2]

    # Listen to websocket and display results
    async with websockets.connect(ws_url) as websocket:
        while True:
            message = await websocket.recv()
            print(message)


if __name__ == "__main__":
    try:
        asyncio.run(listen())
    except (KeyboardInterrupt, WebSocketException):
        pass
