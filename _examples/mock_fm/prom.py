import os
import sys
from prometheus_client import start_http_server
from prometheus_client import Counter, Gauge

import asyncio
import nats
from nats.js import api


nats_url = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
device = os.environ.get("DEVICE", "")
data_counter = Counter(f'fm_data_bytes_total',
                       'FM audio data size', ['device'])


async def main():
    if device == "":
        print("device name cannot be empty")
        sys.exit(1)
    nats_subject = f'specpipe.data.fm.{device}'
    start_http_server(6060)

    nc = await nats.connect(nats_url)
    js = nc.jetstream()

    sub = await js.pull_subscribe(nats_subject, None,
                                  config=api.ConsumerConfig(api.DeliverPolicy.NEW, inactive_threshold=30))

    try:
        while True:
            msgs = await sub.fetch(100)
            for msg in msgs:
                data_counter.labels(device=device).inc(len(msg.data))
                await msg.ack()
    except Exception as e:
        print(e)
    finally:
        await sub.unsubscribe()
        await nc.close()

if __name__ == '__main__':
    asyncio.run(main())
