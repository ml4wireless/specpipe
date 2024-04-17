import datetime
from datetime import datetime, date
import pytz
import sys

def form_number(pInputText):
    try:
        return float(pInputText.replace('\r', ''))
    except:
        return float(0)


def form_text(pInputText):
    return pInputText.replace('\r', '')


def print_stuff(pText):
    print("{:%Y%m%d %H:%M:%S} {}".format(datetime.now(), pText))

def timestamp():
    now = datetime.now(pytz.timezone('UTC'))
    return now.isoformat()

def log_msg(pText):
    print("{:%Y%m%d %H:%M:%S} {}".format(datetime.now(), pText))
    sys.stdout.flush()