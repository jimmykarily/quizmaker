#!/bin/bash

startsudo() {
    sudo -v
    ( while true; do sudo -v; sleep 50; done; ) &
    SUDO_PID="$!"
    trap stopsudo SIGINT SIGTERM
}
stopsudo() {
    kill "$SUDO_PID"
    trap - SIGINT SIGTERM
    sudo -k
}

startsudo

export IMAGE=quay.io/jimmykarily/kairos-3.1.1-rpi-kubecon-kiosk:latest

docker buildx build --load --platform linux/arm64 -t "${IMAGE}" . && \
  sudo rm -rf unpacked && \
  sudo luet util unpack --local "${IMAGE}" unpacked && \
  docker run --privileged -it \
    -v $PWD/unpacked:/image \
    -v $PWD/build:/build \
    --entrypoint "/bin/sh" \
    -e MODEL=rpi4 \
    -e SIZE="30000" \
    -e DEFAULT_ACTIVE_SIZE="5000" \
    -e STATE_SIZE="15000" \
    -e RECOVERY_SIZE="10000" \
    quay.io/kairos/osbuilder-tools:latest \
    /build-arm-image.sh --model rpi4 --directory "/image" /build/image.img

stopsudo
