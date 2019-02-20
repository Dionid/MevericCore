
# Meveric IoT

Set of Golang applications that give you control other your IoT devices.

## What is inside

1. MongoDB
1. NATS
1. MQTT
1. Device as Microservice
1. API as RPC other WebSockets and HTTP
1. MQTT routing
1. Echo extensions
1. Golang MongoDB Models and Collections overlay
1. RPC router
1. `Device shadow` similar to [AWS IoT architecture](https://docs.aws.amazon.com/en_us/iot/latest/developerguide/iot-device-shadows.html)

## Requirements

1. MongoDB
2. NATS server
3. Golang

## How to Start

1. Start `$ mongod`
1. Restore sample DB `$  mongorestore -d tztatom ./DB/tztatom/`
2. Start `$ gnatsd`
3. Start main Dashboard API `$ go run mcdashboard/main/main.go`
4. Start Plantainer API (as example of Device) `$ go run mcplantainer/main/main.go`