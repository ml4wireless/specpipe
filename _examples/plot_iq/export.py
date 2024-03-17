import os
import sys
import requests
import logging

import asyncio
import nats
from nats.js import api


sigmf_meta_template = '''
{
    "global": {
        "core:datatype": "cu8",
        "core:sample_rate": %s,
        "core:hw": "RTL-SDR IQ example",
        "core:author": "minghsu0107",
        "core:version": "1.0.0"
    },
    "captures": [
        {
            "core:sample_start": 0,
            "core:frequency": %s
        }
    ],
    "annotations": []
}
'''


async def main():
    controller_api_base_url = os.environ.get(
        "CONTROLLER_API_BASE_URL", "http://localhost/v0")
    nats_url = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
    device = os.environ.get("DEVICE", "")
    if device == "":
        print("device name cannot be empty")
        sys.exit(1)

    sigmf_meta = ""
    try:
        resp = requests.get(f"{controller_api_base_url}/iq/devices/{device}")
        resp.raise_for_status()

        resp_data = resp.json()

        sample_rate = resp_data['device']['sample_rate']
        frequency = resp_data['device']['freq']

        sigmf_meta = sigmf_meta_template % (sample_rate, frequency)
    except requests.HTTPError as errh:
        logging.error("An Http Error occurred:" + repr(errh))
        sys.exit(1)
    except requests.exceptions.ConnectionError as errc:
        logging.error(
            "An Error Connecting to the API occurred:" + repr(errc))
        sys.exit(1)
    except requests.exceptions.Timeout as errt:
        logging.error("A Timeout Error occurred:" + repr(errt))
        sys.exit(1)
    except requests.exceptions.RequestException as err:
        logging.error("An Unknown Error occurred" + repr(err))
        sys.exit(1)

    nats_subject = f'specpipe-iq.data.iq.{device}'

    nc = await nats.connect(nats_url)
    js = nc.jetstream()

    sub = await js.subscribe(nats_subject, deliver_policy=api.DeliverPolicy.NEW)
    data = bytearray()
    try:
        for _ in range(191):
            msg = await sub.next_msg(timeout=3)
            await msg.ack()

            data.extend(msg.data)
        print("total data size (bytes): ", len(data))

        with open('iq_example.sigmf-data', 'wb') as f:
            f.write(data)

        with open('iq_example.sigmf-meta', 'w') as f:
            f.write(sigmf_meta)
    except Exception as e:
        print(e)
    finally:
        await sub.unsubscribe()
        await nc.close()


if __name__ == '__main__':
    asyncio.run(main())
