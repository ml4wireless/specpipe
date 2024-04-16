import os
import asyncio
import nats
import subprocess
from nats.aio.errors import ErrNoServers
import tempfile
import pickle
import re
from utils import *

async def run():
    # Configuration based on environment variables
    nats_url = os.getenv('NATS_URL', 'nats://127.0.0.1:4222')
    device_name = os.getenv('DEVICE', "dev1")
    token = os.getenv('TOKEN', "mytoken")
    if not device_name:
        print("device name cannot be empty")
        return
    
    subject = f'specpipe-iq.data.iq.{device_name}'

    # Connect to NATS
    try:
        nc = await nats.connect(nats_url, token = "mytoken")
        js = nc.jetstream()
    except ErrNoServers:
        print("No NATS servers available")
        return

    async def message_handler(msg):
        data = msg.data
        process_iq_data(data)
        
    # Subscribe to the subject
    sub = await js.pull_subscribe(subject)

    try:
        while True:
            # Wait for a message
            msgs = await sub.fetch(100, timeout=30)
            for msg in msgs:
                await message_handler(msg)
    finally:
        await sub.unsubscribe()
        await nc.close()

def process_iq_data(data):
    command = './dump1090'

    temp_file = tempfile.NamedTemporaryFile(delete=False, mode='wb')
    pickle.dump(data, temp_file)
    temp_file.close()

    temp_file_path = temp_file.name

    full_command = [command, "--ifile", temp_file_path]

    try:
        result = subprocess.run(full_command, check=True, text=True, capture_output=True)
        parse_adsb_message_reports(result.stdout)
        # if result.stdout:
        #     print(result.stdout)
    except subprocess.CalledProcessError as e:
        print("An error occurred while executing the command.")
        print("Error message:", e.stderr)
    finally:
        os.remove(temp_file.name)

def parse_adsb_message_reports(data):
    if len(data) >= 1:
        print("----- Decoded Signal -----")
        print(data)
        # End(?) of block of info
        searchICAO = re.search(
            r'(ICAO Address   :\s+|ICAO Address:\s+)([\w\d]+)( \(Mode S / ADS-B\))?$', data, re.M | re.I)
        searchFeet = re.search(
            r'(Altitude :\s+|Baro altitude:\s+)([\d.]+)( feet| ft)(.*$)', data, re.M | re.I)
        searchLatitude = re.search(
            r'(Latitude :\s+|CPR latitude:\s+)([\d.-]+)( \(\d+\))?$', data, re.M | re.I)
        searchLongitude = re.search(
            r'(Longitude:\s+|CPR longitude:\s+)([\d.-]+)( \(\d+\))?$', data, re.M | re.I)
        searchIdent = re.search(
            r'(Identification :\s+|Ident:\s+)(.*$)', data, re.M | re.I)
        
        if searchICAO and searchIdent:
            valICAO = formText(searchICAO.group(2))
            valIdent = formText(searchIdent.group(2)).strip()
            logmsg(f"valICAO:{valICAO} valIdent:{valIdent}")

        if searchFeet and searchICAO and searchLatitude and searchLongitude:
            # Found a valid combination
            valICAO = formText(searchICAO.group(2))
            valFeet = formNumber(searchFeet.group(2))
            valLatitude = formNumber(searchLatitude.group(2))
            valLongitude = formNumber(searchLongitude.group(2))
            logmsg(
                f"ICAO:{valICAO} Feet:{valFeet} Latitude:{valLatitude} Longitude:{valLongitude}")


if __name__ == '__main__':
    asyncio.run(run())