### Server Setup

#### Locally

- create `.env` see [example](.env.example)
- run `go build .`
- run `./iot-device-manager`

#### Docker

- create `.env` see [example](.env.example)
- ensure `MOCHI_BROKER` is not assigned in `.env`
- run `docker-compose up -d`
- run `docker logs iot-device-manager` to ensure everything is started correctly

### Mochi Setup

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
