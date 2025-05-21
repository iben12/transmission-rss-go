FROM golang:1.19 AS backend

WORKDIR /app

RUN apt-get clean && apt-get update && apt-get install -y ca-certificates

COPY go.mod go.sum main.go ./
COPY . .

RUN go get

RUN CGO_ENABLED=0 go build -o ./bin/trss ./main.go


FROM node:16.10 AS frontend

WORKDIR /app

COPY package.json package-lock.json webpack.config.js ./
COPY frontend ./frontend

RUN npm i
RUN npm run build


FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Copy Tailscale binaries from the tailscale image on Docker Hub.
COPY --from=docker.io/tailscale/tailscale:stable /usr/local/bin/tailscaled /app/tailscaled
COPY --from=docker.io/tailscale/tailscale:stable /usr/local/bin/tailscale /app/tailscale
RUN mkdir -p /var/run/tailscale /var/cache/tailscale /var/lib/tailscale

WORKDIR /app

COPY --from=backend /app/bin/trss .
COPY static ./static
COPY --from=frontend /app/static/assets /app/static/assets

COPY app-start.sh /app


CMD ["/app/app-start.sh"]