import asyncio
import websockets

WS_URL = "ws://0.0.0.0:8000/ws_fm"

async def listen():
  async with websockets.connect(WS_URL) as websocket:
    while True:
      message = await websocket.recv()
      print(message)


if __name__ == "__main__":
  try:
    asyncio.run(listen())
  except:
    pass
