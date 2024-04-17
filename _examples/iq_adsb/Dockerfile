# Use an official Ubuntu 18.04 base image
FROM python:3.11-slim-bullseye

# Install necessary packages
RUN apt-get update && \
    apt-get install -y \
    pkg-config \
    librtlsdr-dev \
    git \
    build-essential

RUN mkdir -p /app

# Clone the dump1090 repository
RUN git clone https://github.com/antirez/dump1090.git /app

# Change working directory
WORKDIR /app

# Compile dump1090
RUN make

# Copy your Python script into the Docker image
COPY main.py utils.py requirements.txt /app/

RUN pip3 install -r requirements.txt

# Set the command to run your script
CMD ["python3", "main.py"]