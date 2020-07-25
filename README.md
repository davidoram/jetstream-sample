# NATS / Jetstream samples

Install nats-server (Jetstream ) locally

```
nats-server -v
nats-server: v2.2.0-beta.19
```


# Sample 1. IOT message flows

2 iot devices
1 backend server

Messages:

- `iot.event` Sent from device -> backend
- `config.device` Sent from backend -> device

# To run

In 4 sessions run the following:

`nats-server --config 1.config`

`go run server1/main.go -max-client 2`

`go run client1/main.go -client-num 1`

`go run client1/main.go -client-num 1`
