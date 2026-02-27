#!/bin/sh

# if [ "$USE_TAILSCALE" == "true" ]; then
    /tailscale/tailscaled --tun=userspace-networking --socks5-server=localhost:1055 &
    /tailscale/tailscale up --auth-key=${TAILSCALE_AUTHKEY} --hostname=${TAILSCALE_HOST}

    /tailscale set --exit-node=${TAILSCALE_EXIT_NODE} --exit-node-allow-lan-access
    echo Tailscale started
    export http_proxy=socks5://localhost:1055/
    # HTTP_PROXY=socks5://localhost:1055/ HTTPS_PROXY=socks5://localhost:1055/ /app/trss
# else
    # /app/trss
# fi