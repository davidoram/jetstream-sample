# NATS / Jetstream samples

Install nats-server (Jetstream ) locally

```
nats-server -v
nats-server: v2.2.0-beta.19
```


# Sample 1. Simplest config

## IOT message flows

2 iot devices
1 backend server

Messages:

- `iot.event` Sent from device -> backend
- `config.device` Sent from backend -> device

## Architecture

- All messages are public, no isolation of client1 from client2
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

## Architecture

- Messages for each client are isolated from one another
- Config moved to NATS
- Each client has to know and publish/subscribe to a separate `Subject` eg: `device.event.{client-id}`, & `config.changed.{device-id}`

## To run

In 4 sessions run the following:

`nats-server --config 2.config`

`go run server2/main.go -m 2 -u server -p a`

`go run client2/main.go -n 1 -u client1 -p b`

`go run client2/main.go -n 2 -u client2 -p c`
