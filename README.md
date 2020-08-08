# NATS / Jetstream samples

Install nats-server (Jetstream ) locally

```
git checkout git@github.com:nats-io/nats-server.git
git pull
go install
nats-server -v
nats-server: v2.2.0-beta.20
```


# Sample 1. Simplest config

## IOT message flows

2 iot devices
1 backend server

Messages:

- `iot.event` Sent from device -> backend
- `config.device` Sent from backend -> device

## Goal

- Messages are public
- Each client has to know and publish/subscribe to a separate `Subject` eg: `device.event.{client-id}`, & `config.changed.{device-id}`

## To run

In 4 sessions run the following:

`nats-server --config 1.config`

`go run server1/main.go -max-client 2`

`go run client1/main.go -client-num 1`

`go run client1/main.go -client-num 2`

# Sample 2. Use NATS Accounts for multi-tenancy

## IOT message flows

2 iot devices
1 backend server

Messages:

- `iot.event` Sent from device -> backend
- `config.device` Sent from backend -> device

## Goal

- Messages are private to each client
- Config moved to NATS
- Each client must know about its own unique identifier:
  - publishes to `iot.event.{client-id}`
  - subscribes to `config.changed.{client-id}`

## To run

In 4 sessions run the following:

```
cd sample2
docker image rm davidoram/sample2-client:latest davidoram/sample2-server:latest
docker-compose build
docker-compose up
```

# Sample 3. Use NATS Accounts for multi-tenancy

## IOT message flows

2 iot devices
1 backend server

Messages:

- `iot.event` Sent from device -> backend
- `config.device` Sent from backend -> device

## Goal

- Messages are private to each client
- Config moved to NATS
- Move entire knowlege of separate `{client-id}` into the configuration layer inside NATS
- Each client:
  - publishes to `iot.event`
  - subscribes to `config.changed`

## To run

In 4 sessions run the following:

```
cd sample3
docker image rm davidoram/sample3-client:latest davidoram/sample3-server:latest
docker-compose build
docker-compose up
```

https://max.pfingsthorn.de/news/2016/05/code-simulating-network-links-with-docker/
https://github.com/maxpfingsthorn/mini-network-simulator
