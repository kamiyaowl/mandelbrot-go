# minimal build image
FROM golang:1.15-buster as build_base
ENV DEBIAN_FRONTEND noninteractive
WORKDIR /go/src

# x11, glfw for ebiten
FROM build_base as build_ubuntu
RUN apt update \
    && apt install -y --no-install-recommends \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libxcursor-dev \
    libgl1-mesa-dev \
    libxxf86vm-dev \
    && apt -y clean

# dummy display for ebiten
FROM build_ubuntu as test_ubuntu
RUN apt update \
    && apt install -y --no-install-recommends \
    xorg \
    xserver-xorg-video-dummy \
    x11-apps \
    && apt -y clean