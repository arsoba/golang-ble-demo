Demo Golang BLE project

```bash
armhf: snapcraft --target-arch=armhf
```

## Start

 * Copy `./golang-ble_0.1_amd64.snap` in to the gateway home directory.
 * Run install command: `sudo snap install ~/golang-ble_0.1_amd64.snap --dangerous`
 * Connect plugs:

```bash
snap stop --disable golang-ble
snap connect golang-ble:network :network
snap connect golang-ble:bluetooth-control :bluetooth-control
snap connect golang-ble:network-control :network-control
snap start --enable golang-ble
```
