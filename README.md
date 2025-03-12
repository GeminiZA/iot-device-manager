## Setup

### Server Setup

#### Locally

- create `.env` see [example](.env.example)
- run `go build .`
- run `./iot-device-manager`

#### Docker-Compose

- create `.env` see [example](.env.example)
- run `docker-compose up -d` or `MQTT_BROKER_PORT=1883 MQTT_BROKER_WD_PORT=8000 API_PORT=8080 docker-compose up -d` (defaults) to specify the ports to be used
- run `docker logs iot-device-manager` to ensure everything is started correctly

### Standalone Mochi Setup

- Create `config.yaml` in `startMochi/`
- Add users needed for devices and server
- run docker-compose up -d
- will always restart on unexpected shutdown

#### config.yaml

```
listeners:
  - type: "tcp"
    id: "mqtt"
    address: ":1883"

hooks:
  auth:
    ledger:
      auth:
        - username: "[username]"
          password: "[password]"
          allow: true
        - username: "[username]"
          password: "[password]"
          allow: true
```

# Usage

## Devices

Devices are references by `id` which is a `uint` and need to be added to the server by the http endpoint before the server will keep track of it's status and telemetry

## MQTT

All assets updates should be published to `assets/:id`
MQTT username and password must match that defined in `mochi-mqtt/config.yaml`
update messages are _json_ and of the following form:

```
{
  "status": string,
  "telemetry": object
}
```

Status is a string to allow extension on it's functionality
Telemetry is a json object to allow any sort of data to stored in it regardless of what device needs to publish

## HTTP

### Endpoints

#### GET (/assets/:id)

Returns the asset details including as _json_:

- Name
- Status
- Telemetry

```
{
  "name": string,
  "status": string,
  "telemetry": object
}
```

#### POST (/assets)

Creates a new device and subscribed to the relevant MQTT topic
Request body:

```
{
  "name": string,
  "id": int
}
```

id must be a _positive 32 bit integer_ and must be unique

Returns as _json_:

- Name
- ID

```
{
  "name": string,
  "id": int
}
```

#### PUT (/assets/:id)

Updates the details of the device and publishes the update the the relevant MQTT topic
Request body:

```
{
  "status": string,
  "telemetry": object
}
```

`telemetry` is a _json object_ as to allow any data to be stored

Returns only a status code

#### DELETE (/assets/:id)

Deletes the device from the database and subsequently disallows updates to that device from the MQTT broker
No Request Body

Returns only a status code
