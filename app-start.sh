#!/bin/sh

/app/tailscaled --tun=userspace-networking --socks5-server=localhost:1055 &
/app/tailscale up --auth-key=${TAILSCALE_AUTHKEY} --hostname=cloudrun-app
echo Tailscale started
ALL_PROXY=socks5://localhost:1055/ /app/trss