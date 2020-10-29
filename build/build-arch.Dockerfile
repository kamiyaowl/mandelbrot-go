FROM golang:1.15-buster

RUN apt update \
    && apt install -y --no-install-recommends \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libxcursor-dev \
    libgl1-mesa-dev \
    libxxf86vm-dev \
    && apt -y clean

WORKDIR /go/src