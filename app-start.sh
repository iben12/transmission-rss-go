#!/bin/sh

if [ "$USE_PROXY" == "true" ]; then
    /app/tailscaled --tun=userspace-networking --socks5-server=localhost:1055 &
    /app/tailscale up --auth-key=${TAILSCALE_AUTHKEY} --hostname=cloudrun-app
    echo Tailscale started
    HTTP_PROXY=socks5://localhost:1055/ HTTPS_PROXY=socks5://localhost:1055/ /app/trss
else
    /app/trss
fi