FROM golang:1.19 as backend

WORKDIR /app

RUN apt-get clean && apt-get update && apt-get install -y ca-certificates

COPY go.mod go.sum main.go ./
COPY . .

RUN go get

RUN CGO_ENABLED=0 go build -o ./bin/trss ./main.go


FROM node:16.10 as frontend

WORKDIR /app

COPY package.json package-lock.json webpack.config.js ./
COPY frontend ./frontend

RUN npm i
RUN npm run build


FROM busybox

WORKDIR /app

COPY --from=backend /app/bin/trss .
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY static ./static
COPY --from=frontend /app/static/assets /app/static/assets

CMD ["/app/trss"]