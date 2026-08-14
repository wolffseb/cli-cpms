# cli-cpms

A terminal charge point management system for a single OCPP charging station.

`cpms` is one static binary that:

- speaks **OCPP 1.6-J** (2.0.1 later) to the station, so you can reserve a connector,
  redeem the reservation with a known RFID tag, stop the session, and unlock the cable;
- exposes a minimal **OCPI 2.3.0 CPO** interface on the LAN — `versions`, `credentials`,
  `locations` (sender) and `commands` (receiver) — so a backend can register, pull static
  location data, and send `RESERVE_NOW` that ends up at the real charger;
- **PATCHes** EVSE status to the registered counterparty whenever it changes, whether the
  change came from the terminal or from someone physically plugging in a car;
- needs no database, no container and no web UI — one hand-edited config file, one binary.

## How the charger connects

This trips people up, so it is worth stating plainly: **in OCPP-J the charger is the
WebSocket client.** It dials us; we do not dial it. Configure the station's CSMS URL as:

```
ws://<host running cpms>:9000/ocpp/<charge point id>
```

`charger.ip` in the config is only used for a reachability pre-check and for display.

## Quickstart

```sh
make build                                   # -> bin/cpms
cp config.example.yaml config.yaml           # then edit it
./bin/cpms config validate -c config.yaml
```

A valid config prints the facts you need to get the station and the counterparty
talking to you:

```
config.yaml is valid.

  Charge point        ALP-HYC-001  (OCPP 1.6)  at 192.168.1.42
  OCPP listener       0.0.0.0:9000
  Station CSMS URL    ws://<this-host>:9000/ocpp/ALP-HYC-001
  Heartbeat timeout   90s
  OCPI listener       0.0.0.0:8080
  OCPI advertised as  http://192.168.1.10:8080/ocpi
  OCPI party          DE*FRY
  Default RFID tag    04A1B2C3D4
  Location            OFFICE-01 "Fryte HQ", Berlin (DEU) — 2 EVSEs, 2 connectors

  EVSE           EVSE ID        OCPP CONNECTOR  CONNECTORS
  ALP-HYC-001-1  DE*FRY*E001*1  1               IEC_62196_T2_COMBO/CABLE/DC 920V 500A
  ALP-HYC-001-2  DE*FRY*E001*2  2               IEC_62196_T2_COMBO/CABLE/DC 920V 500A
```

An invalid one exits 1 and names every problem by its field path, all in one pass:

```
config.yaml is not valid:

  auth.default_id_tag: "04A1 B2C3" is not a valid OCPP idTag; want 1-20 printable ASCII characters without spaces
  charger.id: is required
  location.evses[1].ocpp_connector_id: 1 is already mapped by location.evses[0]

3 problems found.
```

## Configuration

Everything lives in one YAML file — see [`config.example.yaml`](config.example.yaml),
which is commented field by field and is itself checked by the test suite.

Two files, one rule: **`config.yaml` is read-only to the tool.** It holds what you decide
(the charger, the static location data, and the two OCPI tokens that are ours to choose).
`state.json` holds what the tool learns at runtime (the counterparty's token and endpoints,
active reservations, last known status) and is the only file cpms writes.

That split is why `ocpi.token_c` is a fixed value in config rather than being generated
during the OCPI handshake as implementations usually do.

## Development

```sh
make ci      # gofmt check, go vet, golangci-lint, go test -race
make test
make lint    # needs: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2
```

## Status

Built in tracked steps; each lands as its own PR.

| # | Scope | State |
|---|---|---|
| 1 | Project skeleton, config loader, CI | done |
| 2 | OCPP 1.6-J CSMS: WebSocket server and message routing | next |
| 3 | Built-in charge point simulator | |
| 4 | Outbound OCPP commands, reservation lifecycle, one-shot CLI | |
| 5 | OCPI foundation: versions + credentials handshake | |
| 6 | OCPI Locations module (sender) | |
| 7 | PATCH location on status change | |
| 8 | OCPI Commands module (receiver) | |
| 9 | Terminal UI | |
| 10 | OCPP 2.0.1 | |
| 11 | Docs and release builds | |

Out of scope for now: the OCPI Booking module, CDRs, sessions, tariffs, payment terminals,
multiple OCPI versions, and any web UI.
