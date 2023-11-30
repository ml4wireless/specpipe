import os
import sys

import asyncio
import nats
from nats.js import api

import speech_recognition as sr
import io


async def main():
    nats_url = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
    device = os.environ.get("DEVICE", "")
    if device == "":
        print("device name cannot be empty")
        sys.exit(1)
    nats_subject = f'specpipe.data.fm.{device}'

    nc = await nats.connect(nats_url)
    js = nc.jetstream()

    sub = await js.pull_subscribe(nats_subject, None,
                                  config=api.ConsumerConfig(api.DeliverPolicy.NEW, inactive_threshold=30))
    r = sr.Recognizer()

    try:
        data = bytearray()
        msgs = await sub.fetch(100)
        for msg in msgs:
            data.extend(msg.data)
            await msg.ack()

        b = io.BytesIO(data)
        with sr.AudioFile(b) as source:
            audio = r.record(source)
            text = r.recognize_google(audio)
            print(text)
        b.close()

    except Exception as e:
        print(e)
    finally:
        await sub.unsubscribe()
        await nc.close()

if __name__ == '__main__':
    asyncio.run(main())

