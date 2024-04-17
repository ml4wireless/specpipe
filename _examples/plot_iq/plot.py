import os
import sys

import asyncio
import nats
from nats.js import api

import numpy as np
import matplotlib.pyplot as plt


def load_data(data):
    # Read uint8 data and convert to double
    y = np.frombuffer(data, dtype='uint8')

    # Convert to range -127.5 to 127.5
    y = y - 127.5

    # Split into real and imaginary parts
    y_real = y[0::2]
    y_imag = y[1::2]

    # Create complex array
    y = y_real + 1j*y_imag

    return y


def plot_FFT_IQ(x, n0, nf, fs, f0, dst='iq.png'):
    global fig

    # Extract segment
    x_segment = x[n0:n0+nf]

    # Take FFT
    p = np.fft.fftshift(np.fft.fft(x_segment))

    # Normalize
    p = 20*np.log10(np.abs(p)/np.max(np.abs(p)))

    # Frequency axis
    low_freq = f0 - fs/2
    # high_freq = f0 + fs/2
    N = len(p)
    freq = np.arange(N)*(fs/N) + low_freq

    # Plot
    plt.plot(freq, p)
    plt.ylim(ymin=np.min(p), ymax=np.max(p))
    plt.xlabel('Frequency [MHz]', fontsize=14)
    plt.ylabel('Relative amplitude [dB]', fontsize=14)
    plt.grid()
    plt.rcParams['axes.facecolor'] = 'white'

    plt.title(f'Spectrum\nCenter frequency = {f0} MHz', fontsize=14)

    # Add center line
    plt.vlines(f0, ymin=np.min(p), ymax=np.max(p), colors='r', linewidth=2)

    plt.savefig(dst)


async def main():
    nats_url = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
    device = os.environ.get("DEVICE", "")
    if device == "":
        print("device name cannot be empty")
        sys.exit(1)
    nats_subject = f'specpipe-iq.data.iq.{device}'

    nc = await nats.connect(nats_url)
    js = nc.jetstream()

    sub = await js.subscribe(nats_subject, deliver_policy=api.DeliverPolicy.NEW)
    data = bytearray()
    try:
        for _ in range(3):
            msg = await sub.next_msg(timeout=3)
            await msg.ack()

            data.extend(msg.data)
        print("total data size of 3 messages (bytes): ", len(data))

        y = load_data(data)
        plot_FFT_IQ(y, 0, len(y), 2.5, 99.7)
    except Exception as e:
        print(e)
    finally:
        await sub.unsubscribe()
        await nc.close()


if __name__ == '__main__':
    asyncio.run(main())
